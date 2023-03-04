#data "k3d_node_list" "k3s_default" {
#  depends_on = [
#    k3d_node_create.node-1,
#    k3d_node_create.node-2
#  ]
#  cluster = "k3s-default"
#  all     = true
##  nodes = ["test-node-from-terraform"]
#}

##data "k3d_node_list" "k3s_default_server" {
##  cluster = "k3s-default"
##  nodes   = ["k3d-k3s-default-server-0", "k3d-k3s-default-serverlb"]
##}
#
#output "rancher_nodes_list" {
#  value = data.k3d_node_list.k3s_default.node_list
#}