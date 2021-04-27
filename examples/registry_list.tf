data "rancherk3d_registry_list" "registry-1" {
  //  cluster = "k3s-default"
  //  registries = ["k3d-k3s-sample-registr"]
  all = true
}

output "rancher_registry_list" {
  value = data.rancherk3d_registry_list.registry-1.registries_list
}