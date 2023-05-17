resource "gocd_agent" "sample_agent" {
  uuid         = "fe61ac36-1c16-4260-93d3-23110d94b38a"
  environments = ["sample_environment_3"]
}

data "gocd_agent" "sample_agent" {
  uuid = "aaf50aed-cfdb-4c20-8989-f17a2ba54739"
}