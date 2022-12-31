output "encrypted_value" {
  value     = gocd_encrypt_value.new_value.encrypted_value
  sensitive = true
}