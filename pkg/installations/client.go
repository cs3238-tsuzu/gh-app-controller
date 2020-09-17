package installations

import (
	"context"
	"net/http"
	"time"

	"github.com/google/go-github/v30/github"
	"github.com/modoki-paas/ghapp-controller/api/v1alpha1"
	"github.com/modoki-paas/ghapp-controller/pkg/appstransport"
	"github.com/modoki-paas/ghapp-controller/pkg/ghatypes"
	"golang.org/x/xerrors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	errUnknownResousrce = xerrors.New("unknown resource")
)

const (
	// SecretHashAnnotationKey is a key for 'sha-1' in annotations of generated secrets
	// The controller determines the secret is up-to-date by this hash value
	SecretHashAnnotationKey = "ghapp.tsuzu.dev/sha1"

	// SecretExpiredAtAnnotationKey is a key for 'expired-at' in annotations of generated secrets
	SecretExpiredAtAnnotationKey = "ghapp.tsuzu.dev/expired-at"
)

// SecretStatus is a status of the generated secret
type SecretStatus int

const (
	NotExisting SecretStatus = iota
	Desired
	Undesirable
)

type Client struct {
	Client        client.Client
	Installation  *v1alpha1.Installation
	RefreshBefore time.Duration
}

func (c *Client) getGitHubApp(ctx context.Context) (ghatypes.GitHubAppInterface, error) {
	var gha ghatypes.GitHubAppInterface
	var err error

	ref := c.Installation.Spec.AppRef
	switch ref.APIVersion {
	case v1alpha1.GroupVersion.String(), "":
		switch ref.Kind {
		case "ClusterGitHubApp":
			cgha := &v1alpha1.ClusterGitHubApp{}
			err = c.Client.Get(ctx, client.ObjectKey{
				Name: ref.Name,
			}, cgha)

			gha = cgha
		default:
			err = errUnknownResousrce
		}
	default:
		err = errUnknownResousrce
	}

	if err != nil {
		return nil, xerrors.Errorf("failed to get GitHub App: %w", err)
	}

	return gha, nil
}

func (c *Client) getGitHubAppClient(ctx context.Context, gha ghatypes.GitHubAppInterface) (ghclient *github.Client, err error) {
	b, err := gha.GetPrivateKey(ctx, c.Client)

	if err != nil {
		return nil, xerrors.Errorf("failed to get private key: %w", err)
	}

	tr, err := appstransport.NewAppsTransportFromBytes(gha.GetAppID(), b)

	if err != nil {
		return nil, xerrors.Errorf("failed to initialize apps transport: %w", err)
	}

	url := gha.GetURL()
	if url == "" {
		ghclient = github.NewClient(&http.Client{
			Transport: tr,
		})
	} else {
		ghclient, err = github.NewEnterpriseClient(
			url, url,
			&http.Client{
				Transport: tr,
			},
		)

		if err != nil {
			return nil, xerrors.Errorf("failed to parse GitHub url: %w", err)
		}
	}

	return ghclient, nil
}

func (c *Client) hasDesiredSecret(ctx context.Context, gha ghatypes.GitHubAppInterface, key []byte) (SecretStatus, *time.Time, error) {
	secret := &corev1.Secret{}
	err := c.Client.Get(ctx, client.ObjectKey{
		Name:      c.Installation.GetName(),
		Namespace: c.Installation.GetNamespace(),
	}, secret)

	if errors.IsNotFound(err) {
		return NotExisting, nil, nil
	}

	if err != nil {
		return NotExisting, nil, xerrors.Errorf("failed to get secret: %w", err)
	}

	hash, ok := secret.Annotations[SecretHashAnnotationKey]

	if !ok {
		return Undesirable, nil, nil
	}

	et, ok := secret.Annotations[SecretExpiredAtAnnotationKey]

	if !ok {
		return Undesirable, nil, nil
	}

	expiredAt, err := time.Parse(et, time.RFC3339)

	calculatedHash := getInstallationTokenHash(
		gha,
		c.Installation,
		key,
		expiredAt,
	)

	if hash != calculatedHash {
		return Undesirable, nil, nil
	}

	if time.Now().Add(c.RefreshBefore).After(expiredAt) {
		return Undesirable, nil, nil
	}

	return Desired, &expiredAt, nil
}

func (c *Client) generate(ctx context.Context, gha ghatypes.GitHubAppInterface, token *github.InstallationToken, privateKey []byte) *corev1.Secret {
	generated := c.Installation.Spec.Template

	generated.Name = c.Installation.Name
	generated.Namespace = c.Installation.Namespace
	generated.GenerateName = ""

	if generated.Annotations == nil {
		generated.Annotations = map[string]string{}
	}

	hash := getInstallationTokenHash(
		gha,
		c.Installation,
		privateKey,
		token.GetExpiresAt(),
	)

	generated.Annotations[SecretExpiredAtAnnotationKey] = token.GetExpiresAt().Format(time.RFC3339)
	generated.Annotations[SecretHashAnnotationKey] = hash
	generated.StringData[c.Installation.Spec.Key] = token.GetToken()

	return &generated
}

func (c *Client) Run(ctx context.Context) (SecretStatus, *corev1.Secret, *time.Time, error) {
	gha, err := c.getGitHubApp(ctx)

	if err != nil {
		return NotExisting, nil, nil, xerrors.Errorf("GitHub App is not available: %w", err)
	}

	ghclient, err := c.getGitHubAppClient(ctx, gha)

	if err != nil {
		return NotExisting, nil, nil, xerrors.Errorf("failed to initialize GitHub client: %w", err)
	}

	privateKey, err := gha.GetPrivateKey(ctx, c.Client)

	if err != nil {
		return NotExisting, nil, nil, xerrors.Errorf("failed to get private key for GitHub App: %w", err)
	}

	status, expiredAt, err := c.hasDesiredSecret(ctx, gha, privateKey)

	if err != nil {
		return NotExisting, nil, nil, xerrors.Errorf("checking current secret state failed: %w", err)
	}

	if status == Desired {
		return Desired, nil, expiredAt, nil
	}

	token, _, err := ghclient.Apps.CreateInstallationToken(
		ctx,
		c.Installation.Spec.InstallationID,
		&github.InstallationTokenOptions{
			RepositoryIDs: c.Installation.Spec.RepositoryIDs,
			Permissions:   c.Installation.Spec.Permissions.GetGitHubPermissions(),
		},
	)

	if err != nil {
		return status, nil, nil, xerrors.Errorf("failed to create installation token: %w", err)
	}

	generated := c.generate(ctx, gha, token, privateKey)

	return status, generated, token.ExpiresAt, nil
}
