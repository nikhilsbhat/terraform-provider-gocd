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

func dataSourcePipelineTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourcePipelineTemplateRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The name of the pipeline template to be retrieved.",
			},
			"yaml": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "When set, YAML equivalent template config will be set under `config`.",
			},
			"config": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The config of the selected pipeline template. It will be YAML or JSON based on the `yaml` attribute.",
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Etag used to track the pipeline template config.",
			},
		},
	}
}

func datasourcePipelineTemplateRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()
	if len(id) == 0 {
		id = utils.String(d.Get(utils.TerraformResourceName))
	}

	response, err := defaultConfig.GetTemplate(id)
	if err != nil {
		return diag.Errorf("getting pipeline template %s errored with: %v", id, err)
	}

	templateCfg, err := getPipelineTemplateConfigString(response, utils.Bool(d.Get(utils.TerraformResourceYAML)))
	if err != nil {
		return diag.Errorf("translating pipeline template config to json/yaml errored with: %v", err)
	}

	err = d.Set(utils.TerraformResourceConfig, templateCfg)
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceConfig, err)
	}

	err = d.Set(utils.TerraformResourceEtag, response.ETAG)
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	d.SetId(id)

	return nil
}

func getPipelineTemplateConfigString(templateCfg gocd.Template, isYAML bool) (string, error) {
	templateConfig := map[string]any{
		utils.TerraformResourceName: templateCfg.Name,
		"stages":                    templateCfg.Stages,
	}

	if !isYAML {
		valueJSON, err := json.Marshal(templateConfig)
		if err != nil {
			return "", err
		}

		return string(valueJSON), nil
	}

	valueYAML, err := yaml.Marshal(templateConfig)
	if err != nil {
		return "", err
	}

	return string(valueYAML), nil
}
