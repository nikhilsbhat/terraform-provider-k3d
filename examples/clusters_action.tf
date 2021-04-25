resource rancherk3d_cluster_action "stop-k3s-cluster" {
  depends_on = [
    rancherk3d_node_create.node-1
  ]
  clusters = ["k3s-default"]
  stop     = true
}

resource rancherk3d_cluster_action "start-k3s-cluster" {
  depends_on = [
    rancherk3d_cluster_action.stop-k3s-cluster,
  ]
  clusters = ["k3s-default"]
  start    = true
}