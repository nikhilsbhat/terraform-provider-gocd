package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		DeleteContext: resourceEnvironmentDelete,
		UpdateContext: resourceEnvironmentUpdate,
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
				Optional:    true,
				Computed:    false,
				Description: "List of pipeline names that should be added to this environment.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"environment_variables": environmentsSchemaResource(),
			"etag": {
				Type:        schema.TypeString,
				Required:    false,
				Computed:    true,
				ForceNew:    false,
				Description: "etag used to track the environment configurations.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceEnvironmentImport,
		},
	}
}

type environmentChanges struct {
	envVarsChanges  []gocd.EnvVars
	pipelineChanges []gocd.Pipeline
	equal           bool
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !d.IsNewResource() {
		return nil
	}
	id := d.Id()

	if len(id) == 0 {
		resourceID := utils.String(d.Get(utils.TerraformResourceName))
		id = resourceID
	}

	envVars, err := getEnvironments(d.Get(utils.TerraformResourceEnvVar))
	if err != nil {
		return diag.Errorf("reading environment errored with %v", err)
	}

	cfg := gocd.Environment{
		Name:      utils.String(d.Get(utils.TerraformResourceName)),
		Pipelines: getPipelines(d.Get(utils.TerraformResourcePipelines)),
		EnvVars:   envVars,
	}

	if err = defaultConfig.CreateEnvironment(cfg); err != nil {
		return diag.Errorf("creating environment %s errored with %v", cfg.Name, err)
	}

	d.SetId(id)

	return resourceEnvironmentRead(ctx, d, meta)
}

func resourceEnvironmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	envName := utils.String(d.Get(utils.TerraformResourceName))
	response, err := defaultConfig.GetEnvironment(envName)
	if err != nil {
		return diag.Errorf("getting environment %s errored with: %v", envName, err)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	return nil
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if d.HasChange(utils.TerraformResourcePipelines) || d.HasChange(utils.TerraformResourceEnvVar) {
		changes, err := getEnvChanges(d)
		if err != nil {
			return diag.Errorf("fetching changes errored with %v", err)
		}

		if changes.equal {
			return nil
		}

		cfg := gocd.Environment{
			Name:      utils.String(d.Get(utils.TerraformResourceName)),
			Pipelines: changes.pipelineChanges,
			EnvVars:   changes.envVarsChanges,
			ETAG:      utils.String(d.Get(utils.TerraformResourceEtag)),
		}

		_, err = defaultConfig.UpdateEnvironment(cfg)
		if err != nil {
			return diag.Errorf("updating environment %s errored with: %v", cfg.Name, err)
		}

		return resourceEnvironmentRead(ctx, d, meta)
	}

	log.Printf("nothing to update so skipping")

	return nil
}

func resourceEnvironmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if id := d.Id(); len(id) == 0 {
		return diag.Errorf("resource with the ID '%s' not found", id)
	}

	envName := utils.String(d.Get(utils.TerraformResourceName))

	err := defaultConfig.DeleteEnvironment(envName)
	if err != nil {
		return diag.Errorf("deleting environment %s errored with: %v", envName, err)
	}

	d.SetId("")

	return nil
}

func resourceEnvironmentImport(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	defaultConfig := meta.(gocd.GoCd)

	envName := utils.String(d.Id())
	response, err := defaultConfig.GetEnvironment(envName)
	if err != nil {
		return nil, fmt.Errorf("getting environment %s errored with: %w", envName, err)
	}

	if err = d.Set(utils.TerraformResourceName, envName); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceName)
	}

	if err = d.Set(utils.TerraformResourcePipelines, flattenPipelines(response.Pipelines)); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, err, utils.TerraformResourcePipelines)
	}

	flattenedEnvVars, err := utils.MapSlice(response.EnvVars)
	if err != nil {
		d.SetId("")

		return nil, fmt.Errorf("errored while flattening environment variable obtained: %w", err)
	}

	if err = d.Set(utils.TerraformResourceEnvVar, flattenedEnvVars); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceEnvVar)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	return []*schema.ResourceData{d}, nil
}

func getEnvironments(configs interface{}) ([]gocd.EnvVars, error) {
	var envVars []gocd.EnvVars
	envs := configs.(*schema.Set).List()
	if err := mapstructure.Decode(envs, &envVars); err != nil {
		return nil, err
	}

	return envVars, nil
}

func getPipelines(configs interface{}) []gocd.Pipeline {
	pipelines := make([]gocd.Pipeline, 0)
	for _, pipeline := range configs.([]interface{}) {
		pipelines = append(pipelines, gocd.Pipeline{Name: utils.String(pipeline)})
	}

	return pipelines
}

func getEnvChanges(d *schema.ResourceData) (environmentChanges, error) {
	var changes environmentChanges
	oldVars, newVars := d.GetChange(utils.TerraformResourceEnvVar)
	oldPipelines, newPipelines := d.GetChange(utils.TerraformResourcePipelines)

	envVars, err := getEnvironments(oldVars)
	if err != nil {
		return changes, fmt.Errorf("reading environment errored with %w", err)
	}

	changes.equal = true
	changes.envVarsChanges = envVars
	changes.pipelineChanges = getPipelines(oldPipelines)

	if !cmp.Equal(oldVars, newVars) {
		envVars, err = getEnvironments(newVars)
		if err != nil {
			return changes, fmt.Errorf("reading environment errored with %w", err)
		}
		changes.envVarsChanges = envVars
		changes.equal = false
	}

	if !cmp.Equal(oldPipelines, newPipelines) {
		changes.pipelineChanges = getPipelines(newPipelines)
		changes.equal = false
	}

	return changes, nil
}
