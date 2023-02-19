resource "gocd_environment" "sample_environment" {
  name = "sample_environment"
  pipelines = [
    "gocd-prometheus-exporter",
    "helm-images",
  ]
  environment_variables {
    name  = "TEST_ENV11"
    value = "value_env11"
  }
  environment_variables {
    name  = "TEST_ENV12"
    value = "value_env18"
  }
  environment_variables {
    name  = "TEST_ENV13"
    value = "value_env13"
  }
}

data "gocd_environment" "sample_environment" {
  name = gocd_environment.sample_environment.id
}