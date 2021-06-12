---
layout: default
title: terraform-provider-rancherk3d
nav_order: 1
---
# terraform-provider-rancherk3d

## Backdrop
* Main reason why the provider was created is to simplify handling [`kubernetes`](https://kubernetes.io/) resources creation.
* Especially while working on `kubernetes native` applications, local `kube` cluster provisioned with `k3d` comes in handy. Manging these would be easier with this terraform provider.
* Yes [`k3d`](https://k3d.io/usage/commands/k3d/) cli does the work, but bringing in a layer of terraform will add more control over the resource `planning`/`creation`/`deletion`.
* This provider `terraform-provider-rancherk3d` covers almost every task that one can accomplish with `k3d cli` with no `k3d` installed.

## Prerequisites
* [`Terraform`](https://www.terraform.io/downloads.html) v0.13.x [`tested`]
* [`Docker`](https://www.docker.com/)

## Provider

It is expected that few configurations to be passed while configuring provider.

The configured values would be used by all `data_source` and `resource` available.

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
{: .fw-700 }

* `kubernetes_version`{: .fs-3 }: Rancher k3s kubernetes version to be used across all resource creation, including(cluster, node, registries).
* `k3d_api_version`{: .fs-3 }: Rancher k3d api version to be used while creating a cluster.
* `registry`{: .fs-3 }: Registries to be used, can be used as defaults and can be overridden during resource creation.
* `kind`{: .fs-3 }: Defines the kind of config file to be used, at the moment supports just `Simple`, shall support more once k3d starts supporting.
* `runtime`{: .fs-3 }: Runtime to be used, at the moment supports just `Docker`, shall support more once k3d starts supporting.
* `version`{: .fs-3 }: `terraform-provider-rancherk3d` version to be used.

**Note**: Table on all supported features can be found [here](https://github.com/nikhilsbhat/terraform-provider-rancherk3d#features-supported-by-the-provider-at-the-moment).