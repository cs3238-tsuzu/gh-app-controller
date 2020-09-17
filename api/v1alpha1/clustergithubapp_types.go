/*
Copyright 2020 modoki-paas.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"context"
	"encoding/base64"
	"os"

	"github.com/modoki-paas/ghapp-controller/pkg/ghatypes"
	"golang.org/x/xerrors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type PrivateKeySecretRef struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// ClusterGitHubAppSpec defines the desired state of ClusterGitHubApp
type ClusterGitHubAppSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// URL is the base url for GitHub
	// +kubebuilder:validation:Optional
	URL string `json:"url"`

	// AppID is the id of GitHub App
	AppID int64 `json:"appID"`

	// PrivateKeySecretRef is the reference for the secret of GitHub App's private key
	PrivateKeySecretRef PrivateKeySecretRef `json:"privateKeySecretRef"`
}

// ClusterGitHubAppStatus defines the observed state of ClusterGitHubApp
type ClusterGitHubAppStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ClusterGitHubApp is the Schema for the clustergithubapps API
// +kubebuilder:resource:path=clustergithubapps,scope=Cluster
type ClusterGitHubApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterGitHubAppSpec   `json:"spec,omitempty"`
	Status ClusterGitHubAppStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterGitHubAppList contains a list of ClusterGitHubApp
type ClusterGitHubAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterGitHubApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterGitHubApp{}, &ClusterGitHubAppList{})
}

var _ ghatypes.GitHubAppInterface = &ClusterGitHubApp{}

func (a *ClusterGitHubApp) GetURL() string {
	return a.Spec.URL
}

func (a *ClusterGitHubApp) GetAppID() int64 {
	return a.Spec.AppID
}

func (a *ClusterGitHubApp) GetPrivateKey(ctx context.Context, c client.Client) ([]byte, error) {
	secret := &corev1.Secret{}

	err := c.Get(ctx, client.ObjectKey{
		Name:      a.Spec.PrivateKeySecretRef.Name,
		Namespace: os.Getenv("CONTROLLER_NAMESPACE"),
	}, secret)

	if errors.IsNotFound(err) {
		return nil, ghatypes.ErrSecretNotFound
	}

	if err != nil {
		return nil, err
	}

	data, ok := secret.Data[a.Spec.PrivateKeySecretRef.Key]

	if !ok {
		return nil, ghatypes.ErrKeyNotFound
	}

	decoded, err := base64.StdEncoding.DecodeString(string(data))

	if err != nil {
		return nil, xerrors.Errorf("failed to parse base64-encoded key: %w", err)
	}

	return decoded, nil
}
