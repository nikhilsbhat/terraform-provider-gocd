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

resource "gocd_pipeline_group" "sample_group_3" {
  name      = "sample-group-3"
  pipelines = ["testing"]
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

resource "gocd_pipeline_group" "movies" {
  name = "movies"
  pipelines = [
    "action-movies",
    "action-movies-manual",
    "animation-movies",
    "animation-and-action-movies",
    "both"
  ]
  authorization {
    view {
      users = ["nikhil"]
    }
    operate {
      users = ["nikhil"]
    }
    admins {
      users = ["nikhil"]
    }
  }
}

data "gocd_pipeline_group" "sample_group" {
  group_id = "sample-group"
}