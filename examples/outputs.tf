output "rancher_cluster_list" {
  value = data.k3d_cluster.k3s-default.clusters_list
}

output "rancher_kubeconfig_k3d_sample" {
  value     = data.k3d_kubeconfig.k3s-default.kube_config
  sensitive = true
}

output "rancher_nodes_list" {
  value = data.k3d_node.k3s_default.node_list
}

output "rancher_registry_list" {
  value = data.k3d_registry.registry-1.registries_list
}

output "cluster_config_yaml" {
  value = k3d_cluster.sample_cluster.config_yaml
}