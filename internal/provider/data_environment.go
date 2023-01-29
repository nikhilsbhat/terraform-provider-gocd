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
		newID, err := utils.GetRandomID()
		if err != nil {
			d.SetId("")

			return diag.Errorf("errored while fetching randomID %v", err)
		}
		id = newID
	}

	envName := utils.String(d.Get(utils.TerraformResourceName))

	response, err := defaultConfig.GetEnvironment(envName)
	if err != nil {
		return diag.Errorf("getting environment %s errored with: %v", envName, err)
	}

	if err = d.Set(utils.TerraformPipelines, flattenPipelines(response.Pipelines)); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformPipelines)
	}

	flattenedEnvVars, err := utils.MapSlice(response.EnvVars)
	if err != nil {
		d.SetId("")

		return diag.Errorf("errored while flattening environment variable obtained: %v", err)
	}

	if err = d.Set(utils.TerraformEnvVar, flattenedEnvVars); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformEnvVar)
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
