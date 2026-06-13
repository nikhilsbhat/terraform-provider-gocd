resource "gocd_pipeline_template" "sample" {
  name = "sample-template"
  yaml = true

  config = <<EOF
name: sample-template
stages:
  - name: build
    fetch_materials: true
    clean_working_directory: false
    never_cleanup_artifacts: false
    jobs:
      - name: build
        tasks:
          - type: exec
            attributes:
              command: echo
              arguments:
                - building
              run_if: []
EOF
}

data "gocd_pipeline_template" "sample" {
  depends_on = [gocd_pipeline_template.sample]

  name = "sample-template"
  yaml = true
}
