resource "k3d_node_action" "k3s-default" {
  depends_on = [
  k3d_node.node-1]
  nodes = [
    "test-node-from-terraform-0",
  "test-node-terraform-0"]
  cluster = "k3s-default"
  start   = true
  stop    = true
  all     = true
}

resource "k3d_node_action" "k3s-default-2" {
  depends_on = [
  k3d_node.node-1]
  nodes = [
  "test-node-from-terraform"]
  cluster = "k3s-default"
  start   = true
}

resource "k3d_node_action" "k3s-default-3" {
  nodes = [
    "k3d-k3s-default-agent-0",
  "k3d-k3s-default-agent-1"]
  cluster = "k3s-default"
  stop    = true
}