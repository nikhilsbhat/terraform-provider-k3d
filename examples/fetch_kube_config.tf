data "k3d_kubeconfig" "k3s-default" {
  clusters = ["k3s-default"]
  //  not_encoded = true
  all = true
}

data "k3d_kubeconfig" "all_clusters" {
  clusters = ["k3s-sample"]
  all      = true
}
