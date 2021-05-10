# Document for `terraform-provider-rancherk3d`

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

## Data Source

### usage of datasource `rancherk3d_cluster_list`
___

This data_source helps in retrieving information of all available k3d cluster in a selected runtime or selected set of clusters.
##### Sample `rancherk3d_cluster_list` config:

```terraform
data "rancherk3d_cluster_list" "k3s-default" {
  clusters = ["k3s-default"]
}
```
#### Argument Reference
* `clusters`[`list`] - The list of k3d cluster names of which the information to be retrieved.<br><br>
* `all`[`boolean`] - If enabled retrieves data of all available clusters in the selected runtime.
#### Attributes Reference
* `clusters_list`[`list`] - The information of the retrieved clusters.<br><br>
    * `name`: Cluster name that was retrieved.<br><br>
    * `nodes`: List of nodes present in the cluster.<br><br>
    * `network`: Network associated with the cluster.<br><br>
    * `cluster_token`: Token of the cluster.<br><br>
    * `servers_count`: Numbers of server nodes present in the cluster.<br><br>
    * `servers_running`: Numbers of servers running.<br><br>
    * `agents_count`: Numbers of agents in the cluster.<br><br>
    * `agents_running`: Number of agents running in the cluster.<br><br>
    * `has_loadbalancer`: Attribute that notifies the presence of loadbalancer in the cluster.<br><br>
    * `image_volume`: Volume to import images.<br><br>

### usage of datasource `rancherk3d_node_list`
___
This data_source helps in retrieving information of a selected k3d node/nodes, or of all nodes in a selected cluster, else all available nodes in a selected runtime.
##### Sample `rancherk3d_node_list` config:

```terraform
data "rancherk3d_node_list" "k3s_default_server" {
  cluster = "k3s-default"
  nodes   = ["k3d-k3s-default-server-0", "k3d-k3s-default-serverlb"]
}
```
#### Argument Reference
* `nodes`[`list`] - The list of nodes of which information to be retrieved. If passed empty, details of all nodes from a selected cluster would be retrieved. <br><br>
* `cluster`[`string`] - Name of the cluster from which the node details to be retrieved.<br><br>
* `all`[`boolean`] - If enabled retrieves information of all nodes from a selected cluster or from runtime depending on rest arguments.
#### Attributes Reference
* `node_list`[`list`] - List of node information retrieved. This list would contain below attributes.<br><br>
    * `name`: Name of the retrieved node.<br><br>
    * `role`: Role of node retrieved.<br><br>
    * `cluster`: Cluster to which the node belongs.<br><br>
    * `state`: Current state of node (`running`/`exited`).<br><br>
    * `created`: Creation time-stamp of node.<br><br>
    * `memory`: Memory limit imposed on the node.<br><br>
    * `volumes`: List of volumes associated with the nodes.<br><br>
    * `networks`: List of networks associated with the nodes.<br><br>
    * `env`: List of environment variables set in the node<br><br>

### usage of datasource `rancherk3d_kubeconfig`
___
This data_source helps in retrieving `kube-config` from a selected cluster or from all clusters in a selected runtime.

Retrieved `kube-config` will be base64 encoded by default, can be skipped by enabling argument.
##### Sample `rancherk3d_kubeconfig` config:
```terraform
data "rancherk3d_kubeconfig" "k3s-default" {
  clusters = ["k3s-default"]
}
```
#### Argument Reference
* `clusters`[`list`] - List of cluster from which the `kube-config` to be retrieved.<br><br>
* `all`[`boolean`] - If enabled `kube-config` from all available cluster would be retrieved in the selected runtime.<br><br>
* `not_encoded`[`boolean`] - If enabled retrieved kube-config would not be base64 encoded.
#### Attributes Reference
* `kube_config`[`map`] - Base64 encoded kube-config from a selected cluster.<br><br>

### usage of datasource `rancherk3d_registry_list`
___
This data_source helps in retrieving information of registries present in the environment.
##### Sample `rancherk3d_registry_list` config:
```terraform
data "rancherk3d_registry_list" "registry-1" {
  registries = ["k3d-k3s-sample-registr"]
  cluster = "k3s-default"
}
```
#### Argument Reference
* `registries`[`list`] - list of registries to be retrieved from the cluster selected.<br><br>
* `cluster`[`string`] - if enabled fetches all the registries, if cluster is selected then all registries connected to it.<br><br>
* `all`[`boolean`] - name of the cluster of which registries to be retrieved.
#### Attributes Reference
* `registries_list`[`list`] - list of registries retrieved.<br><br>
    * `name`: Name of the registry retrieved.<br><br>
    * `role`: Role of registry retrieved.<br><br>
    * `image`: Image used for registry creation.<br><br>
    * `cluster`: Cluster to which the registry belongs.<br><br>
    * `state`: Current state of registry node.<br><br>
    * `created`: Creation time-stamp of registry.<br><br>
    * `networks`: Networks associated with the registry node.<br><br>
    * `env`: Environment variables set in the registry node.<br><br>
    * `port_mappings`: Port mappings of the registry.

## Resource

### usage of resource `rancherk3d_cluster_action`

---
This resource helps in `starting`/`stopping` clusters from the selected runtime.
##### Sample `rancherk3d_cluster_action` config:

```terraform
resource rancherk3d_cluster_action "start-k3s-cluster" {
  clusters = ["k3s-default"]
  start    = true
}
```
#### Argument Reference
* `clusters`[`list`] - List of k3s clusters that has to be `started`/`stopped`.<br><br>
* `all`[`boolean`] - If enabled selected clusters would be started/stopped.<br><br>
* `start`[`boolean`] - If enabled it starts a stopped cluster.<br><br>
* `stop`[`boolean`] - If enabled it stops a running cluster.
#### Attributes Reference
* `state`[`string`] - Last state of selected clusters.<br><br>
* `status`[`list`] - Updated status of clusters.<br><br>
    * `name`: Cluster name that was fetched.<br><br>
    * `nodes`: List of nodes present in the cluster.<br><br>
    * `network`: Network associated with the cluster.<br><br>
    * `cluster_token`: Token of the cluster.<br><br>
    * `servers_count`: Number of server nodes present in the cluster.<br><br>
    * `servers_running`: Number of server nodes running.<br><br>
    * `agents_count`: Number of agents in the cluster.<br><br>
    * `agents_running`: Number of agents running in the cluster.<br><br>
    * `has_loadbalancer`: Attribute that notifies the presence of loadbalancer in the cluster.<br><br>
    * `image_volume`: Volume to import images.<br><br>

### usage of resource `rancherk3d_node_action`

---
This resource helps in `starting`/`stopping` nodes from the selected clusters.
##### Sample `rancherk3d_node_action` config:

```terraform
resource "rancherk3d_node_action" "k3s-default-nodes" {
  nodes   = ["k3d-k3s-default-serverlb", "k3d-k3s-default-server-0"]
  cluster = "k3s-default"
  start   = true
}
```
#### Argument Reference
* `nodes`[`list`] - List of k3s nodes that has to be `started`/`stopped`.<br><br>
* `cluster`[`string`] - Name of the cluster from which nodes to be acted upon.<br><br>
* `all`[`boolean`] - If enabled selected nodes from a selected cluster would be started/stopped.<br><br>
* `start`[`boolean`] - If enabled it starts a stopped nodes.<br><br>
* `stop`[`boolean`] - If enabled it stops a running node.
#### Attributes Reference
* `status`[`list`] - Updated status of nodes.<br><br>
    * `node`: Node of which the current status is updated with.<br><br>
    * `role`: Role of updated node.<br><br>
    * `state`: Current state of the node specified.<br><br>
    * `cluster`: Name of the cluster of to which node belongs.<br><br>

### usage of resource `rancherk3d_node_create`

---
This resource helps in creation of k3d nodes with preferred configurations.
##### Sample `rancherk3d_node_create` config:

```terraform
resource "rancherk3d_node_create" "test-nodes" {
  name     = "test-node-terraform"
  cluster  = "k3s-default"
  role     = "agent"
  replicas = 2
  volume = "8g"
}
```
#### Argument Reference
*`name`[`string`] - Name the nodes to be created (index would be used to dynamically compute the names for nodes).<br><br>
*`cluster`[`string`] - Name of the cluster to which these nodes to be connected with.<br><br>
*`image`[`string`] - Image to be used for nodes creation defaults to image declared in the provider.<br><br>
*`role`[`string`] - Role to be assigned to the node(agent).<br><br>
*`replicas`[`int`] - Total number of nodes to be created.<br><br>
*`memory`[`string`] - Memory limit to be imposed on the node.<br><br>
*`wait`[`boolean`] - If enabled waits for the nodes to be ready before returning.<br><br>
*`timeout`[`int`] - Maximum waiting time for before canceling/returning in minutes.
#### Attributes Reference
* `node_list`[`list`] - List of node information created. This list would contain below attributes.<br><br>
    * `name`: Name of the created node.<br><br>
    * `role`: Role of node created.<br><br>
    * `cluster`: Cluster to which the node belongs.<br><br>
    * `state`: Current state of node (`running`/`exited`).<br><br>
    * `created`: Creation time-stamp of node.<br><br>
    * `memory`: Memory limit imposed on the node.<br><br>
    * `volumes`: List of volumes associated with the nodes.<br><br>
    * `networks`: List of networks associated with the nodes.<br><br>
    * `env`: List of environment variables set in the node<br><br>

### usage of resource `rancherk3d_create_registry`

---
This resource helps in creation of k3d registries with preferred configurations and associating it with the selected cluster.
##### Sample `rancherk3d_create_registry` config:

```terraform
resource "rancherk3d_create_registry" "registry" {
  name     = "k3s-registry"
  cluster  = "k3s-default"
  protocol = "http"
  host     = "test-registry.com"
}
```
#### Argument Reference
* `name`[`string`] - Name the registry node to be created.<br><br>
* `image`[`string`] - Image to be used for creation of registry(defaults to docker.io/library/registry:2).<br><br>
* `cluster`[`string`] - Cluster to which the registry to be associated with.<br><br>
* `protocol`[`string`] - Protocol to be used while running registry (defaults to http).<br><br>
* `host`[`string`] - Host name to be assigned to the registry that would be created (defaults to name of registry).<br><br>
* `config_file`[`string`] - Config file to be used for configuring registry.<br><br>
* `expose`[`map`] - Host to port mapping.<br><br>
* `use_proxy`[`boolean`] - If enabled proxy configuration provided at 'proxy' would be used for configuring registry.<br><br>
* `proxy`[`map`] - Proxy configurations to be used while configuring registry if enabled
#### Attributes Reference
* `registries_list`[`list`]: List of registries information those were created, This list would contain below attributes.<br><br>
    * `name`: Name of the registry retrieved.<br><br>
    * `role`: Role of registry retrieved.<br><br>
    * `image`: Image used for registry creation.<br><br>
    * `cluster`: Cluster to which the registry belongs.<br><br>
    * `state`: Current state of registry node.<br><br>
    * `created`: Creation time-stamp of registry.<br><br>
    * `networks`: Networks associated with the registry node.<br><br>
    * `env`: Environment variables set in the registry node.<br><br>
    * `port_mappings`: Port mappings of the registry.<br><br>

### usage of resource `rancherk3d_connect_registry`

---
This resource helps in coupling/decoupling registry from the selected cluster.
##### Sample `rancherk3d_connect_registry` config:

```terraform
resource "rancherk3d_connect_registry" "k3s-registry-1" {
  registries = ["k3s-registry-1"]
  cluster    = "k3s-default"
  connect    = true
}
```
#### Argument Reference
* `registries`[`list`] - List of registries to be connected/disconnected from the selected cluster.<br><br>
* `cluster`[`string`] - Cluster to which registries to be associated with.<br><br>
* `connect`[`boolean`] - Enable this flag if registries to be connected with specified cluster, when disabled it disconnects the registry from cluster.
#### Attributes Reference
* `status`[`list`] - Updated status of registry. This list would contain below attributes.<br><br>
    * `registry`[`string`]: Name of the registry.<br><br>
    * `cluster`[`string`]: Cluster to which the registry is either connected/disconnected.<br><br>
    * `state`[`string`]: Updated state of registry, `connected`/`disconnected`.<br><br>
### usage of resource `rancherk3d_load_image`

---
This resource helps in loading a list of images to all selected clusters.
##### Sample `rancherk3d_load_image` config:

```terraform
resource "rancherk3d_load_image" "k3s-default" {
  images       = ["basnik/terragen:latest"]
  cluster      = "k3s-default"
  keep_tarball = false
}
```
#### Argument Reference
* `images`[`list`] - List of images to be imported to the existing cluster.<br><br>
* `cluster`[`string`] - Name of the existing cluster to which the images has to be imported to.<br><br>
* `all`: [`boolean`] - If enabled loads images to all available clusters in the selected runtime.<br><br>
* `keep_tarball`[`boolean`] - Enable to keep the tarball of the loaded images locally.
#### Attributes Reference
* `images_stored`[`list`] - list of images loaded to the cluster.<br><br>
    * `cluster`[`string`]: Cluster to which the below images are stored.<br><br>
    * `images`[`list`]: List of images and its tarball stored, if in case keep_tarball is enabled.<br><br>

