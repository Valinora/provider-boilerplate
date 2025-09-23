terraform {
  required_providers {
    devops = {
      source = "liatr.io/terraform/devops-bootcamp"
    }
  }
}

provider "devops" {
  host = "http://localhost:8080"
}

data "devops_engineer" "engineer" {
}

output "engineer" {
  value = data.devops_engineer.engineer
}
