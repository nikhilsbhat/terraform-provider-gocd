resource "gocd_artifact_store" "docker" {
  store_id  = "docker"
  plugin_id = "cd.go.artifact.s3"
  properties {
    key   = "S3Bucket"
    value = "sample"
  }
  properties {
    key   = "Region"
    value = "ap-south-1"
  }
}

data "gocd_artifact_store" "docker" {
  store_id = gocd_artifact_store.docker.id
}