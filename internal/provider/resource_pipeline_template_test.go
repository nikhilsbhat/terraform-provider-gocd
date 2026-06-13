//nolint:testpackage
package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func TestGetPipelineTemplateConfigParsesYAMLAndJSON(t *testing.T) {
	tests := map[string]struct {
		rawConfig string
		isYAML    bool
	}{
		"yaml": {
			rawConfig: `
name: sample-template
stages:
  - name: build
    jobs:
      - name: build
        tasks:
          - type: exec
            attributes:
              command: echo
              arguments:
                - building
              run_if: []
`,
			isYAML: true,
		},
		"json": {
			rawConfig: `{
				"name": "sample-template",
				"stages": [{
					"name": "build",
					"jobs": [{
						"name": "build",
						"tasks": [{
							"type": "exec",
							"attributes": {
								"command": "echo",
								"arguments": ["building"],
								"run_if": []
							}
						}]
					}]
				}]
			}`,
			isYAML: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			resourceData := schema.TestResourceDataRaw(t, resourcePipelineTemplate().Schema, map[string]any{
				utils.TerraformResourceName:   "sample-template",
				utils.TerraformResourceConfig: tt.rawConfig,
				utils.TerraformResourceYAML:   tt.isYAML,
			})

			templateCfg, err := getPipelineTemplateConfig(resourceData)
			if err != nil {
				t.Fatalf("unexpected error parsing template config: %v", err)
			}

			if templateCfg.Name != "sample-template" {
				t.Fatalf("unexpected template name: %s", templateCfg.Name)
			}

			if len(templateCfg.Stages) != 1 || templateCfg.Stages[0].Name != "build" {
				t.Fatalf("unexpected template stages: %#v", templateCfg.Stages)
			}

			rawTemplateCfg, err := getRawPipelineTemplateConfig(resourceData)
			if err != nil {
				t.Fatalf("unexpected error parsing raw template config: %v", err)
			}

			stages := rawTemplateCfg["stages"].([]any)
			stage := stages[0].(map[string]any)
			jobs := stage["jobs"].([]any)
			job := jobs[0].(map[string]any)
			tasks := job["tasks"].([]any)
			task := tasks[0].(map[string]any)

			if task["type"] != "exec" {
				t.Fatalf("unexpected task type in raw template config: %#v", task)
			}
		})
	}
}
