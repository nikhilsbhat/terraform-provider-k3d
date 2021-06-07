# terraform-provider-rancherk3d

## Backdrop
* Main reason why the provider was created is to simplify handling `kubernetes` resources creation.
* Especially while working on `kubernetes native` applications, local cluster comes in handy. Manging these would be easier with this terraform provider.
* Yes with `k3d cli` does the work, but bringing terraform layer will add more control over the resource `planning`/`creation`/`deletion`.
* This provider `terraform-provider-rancherk3d` covers almost every task that one can accomplish with k3d cli with no k3d installed.<br><br>

## Provider

This expects few configurations to be passed while configuring provider.

The configured values would be used by all data_source and resource available.

##### Sample Provider config:

```terraform
provider "rancherk3d" {
  kubernetes_version = "1.20.2-k3s1"
  k3d_api_version    = "k3d.io/v1alpha2"
  registry           = "rancher/k3s"
  kind               = "Simple"
  runtime            = "docker"
  version            = "0.1.0"
}
```

#### Argument Reference
* `kubernetes_version`: Rancher k3s kubernetes version to be used across all resource creation, including(cluster, node, registries)<br><br>
* `k3d_api_version`: Rancher k3d api version to be used while creating a cluster.<br><br>
* `registry`: Registries to be used, can be used as defaults and can be overridden during resource creation.<br><br>
* `kind`: Defines the kind of config file to be used, at the moment supports just `Simple`, shall support more once k3d starts supporting.<br><br>
* `runtime`: Runtime to be used, at the moment supports just `Docker`, shall support more once k3d starts supporting.<br><br>
* `version`: `terraform-provider-rancherk3d` version to be used.<br><br>