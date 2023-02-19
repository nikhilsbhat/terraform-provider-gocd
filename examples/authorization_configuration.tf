resource "gocd_auth_config" "password_file_config" {
  profile_id = "admin_new"
  plugin_id  = "cd.go.authentication.passwordfile"
  properties {
    key   = "PasswordFilePath"
    value = "/Users/nikhil.bhat/idfc/gocd-setup/.gocdadmin2"
  }
}

data "gocd_auth_config" "password_file_config" {
  profile_id = gocd_auth_config.password_file_config.id
}