resource "gocd_agent" "sample_agent" {
  uuid         = "7c253df5-a262-4573-bcc8-aedc7e87317d"
  environments = ["sample_environment_3"]
}

data "gocd_agent" "sample_agent" {
  uuid = "3a3d8e62-6103-4d05-be92-11fbdf21e945"
}