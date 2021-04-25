resource "rancherk3d_load_image" "k3s-default" {
  depends_on = [
    rancherk3d_node_action.k3s-default-2,
  ]

  images       = ["basnik/terragen:latest"]
  cluster      = "k3s-default"
  keep_tarball = false
}
