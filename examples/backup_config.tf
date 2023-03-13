resource "gocd_backup_config" "nightly_backup" {
  schedule           = "0 0 2 * * ?"
  post_backup_script = "path/to/postbackup_script.sh"
}
