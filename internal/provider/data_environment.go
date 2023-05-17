package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func dataSourceEnvironment() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceEnvironmentRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The name of environment.",
			},
			"pipelines": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Required:    false,
				Description: "List of pipeline names that should be added to this environment.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"environment_variables": environmentsSchemaData(),
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Etag used to track the environment configuration",
			},
		},
	}
}

func datasourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()

	if len(id) == 0 {
		resourceID := utils.String(d.Get(utils.TerraformResourceName))
		id = resourceID
	}

	envName := utils.String(d.Get(utils.TerraformResourceName))

	response, err := defaultConfig.GetEnvironment(envName)
	if err != nil {
		return diag.Errorf("getting environment %s errored with: %v", envName, err)
	}

	if err = d.Set(utils.TerraformResourcePipelines, flattenPipelines(response.Pipelines)); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourcePipelines)
	}

	flattenedEnvVars, err := utils.MapSlice(response.EnvVars)
	if err != nil {
		d.SetId("")

		return diag.Errorf("errored while flattening environment variable obtained: %v", err)
	}

	if err = d.Set(utils.TerraformResourceEnvVar, flattenedEnvVars); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceEnvVar)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	d.SetId(id)

	return nil
}

func flattenPipelines(pipelines []gocd.Pipeline) []interface{} {
	pipeline := make([]interface{}, len(pipelines))
	for i, name := range pipelines {
		pipeline[i] = name.Name
	}

	return pipeline
}
