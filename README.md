# Helm Chart Repository Operator

Operator that exposes [Helm](https://helm.sh) charts contained within Repositories declared on an OpenShift environment.

## Installation

Execute the following command to install and deploy the operator:

```shell
make install
make deploy IMG=quay.io/ablock/helm-chart-repository-operator:latest
```

Once the operator has had a chance to reconcile, the list of Helm Charts will be available as shown below:

```shell
oc get helmcharts

NAME                                               REPOSITORY         NAME                              LATEST VERSION
redhat-helm-repo.ibm-b2bi-prod                     redhat-helm-repo   ibm-b2bi-prod                     2.0.0
redhat-helm-repo.ibm-cpq-prod                      redhat-helm-repo   ibm-cpq-prod                      4.0.1
redhat-helm-repo.ibm-mongodb-enterprise-helm       redhat-helm-repo   ibm-mongodb-enterprise-helm       0.1.0
redhat-helm-repo.ibm-object-storage-plugin         redhat-helm-repo   ibm-object-storage-plugin         2.0.7
redhat-helm-repo.ibm-oms-ent-prod                  redhat-helm-repo   ibm-oms-ent-prod                  6.0.0
redhat-helm-repo.ibm-oms-pro-prod                  redhat-helm-repo   ibm-oms-pro-prod                  6.0.0
redhat-helm-repo.ibm-operator-catalog-enablement   redhat-helm-repo   ibm-operator-catalog-enablement   1.1.0
redhat-helm-repo.ibm-sfg-prod                      redhat-helm-repo   ibm-sfg-prod                      2.0.0
redhat-helm-repo.nodejs                            redhat-helm-repo   nodejs                            0.0.1
redhat-helm-repo.nodejs-ex-k                       redhat-helm-repo   nodejs-ex-k                       0.2.1
```
