resource "k3d_cluster_create" "sample_cluster" {
  name          = "default"
  servers_count = 1
  agents_count  = 2
  //  image = "rancher/k3s:v1.24.4-k3s1"
  kube_api {
    host_ip   = "0.0.0.0"
    host_port = 6445
  }
  //
  //  ports {
  //    host_port = 8080
  //    container_port = 80
  //    node_filters = [
  //      "loadbalancer",
  //    ]
  //  }

  k3d_options {
    no_loadbalancer = false
    no_image_volume = false
  }

  kube_config {
    update_default = true
    switch_context = true
  }
}