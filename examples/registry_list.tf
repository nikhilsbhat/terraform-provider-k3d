data "k3d_registry_list" "registry-1" {
  depends_on = [
    k3d_create_registry.registry,
  ]
  cluster = "k3s-default"
  registries = [
  "k3s-registry-2"]
  all = true
}

#output "rancher_registry_list" {
#  value = data.k3d_registry_list.registry-1.registries_list
#}