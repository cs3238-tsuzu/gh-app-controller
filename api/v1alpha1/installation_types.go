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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// InstallationPermissions is the permissions to restrict permissions for tokens
type InstallationPermissions struct {
	Administration              *string `json:"administration,omitempty"`
	Blocking                    *string `json:"blocking,omitempty"`
	Checks                      *string `json:"checks,omitempty"`
	Contents                    *string `json:"contents,omitempty"`
	ContentReferences           *string `json:"content_references,omitempty"`
	Deployments                 *string `json:"deployments,omitempty"`
	Emails                      *string `json:"emails,omitempty"`
	Followers                   *string `json:"followers,omitempty"`
	Issues                      *string `json:"issues,omitempty"`
	Metadata                    *string `json:"metadata,omitempty"`
	Members                     *string `json:"members,omitempty"`
	OrganizationAdministration  *string `json:"organization_administration,omitempty"`
	OrganizationHooks           *string `json:"organization_hooks,omitempty"`
	OrganizationPlan            *string `json:"organization_plan,omitempty"`
	OrganizationPreReceiveHooks *string `json:"organization_pre_receive_hooks,omitempty"`
	OrganizationProjects        *string `json:"organization_projects,omitempty"`
	OrganizationUserBlocking    *string `json:"organization_user_blocking,omitempty"`
	Packages                    *string `json:"packages,omitempty"`
	Pages                       *string `json:"pages,omitempty"`
	PullRequests                *string `json:"pull_requests,omitempty"`
	RepositoryHooks             *string `json:"repository_hooks,omitempty"`
	RepositoryProjects          *string `json:"repository_projects,omitempty"`
	RepositoryPreReceiveHooks   *string `json:"repository_pre_receive_hooks,omitempty"`
	SingleFile                  *string `json:"single_file,omitempty"`
	Statuses                    *string `json:"statuses,omitempty"`
	TeamDiscussions             *string `json:"team_discussions,omitempty"`
	VulnerabilityAlerts         *string `json:"vulnerability_alerts,omitempty"`
}

// InstallationSpec defines the desired state of GitHub installation
type InstallationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// AppRef is a reference to ClusterGitHubApp or GitHubApp
	AppRef corev1.ObjectReference `json:"appRef"`

	// RepositoryIDS are used to restrict permissions for tokens
	// +kubebuilder:validation:Optional
	RepositoryIDS []string `json:"repositoryIDs,omitempty"`

	// InstallationPermissions are used to restrict permissions for tokens
	// +kubebuilder:validation:Optional
	InstallationPermissions *InstallationPermissions `json:"installationPermissions"`

	// Key is the key in the secret to save the token
	Key string `json:"key"`

	// Template is the template to generate secret with the installation token
	Template corev1.Secret `json:"template"`
}

// InstallationStatus defines the observed state of Installation
type InstallationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Ready is the status of the installation token
	Ready bool `json:"bool"`

	// Secret is the secret name to save the installation token
	Secret string `json:"secret"`

	// Message is the error message if something failed
	// +kubebuilder:validation:Optional
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Installation is the Schema for the installations API
type Installation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstallationSpec   `json:"spec,omitempty"`
	Status InstallationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// InstallationList contains a list of Installation
type InstallationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Installation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Installation{}, &InstallationList{})
}
