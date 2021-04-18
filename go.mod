module github.com/redhat-cop/helm-chart-repository-operator

go 1.15

require (
	github.com/go-logr/logr v0.3.0
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/openshift/api v0.0.0-20210202165416-a9e731090f5e
	github.com/redhat-cop/operator-utils v1.1.0
	helm.sh/helm/v3 v3.5.0
	k8s.io/api v0.20.1 //ct
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v0.20.1
	sigs.k8s.io/controller-runtime v0.7.0
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309
	github.com/opencontainers/runc => github.com/opencontainers/runc v1.0.0-rc8.0.20190926150303-84373aaa560b
)
