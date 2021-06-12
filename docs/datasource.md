---
layout: default 
title: Datasource 
nav_order: 3
---

## Datasource

### Usage of datasource `rancherk3d_cluster_list`{: .fs-3 }<br><br>
This data_source helps in retrieving information of all available k3d cluster in a selected runtime or selected set of
clusters.

##### Sample `rancherk3d_cluster_list`{: .fs-3 } config:

```terraform
data "rancherk3d_cluster_list" "k3s-default" {
  clusters = [
    "k3s-default"]
}
```

#### Argument Reference
{: .fw-700 }

* `clusters`{: .fs-3 }[`list`{: .fs-3 }] - The list of k3d cluster names of which the information to be retrieved.
* `all`{: .fs-3 }[`boolean`{: .fs-3 }] - If enabled retrieves data of all available clusters in the selected runtime.<br><br>

#### Attributes Reference 
{: .fw-700 }

* `clusters_list`{: .fs-3 }[`list`{: .fs-3 }] - The information of the retrieved clusters.
    * `name`: Cluster name that was retrieved.
    * `nodes`: List of nodes present in the cluster.
    * `network`: Network associated with the cluster.
    * `cluster_token`: Token of the cluster.
    * `servers_count`: Numbers of server nodes present in the cluster.
    * `servers_running`: Numbers of servers running.
    * `agents_count`: Numbers of agents in the cluster.
    * `agents_running`: Number of agents running in the cluster.
    * `has_loadbalancer`: Attribute that notifies the presence of loadbalancer in the cluster.
    * `image_volume`: Volume to import images.

---
### Usage of datasource `rancherk3d_node_list`{: .fs-3 }<br><br>
This data_source helps in retrieving information of a selected k3d node/nodes, or of all nodes in a selected cluster,
else all available nodes in a selected runtime.

##### Sample `rancherk3d_node_list`{: .fs-3 } config:

```terraform
data "rancherk3d_node_list" "k3s_default_server" {
  cluster = "k3s-default"
  nodes = [
    "k3d-k3s-default-server-0",
    "k3d-k3s-default-serverlb"]
}
```

#### Argument Reference
{: .fw-700 }

* `nodes`{: .fs-3 }[`list`{: .fs-3 }] - The list of nodes of which information to be retrieved. If passed empty, details of all nodes from a
  selected cluster would be retrieved.
* `cluster`{: .fs-3 }[`string`{: .fs-3 }] - Name of the cluster from which the node details to be retrieved.
* `all`{: .fs-3 }[`boolean`{: .fs-3 }] - If enabled retrieves information of all nodes from a selected cluster or from runtime depending on
  rest arguments.<br><br>

#### Attributes Reference
{: .fw-700 }

* `node_list`{: .fs-3 }[`list`{: .fs-3 }] - List of node information retrieved. This list would contain below attributes.
    * `name`: Name of the retrieved node.
    * `role`: Role of node retrieved.
    * `cluster`: Cluster to which the node belongs.
    * `state`: Current state of node (`running`/`exited`).
    * `created`: Creation time-stamp of node.
    * `memory`: Memory limit imposed on the node.
    * `volumes`: List of volumes associated with the nodes.
    * `networks`: List of networks associated with the nodes.
    * `env`: List of environment variables set in the node

---
### Usage of datasource `rancherk3d_kubeconfig`{: .fs-3 }<br><br>
This data_source helps in retrieving `kube-config` from a selected cluster or from all clusters in a selected runtime.

Retrieved `kube-config` will be base64 encoded by default, can be skipped by enabling argument.

##### Sample `rancherk3d_kubeconfig`{: .fs-3 } config:

```terraform
data "rancherk3d_kubeconfig" "k3s-default" {
  clusters = [
    "k3s-default"]
}
```

#### Argument Reference
{: .fw-700 }

* `clusters`{: .fs-3 }[`list`{: .fs-3 }] - List of cluster from which the `kube-config` to be retrieved.
* `all`{: .fs-3 }[`boolean`{: .fs-3 }] - If enabled `kube-config` from all available cluster would be retrieved in the selected runtime.
* `not_encoded`{: .fs-3 }[`boolean`{: .fs-3 }] - If enabled retrieved kube-config would not be base64 encoded.<br><br>

#### Attributes Reference
{: .fw-700 }

* `kube_config`{: .fs-3 }[`map`{: .fs-3 }] - Base64 encoded kube-config from a selected cluster.

---
### Usage of datasource `rancherk3d_registry_list`{: .fs-3 }<br><br>
This data_source helps in retrieving information of registries present in the environment.

##### Sample `rancherk3d_registry_list`{: .fs-3 } config:

```terraform
data "rancherk3d_registry_list" "registry-1" {
  registries = [
    "k3d-k3s-sample-registr"]
  cluster = "k3s-default"
}
```

#### Argument Reference
{: .fw-700 }

* `registries`{: .fs-3 }[`list`{: .fs-3 }] - list of registries to be retrieved from the cluster selected.
* `cluster`{: .fs-3 }[`string`{: .fs-3 }] - if enabled fetches all the registries, if cluster is selected then all registries connected to
  it.
* `all`{: .fs-3 }[`boolean`{: .fs-3 }] - name of the cluster of which registries to be retrieved.<br><br>
#### Attributes Reference
{: .fw-700 }

* `registries_list`{: .fs-3 }[`list`{: .fs-3 }] - list of registries retrieved.
    * `name`: Name of the registry retrieved.
    * `role`: Role of registry retrieved.
    * `image`: Image used for registry creation.
    * `cluster`: Cluster to which the registry belongs.
    * `state`: Current state of registry node.
    * `created`: Creation time-stamp of registry.
    * `networks`: Networks associated with the registry node.
    * `env`: Environment variables set in the registry node.
    * `port_mappings`: Port mappings of the registry.