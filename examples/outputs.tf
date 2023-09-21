output "encrypted_value" {
  value     = gocd_encrypt_value.new_value.encrypted_value
  sensitive = true
}

output "sample_config_repo" {
  value = data.gocd_config_repository.sample_config_repo.material
}

output "password_file_config" {
  value = data.gocd_auth_config.password_file_config.properties
}

output "ec2_cluster_profile" {
  value = data.gocd_cluster_profile.ec2_cluster_profile.properties
}

output "sample_environment" {
  value = data.gocd_environment.sample_environment.pipelines
}

output "sample_ec2" {
  value = data.gocd_elastic_agent_profile.sample_ec2.properties
}

output "yaml_plugin_settings" {
  value = data.gocd_plugin_setting.yaml_plugin_settings.configuration
}

output "sample_kube_secret_config" {
  value = data.gocd_secret_config.sample_kube_secret_config.plugin_id
}

output "kubernetes_plugin" {
  value = data.gocd_plugin_info.kubernetes_plugin
}

output "sample_agent_config" {
  value = data.gocd_agent.sample_agent.hostname
}

output "helm_images" {
  value = data.gocd_pipeline.helm_images.config
}

output "helm_drift" {
  value = gocd_pipeline.helm_drift.config
}

output "docker_artifact_store" {
  value = data.gocd_artifact_store.docker.properties
}

output "sample_role" {
  value = data.gocd_role.sample.users
}

output "sample_ldap_role" {
  value = data.gocd_role.sample_ldap.properties
}

output "pipeline_group_movies" {
  value = data.gocd_pipeline_group.movies.authorization
}