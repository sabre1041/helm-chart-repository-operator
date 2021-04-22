package utils

import (
	"crypto/tls"
	"fmt"
	"time"

	redhatcopv1alpha1 "github.com/redhat-cop/helm-chart-repository-operator/api/v1alpha1"
	"github.com/redhat-cop/helm-chart-repository-operator/pkg/types"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/repo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	repositoryLabelKey = "helm-chart-repository-operator.redhat-cop.io/repository"
)

func DefaultCiphers() []uint16 {
	// HTTP/2 mandates TLS 1.2 or higher with an AEAD cipher
	// suite (GCM, Poly1305) and ephemeral key exchange (ECDHE, DHE) for
	// perfect forward secrecy. Servers may provide additional cipher
	// suites for backwards compatibility with HTTP/1.1 clients.
	// See RFC7540, section 9.2 (Use of TLS Features) and Appendix A
	// (TLS 1.2 Cipher Suite Black List).
	return []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, // required by http/2
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256, // forbidden by http/2, not flagged by http2isBadCipher() in go1.8
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,   // forbidden by http/2, not flagged by http2isBadCipher() in go1.8
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,    // forbidden by http/2
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,    // forbidden by http/2
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,      // forbidden by http/2
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,      // forbidden by http/2
		tls.TLS_RSA_WITH_AES_128_GCM_SHA256,         // forbidden by http/2
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,         // forbidden by http/2
		// the next one is in the intermediate suite, but go1.8 http2isBadCipher() complains when it is included at the recommended index
		// because it comes after ciphers forbidden by the http/2 spec
		// tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
		// tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA, // forbidden by http/2, disabled to mitigate SWEET32 attack
		// tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,       // forbidden by http/2, disabled to mitigate SWEET32 attack
		tls.TLS_RSA_WITH_AES_128_CBC_SHA, // forbidden by http/2
		tls.TLS_RSA_WITH_AES_256_CBC_SHA, // forbidden by http/2
	}
}

func DefaultTLSVersion() uint16 {
	// Can't use SSLv3 because of POODLE and BEAST
	// Can't use TLSv1.0 because of POODLE and BEAST using CBC cipher
	// Can't use TLSv1.1 because of RC4 cipher usage
	return tls.VersionTLS12
}

// SecureTLSConfig enforces the default minimum security settings for the cluster.
func SecureTLSConfig(config *tls.Config) *tls.Config {
	if config.MinVersion == 0 {
		config.MinVersion = DefaultTLSVersion()
	}

	config.PreferServerCipherSuites = true
	if len(config.CipherSuites) == 0 {
		config.CipherSuites = DefaultCiphers()
	}
	return config
}

func MapToHelmChart(helmChartEntry *types.HelmChartEntry) (*redhatcopv1alpha1.HelmChart, error) {

	helmChart := &redhatcopv1alpha1.HelmChart{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HelmChart",
			APIVersion: redhatcopv1alpha1.GroupVersion.String(),
		},
	}

	helmChart.Name = fmt.Sprintf("%s.%s", helmChartEntry.Repository.Name, helmChartEntry.Name)

	helmChart.SetLabels(map[string]string{
		repositoryLabelKey: helmChartEntry.Repository.Name,
	})

	if helmChartEntry.Repository.Spec.DisplayName != "" {
		helmChart.Spec.RepositoryDisplayName = helmChartEntry.Repository.Spec.DisplayName
	}

	helmChart.Spec.RepositoryName = helmChartEntry.Repository.Name
	helmChart.Spec.Name = helmChartEntry.Name

	chartVersions := []redhatcopv1alpha1.HelmChartVersion{}

	if helmChartEntry.ChartVersions != nil {
		for _, chartVersion := range helmChartEntry.ChartVersions {

			if chartVersion.Metadata != nil && chartVersion.Metadata.KubeVersion != "" && helmChartEntry.ServerVersion != "" {
				if !chartutil.IsCompatibleRange(chartVersion.Metadata.KubeVersion, helmChartEntry.ServerVersion) {
					continue
				}
			}

			helmChartVersion, err := mapToHelmChartVersion(chartVersion)

			if err != nil {
				return nil, err
			}

			chartVersions = append(chartVersions, *helmChartVersion)
		}
	}

	helmChart.Spec.Versions = chartVersions

	return helmChart, nil
}

func mapToHelmChartVersion(chartVersion *repo.ChartVersion) (*redhatcopv1alpha1.HelmChartVersion, error) {

	helmChartVersion := &redhatcopv1alpha1.HelmChartVersion{}
	helmChartVersion.ApiVersion = chartVersion.APIVersion
	helmChartVersion.AppVersion = chartVersion.AppVersion

	if (chartVersion.Created != time.Time{}) {
		helmChartVersion.Created = &metav1.Time{Time: chartVersion.Created}
	}

	if chartVersion.Maintainers != nil {

		maintainers := []redhatcopv1alpha1.HelmChartMaintainer{}

		for _, maintainer := range chartVersion.Maintainers {
			maintainers = append(maintainers, redhatcopv1alpha1.HelmChartMaintainer{
				Name:  maintainer.Name,
				Email: maintainer.Email,
				URL:   maintainer.URL,
			})
		}

		helmChartVersion.Maintainers = &maintainers

	}

	if chartVersion.Dependencies != nil {
		dependencies := []redhatcopv1alpha1.HelmChartDependency{}

		for _, dependency := range chartVersion.Dependencies {
			dependencies = append(dependencies, redhatcopv1alpha1.HelmChartDependency{
				Alias:      dependency.Alias,
				Condition:  dependency.Condition,
				Enabled:    dependency.Enabled,
				Name:       dependency.Name,
				Repository: dependency.Repository,
				Tags:       dependency.Tags,
				Version:    dependency.Version,
			})
		}

		helmChartVersion.Dependencies = &dependencies
	}

	helmChartVersion.Description = chartVersion.Description
	helmChartVersion.Digest = chartVersion.Digest
	helmChartVersion.Home = chartVersion.Home
	helmChartVersion.Icon = chartVersion.Icon
	helmChartVersion.Keywords = chartVersion.Keywords
	helmChartVersion.KubeVersion = chartVersion.KubeVersion
	helmChartVersion.Type = chartVersion.Type
	helmChartVersion.URLs = chartVersion.URLs
	helmChartVersion.Version = chartVersion.Version

	return helmChartVersion, nil

}
