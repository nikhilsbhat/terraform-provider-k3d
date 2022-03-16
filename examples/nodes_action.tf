
//resource "rancherk3d_node_action" "k3s-default" {
//  nodes   = ["k3d-k3s-default-agent-0"]
//  cluster = "k3s-default"
//  stop    = true
//  all     = true
//}

#resource "rancherk3d_node_action" "k3s-default-2" {
#  nodes   = ["k3d-k3s-default-serverlb", "k3d-k3s-default-server-0"]
#  cluster = "k3s-default"
#  start   = true
#}

//resource "rancherk3d_node_action" "k3s-default-3" {
//  nodes   = ["k3d-k3s-default-agent-0", "k3d-k3s-default-agent-1"]
//  cluster = "k3s-default"
//  stop    = true
//}