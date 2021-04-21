data "rancherk3d_node_list" "k3s_default" {
  cluster = "k3s-default"
  all     = true
}

output "rancher_nodes_list" {
  value = data.rancherk3d_node_list.k3s_default.node_list
}