resource "gocd_secret_config" "sample_kube_secret_config" {
  profile_id  = "sample-kube-secret-config-new"
  plugin_id   = "cd.go.contrib.secrets.kubernetes"
  description = "sample secret new config"
  properties {
    key   = "kubernetes_secret_name"
    value = "ci_secret"
  }
  properties {
    key   = "kubernetes_cluster_url"
    value = "https://0.0.0.0:64527"
  }
  properties {
    key   = "security_token"
    value = "AES:1poYVwRFAh8geGFaeY0GiQ==:7YQOCRx6sIMG9OjBkH0pdUvr1qJuokihboN+D0JBXzsrrFrmooItzZbyDFav/EcO"
  }
  properties {
    key   = "kubernetes_cluster_ca_cert"
    value = "AES:frDz530rq4p6ZgzQUD0X5Q==:Ra9Ldo9TgwcvrzSwPoU4g1KqgDi0ByZXPV7oayJYMxsSXqsIUsDUrS4cyAtq5pPz"
  }
  properties {
    key   = "namespace"
    value = "default"
  }
  rules = [
    {
      action    = "refer",
      directive = "allow",
      resource  = "*",
      type      = "*"
    },
  ]
}

data "gocd_secret_config" "sample_kube_secret_config" {
  profile_id = "sample-kube-secret-config"
}