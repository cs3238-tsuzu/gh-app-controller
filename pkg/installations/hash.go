package installations

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/modoki-paas/ghapp-controller/api/v1alpha1"
	"github.com/modoki-paas/ghapp-controller/pkg/ghatypes"
)

type installationTokenClaims struct {
	URL            string
	AppID          int64
	PrivateKey     string
	InstallationID int64
	RepositoryIDs  []int64
	Permissions    *v1alpha1.InstallationPermissions
	ExpiresAt      int64
	Template       v1alpha1.SecretTemplateSpec
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
