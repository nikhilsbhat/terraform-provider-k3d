output "cluster_config_yaml" {
  value = k3d_cluster_create.sample_cluster.config_yaml
}