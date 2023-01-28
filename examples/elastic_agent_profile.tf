resource "gocd_elastic_agent_profile" "sample_kubernetes" {
  profile_id         = "sample_kubernetes"
  cluster_profile_id = gocd_cluster_profile.kube_cluster_profile.profile_id
  properties {
    key   = "Image"
    value = "basnik/thanosbench:with-thanos"
  }
  properties {
    key   = "MaxMemory"
    value = "2G"
  }
  properties {
    key   = "MaxCPU"
    value = "500M"
  }
  properties {
    key   = "Privileged"
    value = "false"
  }
  properties {
    key   = "PodSpecType"
    value = "properties"
  }
}

resource "gocd_elastic_agent_profile" "sample_ec2" {
  profile_id         = "sample_ec2"
  cluster_profile_id = gocd_cluster_profile.ec2_cluster_profile.profile_id
  properties {
    key   = "ec2_ami"
    value = "test-image"
  }
  properties {
    key   = "ec2_instance_type"
    value = "t2.micro"
  }
  properties {
    key   = "ec2_sg"
    value = "sg-12pjexw8121uwj"
  }
  properties {
    key   = "ec2_subnets"
    value = "aws-subnet-sddaekjkjddlfg,aws-subnet-sksdcndjfergfg"
  }
  properties {
    key   = "ec2_key"
    value = "test-key.pem"
  }
  properties {
    key   = "ec2_instance_profile"
    value = "default"
  }
  properties {
    key   = "go_agent_work_dir"
    value = "/data"
  }
  properties {
    key   = "ec2_user_data"
    value = "echo hi"
  }
}

data "gocd_elastic_agent_profile" "sample_ec2" {
  profile_id         = "sample_ec2"
}