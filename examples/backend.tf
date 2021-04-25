terraform {
  backend "local" {
    path = "/data/terraform-backend/terraform.tfstate"
  }
}