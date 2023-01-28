//resource "gocd_plugin_setting" "yaml_plugin_settings" {
//  plugin_id = "yaml.config.plugin"
//  plugin_configurations {
//    key = "file_pattern"
//    value = "*.gocd.yaml"
//  }
//  plugin_configurations {
//    key = "file_pattern"
//    value = "*.gocd.yam"
//  }
//}

data "gocd_plugin_setting" "yaml_plugin_settings" {
  plugin_id = "yaml.config.plugin"
}