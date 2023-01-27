output "encrypted_value" {
  value     = gocd_encrypt_value.new_value.encrypted_value
  sensitive = true
}

output "sample_config_repo" {
  value = data.gocd_config_repository.sample_config_repo.material
}