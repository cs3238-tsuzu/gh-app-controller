/*
Copyright 2020 cs3238-tsuzu.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	URL                 string              `json:"url"`
	AppID               string              `json:"appID"`
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
