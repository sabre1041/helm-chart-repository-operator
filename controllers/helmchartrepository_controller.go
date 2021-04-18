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

package controllers

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"sigs.k8s.io/yaml"

	"github.com/go-logr/logr"
	configv1 "github.com/openshift/api/config/v1"
	helmv1beta1 "github.com/openshift/api/helm/v1beta1"
	"github.com/redhat-cop/helm-chart-repository-operator/pkg/types"
	"github.com/redhat-cop/helm-chart-repository-operator/pkg/utils"

	"github.com/redhat-cop/operator-utils/pkg/util"
	"helm.sh/helm/v3/pkg/repo"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	kubeclock "k8s.io/apimachinery/pkg/util/clock"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	configNamespace = "openshift-config"
)

var clock kubeclock.Clock = &kubeclock.RealClock{}

// HelmChartRepositoryReconciler reconciles a HelmChartRepository object
type HelmChartRepositoryReconciler struct {
	util.ReconcilerBase
	Log             logr.Logger
	ReconcilePeriod int
}

//+kubebuilder:rbac:groups=redhatcop.redhat.io,resources=helmchartrepositories,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=redhatcop.redhat.io,resources=helmchartrepositories/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=redhatcop.redhat.io,resources=helmchartrepositories/finalizers,verbs=update
//+kubebuilder:rbac:groups="",namespace=openshift-config,resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",namespace=openshift-config,resources=configmaps,verbs=get;list;watch

func (r *HelmChartRepositoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("helmchartrepository", req.NamespacedName)

	instance := &helmv1beta1.HelmChartRepository{}
	err := r.GetClient().Get(ctx, req.NamespacedName, instance)

	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	r.Log.Info("Reconciling Helm Chart Repository", "Name", instance.Name)

	if !instance.Spec.Disabled {

		var indexFile repo.IndexFile

		httpClient, err := r.getHttpClient(ctx, instance)
		if err != nil {
			return reconcile.Result{}, err
		}

		repositoryURL, err := url.Parse(instance.Spec.ConnectionConfig.URL)
		if err != nil {
			return reconcile.Result{}, errors.New(fmt.Sprintf("Unable to parse repository URL %v", instance.Spec.ConnectionConfig.URL))
		}

		indexURL := repositoryURL.String()
		if !strings.HasSuffix(indexURL, "/index.yaml") {
			indexURL += "/index.yaml"
		}
		resp, err := httpClient.Get(indexURL)
		if err != nil {
			return reconcile.Result{}, err
		}
		if resp.StatusCode != 200 {
			return reconcile.Result{}, errors.New(fmt.Sprintf("Response for %v returned %v with status code %v", indexURL, resp, resp.StatusCode))
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return reconcile.Result{}, err
		}
		err = yaml.Unmarshal(body, &indexFile)
		if err != nil {
			return reconcile.Result{}, err
		}
		for _, chartVersions := range indexFile.Entries {
			for _, chartVersion := range chartVersions {
				for i, url := range chartVersion.URLs {
					chartVersion.URLs[i], err = repo.ResolveReferenceURL(indexURL, url)
					if err != nil {
						r.Log.Error(err, "Error resolving chart url", instance.Name)
					}
				}
			}
		}

		// Sort Entries
		indexFile.SortEntries()

		for chartName, versions := range indexFile.Entries {

			helmChart, err := utils.MapToHelmChart(&types.HelmChartEntry{Name: chartName, Repository: instance, ChartVersions: versions})

			if err != nil {
				r.Log.Error(err, "Failed to map to Helm Chart")
				return reconcile.Result{}, err
			}

			err = r.CreateOrUpdateResource(ctx, instance, "", helmChart)

			if err != nil {
				r.Log.Error(err, "Failed to Update Chart", "Name", helmChart.Name)
				return reconcile.Result{}, err
			}

			helmChart.Status.LastUpdateTimestamp = &metav1.Time{Time: clock.Now()}

			err = r.GetClient().Status().Update(ctx, helmChart)

			if err != nil {
				r.Log.Error(err, "Failed to Update Chart Status", "Name", helmChart.Name)
				return reconcile.Result{}, err
			}

		}

	} else {
		r.Log.Info("Skipping Disabled Chart Repository", "Name", instance.Name)
	}

	return ctrl.Result{Requeue: true, RequeueAfter: time.Second * time.Duration(r.ReconcilePeriod)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HelmChartRepositoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&helmv1beta1.HelmChartRepository{}).
		Complete(r)
}

func (r *HelmChartRepositoryReconciler) getHttpClient(ctx context.Context, helmChartRepository *helmv1beta1.HelmChartRepository) (*http.Client, error) {

	var err error

	var rootCAs *x509.CertPool

	if (helmChartRepository.Spec.ConnectionConfig.CA != configv1.ConfigMapNameReference{}) {
		caName := helmChartRepository.Spec.ConnectionConfig.CA.Name

		configMap := &corev1.ConfigMap{}
		err = r.GetClient().Get(ctx, k8stypes.NamespacedName{Name: caName, Namespace: configNamespace}, configMap)

		if err != nil {
			r.Log.Error(err, "Unable to access ConfigMap from OpenShift Config Namespace", "Name", helmChartRepository.Spec.ConnectionConfig.CA.Name)
		}
		caBundleKey := "ca-bundle.crt"
		caCert, found := configMap.Data[caBundleKey]

		if !found {
			return nil, errors.New(fmt.Sprintf("Failed to find %s key in configmap %s", caBundleKey, caName))
		}

		if caCert != "" {
			rootCAs = x509.NewCertPool()
			if ok := rootCAs.AppendCertsFromPEM([]byte(caCert)); !ok {
				return nil, errors.New("Failed to append caCert")
			}
		}
	}

	if rootCAs == nil {
		rootCAs, err = x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
	}

	tlsClientConfig := utils.SecureTLSConfig(&tls.Config{
		RootCAs: rootCAs,
	})

	if (helmChartRepository.Spec.ConnectionConfig.TLSClientConfig != configv1.SecretNameReference{}) {

		secretName := helmChartRepository.Spec.ConnectionConfig.TLSClientConfig.Name

		secret := &corev1.Secret{}
		err := r.GetClient().Get(ctx, k8stypes.NamespacedName{Name: secretName, Namespace: configNamespace}, secret)

		if err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to GET secret %s reason %v", secretName, err))
		}
		tlsCertSecretKey := "tls.crt"
		tlsCert, ok := secret.Data[tlsCertSecretKey]
		if !ok {
			return nil, errors.New(fmt.Sprintf("Failed to find %s key in secret %s", tlsCertSecretKey, secretName))
		}
		tlsSecretKey := "tls.key"
		tlsKey, ok := secret.Data[tlsSecretKey]
		if !ok {
			return nil, errors.New(fmt.Sprintf("Failed to find %s key in secret %s", tlsSecretKey, secretName))
		}
		if tlsKey != nil && tlsCert != nil {
			cert, err := tls.X509KeyPair(tlsCert, tlsKey)
			if err != nil {
				return nil, err
			}
			tlsClientConfig.Certificates = []tls.Certificate{cert}
		}

	}

	tr := &http.Transport{
		TLSClientConfig: tlsClientConfig,
		Proxy:           http.ProxyFromEnvironment,
	}

	return &http.Client{Transport: tr}, nil

}
