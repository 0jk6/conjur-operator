# conjur-operator
A custom kubernetes operator that pulls secrets from a CyberArk safe and stores them as a native k8s secret object.

## Description.
To install this operator in your cluster, you must follow the getting started steps available down below.
To run this locally, clone this repo and run the following commands.

`make manifests` - generates manifest files.
`make install` - installs the CRD, RBAC, etc in your cluster
`make run` - runs the controller on your local machine, which will reconcile the Conjur objects.

Assuming that you have followed all the steps mentioned above, you can apply the following yaml files in your namespace.

Apply the following yaml to create a secret that stores the Conjur API key in base64 encoded format
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: conjur-api-key
type: Opaque
stringData:
  apikey: "---put your base64 encoded conjur api key here---"
```

Apply the following yaml to create the Conjur object
```yaml
apiVersion: api.0jk6.github.io/v1alpha1
kind: Conjur
metadata:
  labels:
    app.kubernetes.io/name: conjur-operator
    app.kubernetes.io/managed-by: kustomize
  name: conjur-sample
spec:
  refreshInterval: 60

  apiKeyFromSecret: "conjur-api-key" #another secret which has base64 encoded conjur api key

  conjurHost: "your conjur host url"
  conjurAcct: "your conjur account" #eg: cyberarkprd
  hostname: "your host name" #eg: host/apps/some-name

  data:
    dev-database: #a secret will be created in your namespace with this name storing the value of the secretIdentfier
      secretIdentifier: "your secret identifier"
    prd-db: #another secret
      secretIdentifier: "your secret identifier"
```

## Getting Started

### Prerequisites
- go version v1.21.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/conjur-operator:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/conjur-operator:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following are the steps to build the installer and distribute this project to users.

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/conjur-operator:tag
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

2. Using the installer

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/conjur-operator/<tag or branch>/dist/install.yaml
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2024 0jk6.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

