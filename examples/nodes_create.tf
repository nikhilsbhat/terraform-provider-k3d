resource "k3d_node" "node-1" {
  name     = "sample-node-2"
  cluster  = k3d_cluster.sample_cluster.name
  role     = "agent"
  replicas = 1
  #  memory   = "8g"
  //  wait     = false
  //  timeout  = 1
}

resource "k3d_node" "node-2" {
  name     = "test-node-terraform"
  cluster  = "k3s-default"
  role     = "agent"
  replicas = 2
  // volume = "8g"
  //  wait     = true
  //  timeout  = 3
}