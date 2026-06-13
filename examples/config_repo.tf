resource "gocd_config_repository" "sample_config_repo" {
  profile_id = "sample_config_repo"
  plugin_id  = "yaml.config.plugin"
  configuration {
    key   = "username"
    value = "admin"
  }
  configuration {
    key   = "password"
    value = "admin"
  }
  configuration {
    key       = "url"
    value     = "https://github.com/config-repo/gocd-json-config-example.git"
    is_secure = false
  }
  material {
    type = "git"
    attributes {
      url                = "https://github.com/config-repo/gocd-json-config-example.git"
      username           = "bob"
      encrypted_password = "AES:hpM5puIHTWMI+4hX+Stl1A==:mEpzzXcJrSkiaBQHJ8SCiA=="
      branch             = "main"
      auto_update        = true
    }
  }
  rules = [
    {
      "directive" : "allow",
      "action" : "refer",
      "type" : "pipeline_group",
      "resource" : "*"
    }
  ]
}

resource "gocd_config_repository" "sample_config_repo_yamll" {
  profile_id = "yamll"
  plugin_id  = "yaml.config.plugin"
  material {
    type = "git"
    attributes {
      url                = "https://github.com/nikhilsbhat/yamll.git"
      username           = "admin"
      encrypted_password = "AES:Y8pqp5DR3XmWr0fc0FlOyg==:bXXYZ3qPPkTqBcIosK3MIQ=="
      branch             = "main"
      auto_update        = false
    }
  }
  rules = [
    {
      "directive" : "allow",
      "action" : "refer",
      "type" : "pipeline_group",
      "resource" : "*"
    }
  ]
}

data "gocd_config_repository" "sample_config_repo" {
  profile_id = gocd_config_repository.sample_config_repo.id
}