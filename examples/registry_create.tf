resource "k3d_registry" "registry" {
  name     = "k3s-registry"
  cluster  = "k3s-default"
  protocol = "http"
  host     = "basnik.com"
}

resource "k3d_registry" "registry-1" {
  name     = "k3s-registry-1"
  cluster  = "k3s-default"
  protocol = "http"
  host     = "basnik.com"
}

resource "k3d_registry" "registry-2" {
  name     = "k3s-registry-2"
  cluster  = "k3s-default"
  protocol = "http"
  expose = {
    "hostIp" : "0.0.0.0",
    "hostPort" : "5300",
  }
}