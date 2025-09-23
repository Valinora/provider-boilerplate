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

resource "devops_engineer" "basic" {
  name  = "Colin"
  email = "colin@liatrio.com"
}

resource "devops_engineer" "other" {
  name  = "Austin"
  email = "austin@liatrio.com"
}


resource "devops_engineer" "newname" {
  name  = "ajdfa"
  email = "djjjdjd@sdasds.com"
}
