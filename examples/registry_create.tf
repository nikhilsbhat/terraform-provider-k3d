#resource "rancherk3d_create_registry" "registry" {
#  name     = "k3s-registry"
#  cluster  = "k3s-default"
#  protocol = "http"
#  host     = "basnik.com"
#}

resource "rancherk3d_create_registry" "registry-1" {
  name     = "k3s-registry-1"
  cluster  = "k3s-default"
  protocol = "http"
}