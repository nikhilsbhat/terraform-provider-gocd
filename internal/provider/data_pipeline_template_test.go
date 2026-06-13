//nolint:testpackage
package provider

import (
	"strings"
	"testing"

	"github.com/nikhilsbhat/gocd-sdk-go"
)

func TestGetPipelineTemplateConfigString(t *testing.T) {
	templateCfg := gocd.Template{
		Name: "sample-template",
		Stages: []gocd.PipelineStageConfig{
			{
				Name: "build",
				Jobs: []gocd.PipelineJobConfig{
					{Name: "build"},
				},
			},
		},
		ETAG: "template-etag",
	}

	jsonConfig, err := getPipelineTemplateConfigString(templateCfg, false)
	if err != nil {
		t.Fatalf("unexpected error marshaling template json: %v", err)
	}

	if !strings.Contains(jsonConfig, `"name":"sample-template"`) {
		t.Fatalf("unexpected json template config: %s", jsonConfig)
	}

	if strings.Contains(jsonConfig, "template-etag") {
		t.Fatalf("template etag leaked into config: %s", jsonConfig)
	}

	yamlConfig, err := getPipelineTemplateConfigString(templateCfg, true)
	if err != nil {
		t.Fatalf("unexpected error marshaling template yaml: %v", err)
	}

	if !strings.Contains(yamlConfig, "name: sample-template") {
		t.Fatalf("unexpected yaml template config: %s", yamlConfig)
	}
}
