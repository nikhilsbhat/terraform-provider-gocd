resource "gocd_agent" "sample_agent" {
  uuid         = "b9101230-daa8-4e47-bc0f-d010b3d49e04"
  environments = ["sample_environment_3"]
}

data "gocd_agent" "sample_agent" {
  uuid = "b9101230-daa8-4e47-bc0f-d010b3d49e04"
}