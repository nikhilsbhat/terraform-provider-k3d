#data "rancherk3d_registry_list" "registry-1" {
#  depends_on = [
#    rancherk3d_create_registry.registry,
##    rancherk3d_create_registry.registry-1,
##    rancherk3d_create_registry.registry-2,
##    rancherk3d_connect_registry.k3s-registry-1
#  ]
#  cluster    = "k3s-default"
#  registries = ["k3s-registry-2"]
#  all        = true
#}
#
#output "rancher_registry_list" {
#  value = data.rancherk3d_registry_list.registry-1.registries_list
#}