resource "gocd_pipeline_group" "sample_group_2" {
  name      = "sample-group-2"
  pipelines = ["helm-images", "helm-drift"]
  authorization {
    view {
      users = ["nikhil"]
      roles = ["sample"]
    }
    operate {
      users = ["nikhil"]
      roles = ["sample"]
    }
    admins {
      users = ["nikhil"]
      roles = ["sample"]
    }
  }
}

data "gocd_pipeline_group" "sample_group" {
  group_id = "sample-group"
}