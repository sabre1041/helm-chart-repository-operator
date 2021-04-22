/*
Copyright 2021.

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

// HelmChartSpec defines the desired state of HelmChart
type HelmChartSpec struct {

	// Name represents the name of the chart
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Chart name"
	Name string `json:"name"`

	// Versions represents the list of chart versions
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Chart versions"
	Versions []HelmChartVersion `json:"versions"`

	// RepositoryName represents the name of the repository
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Repository name"
	RepositoryName string `json:"repositoryName"`

	// RepositoryDisplayName represents a friendly name of the repository
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Repository display name"
	RepositoryDisplayName string `json:"repositoryDisplayName,omitempty"`
}

// HelmChartStatus defines the observed state of HelmChart
type HelmChartStatus struct {

	// LastUpdateTimestamp represents the time the resource was last updated
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=status,displayName="Last Update Time"
	LastUpdateTimestamp *metav1.Time `json:"lastUpdateTimestamp,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Repository",type=string,JSONPath=".spec.repositoryName",description="Chart Repository"
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=".spec.name",description="Chart Name"
// +kubebuilder:printcolumn:name="Latest Version",type=string,JSONPath=".spec.versions[*].version",description="Latest Chart Version"
// +kubebuilder:resource:path=helmcharts,scope=Cluster

// HelmChart is the Schema for the helmcharts API
type HelmChart struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HelmChartSpec   `json:"spec,omitempty"`
	Status HelmChartStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HelmChartList contains a list of HelmChart
type HelmChartList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HelmChart `json:"items"`
}

type HelmChartVersion struct {

	// Version represents the version of the chart
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Chart version"
	Version string `json:"version"`

	// Created represents the time the chart was created
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Chart creation time"
	Created *metav1.Time `json:"created,omitempty"`

	// Description contains a one-sentence description of the chart
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Chart description"
	Description string `json:"description,omitempty"`

	// Digest represents a hash of the chart package archive
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Chart digest"
	Digest string `json:"digest,omitempty"`

	// ApiVersion represents the Chart API
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Chart API version"
	ApiVersion string `json:"apiVersion"`

	// Keywords represents a list of string keywords
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Chart keywords"
	Keywords []string `json:"keyword,omitempty"`

	// AppVersion represents the version of the application enclosed inside of this chart.
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Chart application version"
	AppVersion string `json:"appVersion,omitempty"`

	// Home represents the URL to a relevant project page, git repo, or contact person
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Chart home"
	Home string `json:"home,omitempty"`

	// Icon represents the URL to an icon file.
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Chart icon"
	Icon string `json:"icon,omitempty"`

	// Sources are the URLs to the source code of this chart
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Chart sources"
	Sources *[]string `json:"sources,omitempty"`

	// A list of name and URL/email address combinations for the maintainer(s)
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Chart maintainers"
	Maintainers *[]HelmChartMaintainer `json:"maintainers,omitempty"`

	// Dependencies are a list of dependencies for a chart.
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Chart dependencies"
	Dependencies *[]HelmChartDependency `json:"dependencies,omitempty"`

	// Type specifies the chart type: application or library
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Chart type"
	Type string `json:"type,omitempty"`

	// URLs is the list of Chart URLs
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Chart URL's"
	URLs []string `json:"urls,omitempty"`

	// KubeVersion is a SemVer constraint specifying the version of Kubernetes required.
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Applicable Kubernetes version"
	KubeVersion string `json:"kubeVersion,omitempty"`
}

type HelmChartMaintainer struct {

	// Name is a user name or organization name
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Maintainer name"
	Name string `json:"name,omitempty"`

	// Email is an optional email address to contact the named maintainer
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Maintainer email"
	Email string `json:"email,omitempty"`

	// URL is an optional URL to an address for the named maintainer
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Maintainer URL"
	URL string `json:"url,omitempty"`
}

type HelmChartDependency struct {
	// Name is the name of the dependency.
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Dependency name"
	Name string `json:"name"`

	// Version is the version (range) of this chart.
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Dependency version"
	Version string `json:"version,omitempty"`

	// Repository is the URL to the chart repository.
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Dependency repository"
	Repository string `json:"repository"`

	// Condition is a yaml path that resolves to a boolean, used for enabling/disabling charts
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Dependency conditions"
	Condition string `json:"condition,omitempty"`

	// Tags can be used to group charts for enabling/disabling together
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Dependency tags"
	Tags []string `json:"tags,omitempty"`

	// Enabled bool determines if chart should be loaded
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Enabled dependency"
	Enabled bool `json:"enabled,omitempty"`

	// Alias represents the usable alias to be used for the chart
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Alias of dependency"
	Alias string `json:"alias,omitempty"`
}

func init() {
	SchemeBuilder.Register(&HelmChart{}, &HelmChartList{})
}
