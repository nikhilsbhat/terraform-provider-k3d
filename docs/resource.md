---
layout: default 
title: Resource 
nav_order: 2
---

## Resource

### Usage of resource `rancherk3d_cluster_action`{: .fs-3 }<br><br>
This resource helps in `starting`/`stopping` clusters from the selected runtime.

##### Sample `rancherk3d_cluster_action`{: .fs-3 } config:

```terraform
resource rancherk3d_cluster_action "start-k3s-cluster" {
  clusters = [
    "k3s-default"]
  start = true
}
```

#### Argument Reference
{: .fw-700 }

* `clusters`{: .fs-3 }[`list`{: .fs-3 }] - List of k3s clusters that has to be `started`/`stopped`.
* `all`{: .fs-3 }[`boolean`{: .fs-3 }] - If enabled selected clusters would be started/stopped.
* `start`{: .fs-3 }[`boolean`{: .fs-3 }] - If enabled it starts a stopped cluster.
* `stop`{: .fs-3 }[`boolean`{: .fs-3 }] - If enabled it stops a running cluster.<br><br>

#### Attributes Reference
{: .fw-700 }

* `state`{: .fs-3 }[`string`{: .fs-3 }] - Last state of selected clusters.
* `status`{: .fs-3 }[`list`{: .fs-3 }] - Updated status of clusters.
    * `name`: Cluster name that was fetched.
    * `nodes`: List of nodes present in the cluster.
    * `network`: Network associated with the cluster.
    * `cluster_token`: Token of the cluster.
    * `servers_count`: Number of server nodes present in the cluster.
    * `servers_running`: Number of server nodes running.
    * `agents_count`: Number of agents in the cluster.
    * `agents_running`: Number of agents running in the cluster.
    * `has_loadbalancer`: Attribute that notifies the presence of loadbalancer in the cluster.
    * `image_volume`: Volume to import images.

---  
### Usage of resource `rancherk3d_node_action`{: .fs-3 }<br><br>
This resource helps in `starting`/`stopping` nodes from the selected clusters.

##### Sample `rancherk3d_node_action`{: .fs-3 } config:

```terraform
resource "rancherk3d_node_action" "k3s-default-nodes" {
  nodes = [
    "k3d-k3s-default-serverlb",
    "k3d-k3s-default-server-0"]
  cluster = "k3s-default"
  start = true
}
```

#### Argument Reference
{: .fw-700 }

* `nodes`{: .fs-3 }[`list`{: .fs-3 }] - List of k3s nodes that has to be `started`/`stopped`.
* `nodes`{: .fs-3 }[`list`{: .fs-3 }] - List of k3s nodes that has to be `started`/`stopped`.
* `cluster`{: .fs-3 }[`string`{: .fs-3 }] - Name of the cluster from which nodes to be acted upon.
* `all`{: .fs-3 }[`boolean`{: .fs-3 }] - If enabled selected nodes from a selected cluster would be started/stopped.
* `start`{: .fs-3 }[`boolean`{: .fs-3 }] - If enabled it starts a stopped nodes.
* `stop`{: .fs-3 }[`boolean`{: .fs-3 }] - If enabled it stops a running node.<br><br>

#### Attributes Reference
{: .fw-700 }

* `status`{: .fs-3 }[`list`{: .fs-3 }] - Updated status of nodes.
    * `node`: Node of which the current status is updated with.
    * `role`: Role of updated node.
    * `state`: Current state of the node specified.
    * `cluster`: Name of the cluster of to which node belongs.

---
### Usage of resource `rancherk3d_node_create`{: .fs-3 }<br><br>
This resource helps in creation of k3d nodes with preferred configurations.

##### Sample `rancherk3d_node_create`{: .fs-3 } config:

```terraform
resource "rancherk3d_node_create" "test-nodes" {
  name = "test-node-terraform"
  cluster = "k3s-default"
  role = "agent"
  replicas = 2
  volume = "8g"
}
```

#### Argument Reference
{: .fw-700 }

* `name`{: .fs-3 }[`string`{: .fs-3 }] - Name the nodes to be created (index would be used to dynamically compute the names for nodes).
* `cluster`{: .fs-3 }[`string`{: .fs-3 }] - Name of the cluster to which these nodes to be connected with.
* `image`{: .fs-3 }[`string`{: .fs-3 }] - Image to be used for nodes creation defaults to image declared in the provider.
* `role`{: .fs-3 }[`string`{: .fs-3 }] - Role to be assigned to the node(agent).
* `replicas`{: .fs-3 }[`int`{: .fs-3 }] - Total number of nodes to be created.
* `memory`{: .fs-3 }[`string`{: .fs-3 }] - Memory limit to be imposed on the node.
* `wait`{: .fs-3 }[`boolean`{: .fs-3 }] - If enabled waits for the nodes to be ready before returning.
* `timeout`{: .fs-3 }[`int`{: .fs-3 }] - Maximum waiting time for before canceling/returning in minutes.<br><br>

#### Attributes Reference
{: .fw-700 }

* `node_list`{: .fs-3 }[`list`{: .fs-3 }] - List of node information created. This list would contain below attributes.
    * `name`: Name of the created node.
    * `role`: Role of node created.
    * `cluster`: Cluster to which the node belongs.
    * `state`: Current state of node (`running`/`exited`).
    * `created`: Creation time-stamp of node.
    * `memory`: Memory limit imposed on the node.
    * `volumes`: List of volumes associated with the nodes.
    * `networks`: List of networks associated with the nodes.
    * `env`: List of environment variables set in the node

---
### Usage of resource `rancherk3d_create_registry`{: .fs-3 }<br><br>
This resource helps in creation of k3d registries with preferred configurations and associating it with the selected
cluster.

##### Sample `rancherk3d_create_registry`{: .fs-3 } config:

```terraform
resource "rancherk3d_create_registry" "registry" {
  name = "k3s-registry"
  cluster = "k3s-default"
  protocol = "http"
  host = "test-registry.com"
}
```

#### Argument Reference
{: .fw-700 }

* `name`{: .fs-3 }[`string`{: .fs-3 }] - Name the registry node to be created.
* `image`{: .fs-3 }[`string`{: .fs-3 }] - Image to be used for creation of registry(defaults to docker.io/library/registry:2).
* `cluster`{: .fs-3 }[`string`{: .fs-3 }] - Cluster to which the registry to be associated with.
* `protocol`{: .fs-3 }[`string`{: .fs-3 }] - Protocol to be used while running registry (defaults to http).
* `host`{: .fs-3 }[`string`{: .fs-3 }] - Host name to be assigned to the registry that would be created (defaults to name of registry).
* `config_file`{: .fs-3 }[`string`{: .fs-3 }] - Config file to be used for configuring registry.
* `expose`{: .fs-3 }[`map`{: .fs-3 }] - Host to port mapping.
* `use_proxy`{: .fs-3 }[`boolean`{: .fs-3 }] - If enabled proxy configuration provided at 'proxy' would be used for configuring registry.
* `proxy`{: .fs-3 }[`map`{: .fs-3 }] - Proxy configurations to be used while configuring registry if enabled.<br><br>

#### Attributes Reference
{: .fw-700 }

* `registries_list`{: .fs-3 }[`list`{: .fs-3 }]: List of registries information those were created, This list would contain below
  attributes.
    * `name`: Name of the registry retrieved.
    * `role`: Role of registry retrieved.
    * `image`: Image used for registry creation.
    * `cluster`: Cluster to which the registry belongs.
    * `state`: Current state of registry node.
    * `created`: Creation time-stamp of registry.
    * `networks`: Networks associated with the registry node.
    * `env`: Environment variables set in the registry node.
    * `port_mappings`: Port mappings of the registry.

---
### Usage of resource `rancherk3d_connect_registry`{: .fs-3 }<br><br>
This resource helps in coupling/decoupling registry from the selected cluster.

##### Sample `rancherk3d_connect_registry`{: .fs-3 } config:

```terraform
resource "rancherk3d_connect_registry" "k3s-registry-1" {
  registries = [
    "k3s-registry-1"]
  cluster = "k3s-default"
  connect = true
}
```

#### Argument Reference
{: .fw-700 }

* `registries`{: .fs-3 }[`list`{: .fs-3 }] - List of registries to be connected/disconnected from the selected cluster.
* `cluster`{: .fs-3 }[`string`{: .fs-3 }] - Cluster to which registries to be associated with.
* `connect`{: .fs-3 }[`boolean`{: .fs-3 }] - Enable this flag if registries to be connected with specified cluster, when disabled it
  disconnects the registry from cluster.<br><br>

#### Attributes Reference
{: .fw-700 }

* `status`{: .fs-3 }[`list`{: .fs-3 }] - Updated status of registry. This list would contain below attributes.
    * `registry`{: .fs-3 }[`string`{: .fs-3 }]: Name of the registry.
    * `cluster`{: .fs-3 }[`string`{: .fs-3 }]: Cluster to which the registry is either connected/disconnected.
    * `state`{: .fs-3 }[`string`{: .fs-3 }]: Updated state of registry, `connected`/`disconnected`.

---
### Usage of resource `rancherk3d_load_image`{: .fs-3 }<br><br>
This resource helps in loading a list of images to all selected clusters.

##### Sample `rancherk3d_load_image`{: .fs-3 } config:

```terraform
resource "rancherk3d_load_image" "k3s-default" {
  images = [
    "basnik/terragen:latest"]
  cluster = "k3s-default"
  keep_tarball = false
}
```

#### Argument Reference
{: .fw-700 }

* `images`{: .fs-3 }[`list`{: .fs-3 }] - List of images to be imported to the existing cluster.
* `cluster`{: .fs-3 }[`string`{: .fs-3 }] - Name of the existing cluster to which the images has to be imported to.
* `all`: [`boolean`{: .fs-3 }] - If enabled loads images to all available clusters in the selected runtime.
* `keep_tarball`{: .fs-3 }[`boolean`{: .fs-3 }] - Enable to keep the tarball of the loaded images locally.<br><br>

#### Attributes Reference
{: .fw-700 }

* `images_stored`{: .fs-3 }[`list`{: .fs-3 }] - list of images loaded to the cluster.
    * `cluster`{: .fs-3 }[`string`{: .fs-3 }]: Cluster to which the below images are stored.
    * `images`{: .fs-3 }[`list`{: .fs-3 }]: List of images and its tarball stored, if in case keep_tarball is enabled.

