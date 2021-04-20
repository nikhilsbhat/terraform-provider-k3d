
data "rancherk3d_cluster_list" "k3s-default" {
  clusters = ["k3s-default"]
}

data "rancherk3d_cluster_list" "k3s-sample" {
  clusters = ["k3s-sample"]
  all      = true
}

output "rancher_cluster_list" {
  value = data.rancherk3d_cluster_list.k3s-sample.clusters_list
}