#resource rancherk3d_cluster_action "stop-k3s-cluster" {
#  clusters = ["test"]
#  stop     = false
#}
#
#resource rancherk3d_cluster_action "start-k3s-cluster" {
#  depends_on = [
#    rancherk3d_cluster_action.stop-k3s-cluster,
#  ]
#  clusters = ["k3s-default"]
#  start    = true
#}