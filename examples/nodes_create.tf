resource "rancherk3d_node_create" "node-1" {
  name     = "test-node-from-terraform"
  cluster  = "k3s-default"
  role     = "agent"
  replicas = 1
  //  wait     = false
  //  timeout  = 1
}

resource "rancherk3d_node_create" "node-2" {
  name     = "test-node-terraform"
  cluster  = "k3s-default"
  role     = "agent"
  replicas = 2
  // volume = "8g"
  //  wait     = true
  //  timeout  = 3
}