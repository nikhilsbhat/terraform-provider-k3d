data "rancherk3d_kubeconfig" "k3s-default" {
  clusters = ["k3s-default"]
  //  not_encoded = true
}

//data "rancherk3d_kubeconfig" "all_clusters" {
//  clusters = ["k3s-sample"]
//  all      = true
//}

output "rancher_kubeconfig_k3d_sample" {
  value = data.rancherk3d_kubeconfig.k3s-default.kube_config
}