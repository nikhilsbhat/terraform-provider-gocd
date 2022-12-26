resource "gocd_plugins" "json_plugin_settings" {
  plugin_id = "json.config.plugin"
  plugin_configurations {
    key = "pipeline_pattern"
    value = "*.gocdpipeline.test.json"
  }
  plugin_configurations {
    key = "environment_pattern"
    value = "*.gocdenvironment.test.json"
  }
}