resource "gocd_agent" "sample_agent" {
  uuid         = "3a3d8e62-6103-4d05-be92-11fbdf21e945"
  environments = ["sample_environment_3"]
}

data "gocd_agent" "sample_agent" {
  uuid = "3a3d8e62-6103-4d05-be92-11fbdf21e945"
}