terraform {
  required_providers {
    gocd = {
      source = "hashicorp/gocd"
      version = "~> 0.0.1"
    }
  }
}

provider "gocd" {
  base_url = "http://localhost:8153/go"
  username = "admin"
  password = "admin"
  loglevel = "debug"
}