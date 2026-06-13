package provider

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	gocdclient "github.com/nikhilsbhat/terraform-provider-gocd/pkg/client"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
	"gopkg.in/yaml.v3"
)

func resourcePipelineTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelineTemplateCreate,
		ReadContext:   resourcePipelineTemplateRead,
		UpdateContext: resourcePipelineTemplateUpdate,
		DeleteContext: resourcePipelineTemplateDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The name of the pipeline template to be created or updated.",
			},
			"config": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    false,
				Description: "The pipeline template config to be created or updated. It can be YAML or JSON based on the `yaml` attribute.",
			},
			"yaml": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    false,
				ForceNew:    false,
				Description: "Set to true when the template config declared under `config` is YAML.",
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Etag used to track the pipeline template config.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourcePipelineTemplateImport,
		},
	}
}

func resourcePipelineTemplateCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	templateClient := meta.(gocdclient.PipelineTemplateClient)

	if !d.IsNewResource() {
		return nil
	}

	id := d.Id()
	if len(id) == 0 {
		id = utils.String(d.Get(utils.TerraformResourceName))
	}

	templateCfg, err := getPipelineTemplateConfig(d)
	if err != nil {
		return diag.Errorf("decoding pipeline template '%s' config errored with: %v", id, err)
	}

	if templateCfg.Name != id {
		return diag.Errorf("pipeline template name passed under attribute and template config are not same, make sure to pass the same values, "+
			"current values: 'attribute:%s config:%s'", id, templateCfg.Name)
	}

	rawTemplateCfg, err := getRawPipelineTemplateConfig(d)
	if err != nil {
		return diag.Errorf("decoding pipeline template '%s' raw config errored with: %v", id, err)
	}

	if _, err = templateClient.CreateTemplateRaw(rawTemplateCfg); err != nil {
		return diag.Errorf("creating pipeline template '%s' errored with: %v", id, err)
	}

	d.SetId(id)

	return resourcePipelineTemplateRead(ctx, d, meta)
}

func resourcePipelineTemplateRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	name := utils.String(d.Get(utils.TerraformResourceName))
	if len(name) == 0 {
		name = d.Id()
	}

	response, err := defaultConfig.GetTemplate(name)
	if err != nil {
		return diag.Errorf("getting pipeline template '%s' errored with: %v", name, err)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	return nil
}

func resourcePipelineTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	templateClient := meta.(gocdclient.PipelineTemplateClient)

	if !d.HasChanges(utils.TerraformResourceConfig, utils.TerraformResourceYAML) {
		log.Printf("nothing to update so skipping")

		return nil
	}

	templateCfg, err := getPipelineTemplateConfig(d)
	if err != nil {
		return diag.Errorf("decoding pipeline template '%s' config errored with: %v", d.Id(), err)
	}

	rawTemplateCfg, err := getRawPipelineTemplateConfig(d)
	if err != nil {
		return diag.Errorf("decoding pipeline template '%s' raw config errored with: %v", d.Id(), err)
	}

	if _, err = templateClient.UpdateTemplateRaw(templateCfg.Name, utils.String(d.Get(utils.TerraformResourceEtag)), rawTemplateCfg); err != nil {
		return diag.Errorf("updating pipeline template '%s' errored with: %v", templateCfg.Name, err)
	}

	return resourcePipelineTemplateRead(ctx, d, meta)
}

func resourcePipelineTemplateDelete(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()
	if len(id) == 0 {
		return diag.Errorf("resource with the ID '%s' not found", id)
	}

	if err := defaultConfig.DeleteTemplate(id); err != nil {
		return diag.Errorf("deleting pipeline template '%s' errored with: %v", id, err)
	}

	d.SetId("")

	return nil
}

func resourcePipelineTemplateImport(_ context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	defaultConfig := meta.(gocd.GoCd)

	templateName := utils.String(d.Id())

	response, err := defaultConfig.GetTemplate(templateName)
	if err != nil {
		return nil, err
	}

	if err = d.Set(utils.TerraformResourceName, templateName); err != nil {
		return nil, err
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func getPipelineTemplateConfig(d *schema.ResourceData) (gocd.Template, error) {
	var templateCfg gocd.Template

	config := utils.String(d.Get(utils.TerraformResourceConfig))
	if utils.Bool(d.Get(utils.TerraformResourceYAML)) {
		if err := yaml.Unmarshal([]byte(config), &templateCfg); err != nil {
			return templateCfg, err
		}

		return templateCfg, nil
	}

	if err := json.Unmarshal([]byte(config), &templateCfg); err != nil {
		return templateCfg, err
	}

	return templateCfg, nil
}

func getRawPipelineTemplateConfig(d *schema.ResourceData) (map[string]any, error) {
	var templateCfg map[string]any

	config := utils.String(d.Get(utils.TerraformResourceConfig))
	if utils.Bool(d.Get(utils.TerraformResourceYAML)) {
		if err := yaml.Unmarshal([]byte(config), &templateCfg); err != nil {
			return templateCfg, err
		}

		return templateCfg, nil
	}

	if err := json.Unmarshal([]byte(config), &templateCfg); err != nil {
		return templateCfg, err
	}

	return templateCfg, nil
}
