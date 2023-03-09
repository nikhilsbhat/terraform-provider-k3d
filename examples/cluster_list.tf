data "k3d_cluster" "k3s-default" {
  clusters = ["k3s-default"]
  all      = true
}

data "k3d_cluster" "k3s-sample" {
  clusters = ["k3s-sample"]
  all      = true
}

output "rancher_cluster_list" {
  value = data.k3d_cluster.k3s-default.clusters_list
}