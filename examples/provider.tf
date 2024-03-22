terraform {
  required_providers {
    gocd = {
      source  = "hashicorp/gocd"
      version = "~> 0.0.3"
    }
  }
}

provider "gocd" {
  base_url = "http://localhost:8155/go"
  username = "admin"
  password = "admin"
  //  auth_token = "d8fccbc997d04e917b1490af8e7bf46290ab8c99"
  loglevel = "debug"
  //  skip_check = true
  retries {
    count     = 10
    wait_time = 2
  }
}