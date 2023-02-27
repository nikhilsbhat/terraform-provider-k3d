resource "rancherk3d_connect_registry" "k3s-registry-1" {
  registries = [rancherk3d_create_registry.registry-1.host]
  cluster    = "k3s-default"
  connect    = true
}