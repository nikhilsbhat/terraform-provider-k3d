terraform {
  required_providers {
    rancherk3d = {
      source  = "hashicorp/rancherk3d"
      version = "0.1.2"
    }
  }
}

provider "rancherk3d" {
  kubernetes_version = "1.20.2-k3s1"
  k3d_api_version    = "k3d.io/v1alpha2"
  registry           = "rancher/k3s"
  kind               = "Simple"
  runtime            = "docker"
  version            = "0.1.2"
}