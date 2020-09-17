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

	"github.com/modoki-paas/ghapp-controller/pkg/ghatypes"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GitHubAppSpec defines the desired state of GitHubApp
type GitHubAppSpec struct {
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

// GitHubAppStatus defines the observed state of GitHubApp
type GitHubAppStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// GitHubApp is the Schema for the githubapps API
type GitHubApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GitHubAppSpec   `json:"spec,omitempty"`
	Status GitHubAppStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GitHubAppList contains a list of GitHubApp
type GitHubAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GitHubApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GitHubApp{}, &GitHubAppList{})
}

var _ ghatypes.GitHubAppInterface = &GitHubApp{}

func (a *GitHubApp) GetURL() string {
	return a.Spec.URL
}

func (a *GitHubApp) GetAppID() int64 {
	return a.Spec.AppID
}

func (a *GitHubApp) GetPrivateKey(ctx context.Context, c client.Client) ([]byte, error) {
	secret := &corev1.Secret{}

	err := c.Get(ctx, client.ObjectKey{
		Name:      a.Spec.PrivateKeySecretRef.Name,
		Namespace: a.Namespace,
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

	return data, nil
}
