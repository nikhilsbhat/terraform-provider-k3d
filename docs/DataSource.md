# `terraform-provider-rancherk3d` Datasource Usage
## Usage
### usage for datasource `rancherk3d_cluster_list`
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

### usage for datasource `rancherk3d_node_list`
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

### usage for datasource `rancherk3d_kubeconfig`
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

### usage for datasource `rancherk3d_registry_list`
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
    * `port_mappings`: Port mappings of the registry.<br><br>