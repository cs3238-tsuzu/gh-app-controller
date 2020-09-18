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
	"github.com/google/go-github/v30/github"
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

// GetGitHubPermissions returns github.InstallationPermissions converted from InstallationPermissions
func (p *InstallationPermissions) GetGitHubPermissions() *github.InstallationPermissions {
	if p == nil {
		return nil
	}

	return &github.InstallationPermissions{
		Administration:              p.Administration,
		Blocking:                    p.Blocking,
		Checks:                      p.Checks,
		Contents:                    p.Contents,
		ContentReferences:           p.ContentReferences,
		Deployments:                 p.Deployments,
		Emails:                      p.Emails,
		Followers:                   p.Followers,
		Issues:                      p.Issues,
		Metadata:                    p.Metadata,
		Members:                     p.Members,
		OrganizationAdministration:  p.OrganizationAdministration,
		OrganizationHooks:           p.OrganizationHooks,
		OrganizationPlan:            p.OrganizationPlan,
		OrganizationPreReceiveHooks: p.OrganizationPreReceiveHooks,
		OrganizationProjects:        p.OrganizationProjects,
		OrganizationUserBlocking:    p.OrganizationUserBlocking,
		Packages:                    p.Packages,
		Pages:                       p.Pages,
		PullRequests:                p.PullRequests,
		RepositoryHooks:             p.RepositoryHooks,
		RepositoryProjects:          p.RepositoryProjects,
		RepositoryPreReceiveHooks:   p.RepositoryPreReceiveHooks,
		SingleFile:                  p.SingleFile,
		Statuses:                    p.Statuses,
		TeamDiscussions:             p.TeamDiscussions,
		VulnerabilityAlerts:         p.VulnerabilityAlerts,
	}
}

type MetadataSpec struct {
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +kubebuilder:validation:Optional
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: http://kubernetes.io/docs/user-guide/annotations
	// +kubebuilder:validation:Optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

// SecretTemplateSpec is the template to generate secret with the installation token
type SecretTemplateSpec struct {
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +kubebuilder:validation:Optional
	MetadataSpec `json:"metadata,omitempty"`

	// Immutable, if set to true, ensures that data stored in the Secret cannot
	// be updated (only object metadata can be modified).
	// If not set to true, the field can be modified at any time.
	// Defaulted to nil.
	// This is an alpha field enabled by ImmutableEphemeralVolumes feature gate.
	// +kubebuilder:validation:Optional
	Immutable *bool `json:"immutable,omitempty"`

	// Data contains the secret data. Each key must consist of alphanumeric
	// characters, '-', '_' or '.'. The serialized form of the secret data is a
	// base64 encoded string, representing the arbitrary (possibly non-string)
	// data value here. Described in https://tools.ietf.org/html/rfc4648#section-4
	// +kubebuilder:validation:Optional
	Data map[string][]byte `json:"data,omitempty"`

	// stringData allows specifying non-binary secret data in string form.
	// It is provided as a write-only convenience method.
	// All keys and values are merged into the data field on write, overwriting any existing values.
	// It is never output when reading from the API.
	// +k8s:conversion-gen=false
	// +kubebuilder:validation:Optional
	StringData map[string]string `json:"stringData,omitempty"`

	// Used to facilitate programmatic handling of secret data.
	// +kubebuilder:validation:Optional
	Type corev1.SecretType `json:"type,omitempty"`
}

// InstallationSpec defines the desired state of GitHub installation
type InstallationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// AppRef is a reference to ClusterGitHubApp or GitHubApp
	AppRef corev1.ObjectReference `json:"appRef"`

	// InstallationID is an installation id for GitHub App
	InstallationID int64 `json:"installationID"`

	// RepositoryIDS are used to restrict permissions for tokens
	// +kubebuilder:validation:Optional
	RepositoryIDs []int64 `json:"repositoryIDs,omitempty"`

	// Permissions are used to restrict permissions for tokens
	// +kubebuilder:validation:Optional
	Permissions *InstallationPermissions `json:"permissions"`

	// Key is the key in the secret to save the token
	Key string `json:"key"`

	// Template is the template to generate secret with the installation token
	Template SecretTemplateSpec `json:"template"`
}

// InstallationStatus defines the observed state of Installation
type InstallationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Ready is the status of the installation token
	Ready bool `json:"ready"`

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
