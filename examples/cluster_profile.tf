resource "gocd_cluster_profile" "kube_cluster_profile" {
  profile_id = "kube_cluster_profile"
  plugin_id  = "cd.go.contrib.elasticagent.kubernetes"
  properties {
    key   = "go_server_url"
    value = "https://gocd.sample.com/go"
  }
  properties {
    key   = "auto_register_timeout"
    value = "15"
  }
  properties {
    key   = "pending_pods_count"
    value = "6"
  }
  properties {
    key   = "kubernetes_cluster_url"
    value = "https://0.0.0.0:64527"
  }
  properties {
    key   = "security_token"
    value = "ZGVuaHdlMzQ4NW54ZGN3a3djMjRybmN3ZWZua3dqeGYwMjM0bnhrd2Zqa3NqZGY="
  }
  properties {
    key   = "kubernetes_cluster_ca_cert"
    value = "ZGVuaHdlMzQ4NW54ZGN3a3djMjRybmN3ZWZua3dqeGYwMjM0bnhrd2Zqa3NqZGY="
  }
  properties {
    key   = "namespace"
    value = "default"
  }
}

resource "gocd_cluster_profile" "ec2_cluster_profile" {
  profile_id = "ec2_cluster_profile"
  plugin_id  = "com.continuumsecurity.elasticagent.ec2"
  properties {
    key   = "go_server_url"
    value = "https://gocd.sample.com/go"
  }
  properties {
    key   = "auto_register_timeout"
    value = "60s"
  }
  properties {
    key   = "max_elastic_agents"
    value = "2"
  }
  properties {
    key   = "aws_region"
    value = "ap-south-1"
  }
  properties {
    key   = "aws_profile"
    value = "dev"
  }
  properties {
    key   = "aws_endpoint_url"
    value = "https://0.0.0.0:64527"
  }
  properties {
    key   = "aws_access_key_id"
    value = "dsdfmldfkjgcdfjgcdjfdfjgkcdfgfn"
  }
  properties {
    key   = "aws_secret_access_key"
    value = "dsdfmldfkjgcdfjgcdjfdfjgkcdfgfn"
  }
}

data "gocd_cluster_profile" "ec2_cluster_profile" {
  profile_id = gocd_cluster_profile.ec2_cluster_profile.id
}