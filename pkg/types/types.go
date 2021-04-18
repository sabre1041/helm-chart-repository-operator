package types

import (
	helmv1beta1 "github.com/openshift/api/helm/v1beta1"
	"helm.sh/helm/v3/pkg/repo"
)

type HelmChartEntry struct {
	Name          string
	Repository    *helmv1beta1.HelmChartRepository
	ChartVersions repo.ChartVersions
}
