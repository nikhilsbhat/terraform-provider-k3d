terraform {
  required_providers {
    k3d = {
      source  = "hashicorp/k3d"
      version = "0.1.3"
    }
  }
}

provider "k3d" {
  kubernetes_version = "1.24.4-k3s1"
  k3d_api_version    = "k3d.io/v1alpha4"
  registry           = "rancher/k3s"
  kind               = "Simple"
  runtime            = "docker"
}