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
      url         = "https://github.com/config-repo/gocd-json-config-example.git"
      username    = "bob"
      password    = "aSdiFgRRZ6A="
      branch      = "master"
      auto_update = false
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
