resource "gocd_agent" "sample_agent" {
  uuid         = "bbfe3a75-7fd8-48db-af32-0a91b9efd0ab"
  environments = ["sample_environment_3"]
}

data "gocd_agent" "sample_agent" {
  uuid = "bbfe3a75-7fd8-48db-af32-0a91b9efd0ab"
}