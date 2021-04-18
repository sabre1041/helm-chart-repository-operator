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
	Name                  string             `json:"name"`
	Versions              []HelmChartVersion `json:"versions"`
	RepositoryName        string             `json:"repositoryName"`
	RepositoryDisplayName string             `json:"repositoryDisplayName,omitempty"`
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
	Version      string                 `json:"version"`
	Created      *metav1.Time           `json:"created,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Digest       string                 `json:"digest,omitempty"`
	ApiVersion   string                 `json:"apiVersion"`
	Keywords     []string               `json:"keyword,omitempty"`
	AppVersion   string                 `json:"appVersion,omitempty"`
	Home         string                 `json:"home,omitempty"`
	Icon         string                 `json:"icon,omitempty"`
	Sources      *[]string              `json:"sources,omitempty"`
	Maintainers  *[]HelmChartMaintainer `json:"maintainers,omitempty"`
	Dependencies *[]HelmChartDependency `json:"dependencies,omitempty"`
	Type         string                 `json:"type,omitempty"`
	URLs         []string               `json:"urls,omitempty"`
	KubeVersion  string                 `json:"kubeVersion,omitempty"`
}

type HelmChartMaintainer struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	URL   string `json:"url,omitempty"`
}

type HelmChartDependency struct {
	// Name is the name of the dependency.
	//
	// This must mach the name in the dependency's Chart.yaml.
	Name string `json:"name"`
	// Version is the version (range) of this chart.
	//
	// A lock file will always produce a single version, while a dependency
	// may contain a semantic version range.
	Version string `json:"version,omitempty"`
	// The URL to the repository.
	//
	// Appending `index.yaml` to this string should result in a URL that can be
	// used to fetch the repository index.
	Repository string `json:"repository"`
	// A yaml path that resolves to a boolean, used for enabling/disabling charts (e.g. subchart1.enabled )
	Condition string `json:"condition,omitempty"`
	// Tags can be used to group charts for enabling/disabling together
	Tags []string `json:"tags,omitempty"`
	// Enabled bool determines if chart should be loaded
	Enabled bool `json:"enabled,omitempty"`
	// Alias usable alias to be used for the chart
	Alias string `json:"alias,omitempty"`
}

func init() {
	SchemeBuilder.Register(&HelmChart{}, &HelmChartList{})
}
