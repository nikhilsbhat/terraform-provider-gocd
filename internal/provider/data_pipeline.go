package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
	"gopkg.in/yaml.v3"
)

func dataSourcePipeline() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourcePipelineRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The name of the pipeline to be retrieved.",
			},
			"yaml": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "When set, yaml equivalent config would be set under `config`.",
			},
			"config": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Required:    false,
				Description: "The config of the selected pipeline (it would be in yaml/json based on the attribute set).",
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Etag used to track the pipeline config",
			},
		},
	}
}

func datasourcePipelineRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()

	if len(id) == 0 {
		resourceID := utils.String(d.Get(utils.TerraformResourceName))
		id = resourceID
	}

	isYAML := utils.Bool(d.Get(utils.TerraformResourceYAML))

	response, err := defaultConfig.GetPipelineConfig(id)
	if err != nil {
		return diag.Errorf("getting pipeline configuration %s errored with: %v", id, err)
	}

	pipelineCfg, err := getPipelineConfigYaml(response, isYAML)
	if err != nil {
		return diag.Errorf("translating pipeline config to json/yaml errored with: %v", err)
	}

	if err = d.Set(utils.TerraformResourceConfig, pipelineCfg); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceConfig, err)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	d.SetId(id)

	return nil
}

func getPipelineConfigYaml(pipelineCfg gocd.PipelineConfig, isYaml bool) (string, error) {
	if !isYaml {
		valueJSON, err := json.Marshal(pipelineCfg.Config)
		if err != nil {
			return "", err
		}

		return string(valueJSON), nil
	}

	valueYAML, err := yaml.Marshal(pipelineCfg.Config)
	if err != nil {
		return "", err
	}

	return string(valueYAML), nil
}
