resource "gocd_pipeline" "helm_drift" {
  name              = "helm-drift"
  pause_on_creation = true
  pause_reason      = "better to pause pipeline on creation"
  group             = "sample-group"
  config            = <<EOF
  environment_variables:
    - name: HELM_PLUGIN
      secure: false
      value: "false"
  lock_behavior: none
  materials:
    - attributes:
        auto_update: false
        branch: master
        destination: null
        filter: null
        invert_filter: false
        name: null
        shallow_clone: false
        submodule_folder: null
        url: "https://github.com/nikhilsbhat/helm-drift.git"
      type: git
  name: helm-drift
  parameters: []
  stages:
    - approval:
        allow_only_on_success: false
        authorization:
          roles: []
          users: []
        type: success
      clean_working_directory: false
      environment_variables: []
      fetch_materials: true
      jobs:
        - artifacts: []
          environment_variables: []
          name: lint
          resources: []
          run_instance_count: null
          tabs: []
          tasks:
            - attributes:
                arguments:
                - lint
                command: echo
                run_if: []
              type: exec
          timeout: null
      name: lint
      never_cleanup_artifacts: false
  timer: null
  tracking_tool: null
EOF
}

data "gocd_pipeline" "helm_images" {
  name = "helm-images"
  yaml = true
}

data "gocd_pipeline" "helm_drift" {
  depends_on = [gocd_pipeline.helm_drift]
  name       = "helm-drift"
  yaml       = true
}