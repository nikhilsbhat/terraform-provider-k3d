resource "rancherk3d_load_image" "k3s-default" {
  images       = ["basnik/terragen:latest", "basnik/renderer:latest"]
  cluster      = "k3s-default"
  keep_tarball = false
}
