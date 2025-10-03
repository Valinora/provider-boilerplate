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

resource "devops_engineer" "me" {
  name  = "Colin"
  email = "colin@liatrio.com"
}


resource "devops_engineer" "spam" {
  name  = "ajdfa"
  email = "djjjdjd@sdasds.com"
}

resource "devops_dev" "lmao" {
  name      = "new-goblins"
  engineers = [devops_engineer.me.id, devops_engineer.spam.id]
}

resource "devops_ops" "noops" {
  name      = "losers"
  engineers = [devops_engineer.me.id]
}
