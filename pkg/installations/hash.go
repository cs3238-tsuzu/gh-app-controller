package installations

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/cs3238-tsuzu/ghapp-controller/api/v1alpha1"
	"github.com/cs3238-tsuzu/ghapp-controller/pkg/ghatypes"
	corev1 "k8s.io/api/core/v1"
)

type installationTokenClaims struct {
	URL            string
	AppID          int64
	PrivateKey     string
	InstallationID int64
	RepositoryIDs  []int64
	Permissions    *v1alpha1.InstallationPermissions
	ExpiresAt      int64
	Template       corev1.Secret
}

func getInstallationTokenHash(
	app ghatypes.GitHubAppInterface,
	installation *v1alpha1.Installation,
	privateKey []byte,
	expiredAt time.Time,
) string {

	b, err := json.Marshal(installationTokenClaims{
		URL:            app.GetURL(),
		AppID:          app.GetAppID(),
		PrivateKey:     string(privateKey),
		InstallationID: installation.Spec.InstallationID,
		RepositoryIDs:  installation.Spec.RepositoryIDs,
		Permissions:    installation.Spec.Permissions,
		ExpiresAt:      expiredAt.Unix(),
		Template:       installation.Spec.Template,
	})

	if err != nil {
		panic(err)
	}

	data := sha1.Sum(b)

	return base64.StdEncoding.EncodeToString(data[:])
}
