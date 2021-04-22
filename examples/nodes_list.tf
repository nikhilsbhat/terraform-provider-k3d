data "rancherk3d_node_list" "k3s_default" {
  cluster = "k3s-default"
  all     = true
}

data "rancherk3d_node_list" "k3s_default_server" {
  cluster = "k3s-default"
  nodes     = ["k3d-k3s-default-server-0"]
}

output "rancher_nodes_list" {
  value = data.rancherk3d_node_list.k3s_default.node_list
}