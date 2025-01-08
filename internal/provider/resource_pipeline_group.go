package provider

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func resourcePipelineGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelineGroupCreate,
		ReadContext:   resourcePipelineGroupRead,
		UpdateContext: resourcePipelineGroupUpdate,
		DeleteContext: resourcePipelineGroupDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "Name of the pipeline group to be created or updated.",
			},
			"pipelines": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    false,
				Description: "List of pipelines to be associated with pipeline group.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"authorization": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    false,
				Description: "The authorization configuration for the pipeline group.",
				Elem:        authConfigSchema(),
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Etag used to track the pipeline group.",
			},
		},
	}
}

func resourcePipelineGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !d.IsNewResource() {
		return nil
	}

	id := d.Id()

	if len(id) == 0 {
		resourceID := utils.String(d.Get(utils.TerraformResourceName))
		id = resourceID
	}

	cfg := gocd.PipelineGroup{
		Name:          id,
		Authorization: getPipelineGroupAuthorizationConfig(d.Get(utils.TerraformResourceAuthorization)),
	}

	out, err := json.MarshalIndent(cfg, " ", " ")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("PipelineGroup: %s", string(out))

	if err := defaultConfig.CreatePipelineGroup(cfg); err != nil {
		return diag.Errorf("creating pipeline group '%s' errored with %v", cfg.Name, err)
	}

	d.SetId(id)

	return resourcePipelineGroupRead(ctx, d, meta)
}

func resourcePipelineGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	name := utils.String(d.Get(utils.TerraformResourceName))
	response, err := defaultConfig.GetPipelineGroup(name)
	if err != nil {
		return diag.Errorf("getting pipeline group '%s' errored with: %v", name, err)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	return nil
}

func resourcePipelineGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !d.HasChange(utils.TerraformResourceAuthorization) && !d.HasChange(utils.TerraformResourcePipelines) {
		log.Printf("nothing to update so skipping")

		return nil
	}

	cfg := gocd.PipelineGroup{
		Name:          utils.String(d.Get(utils.TerraformResourceName)),
		Pipelines:     getPipelines(d.Get(utils.TerraformResourcePipelines)),
		Authorization: getPipelineGroupAuthorizationConfig(d.Get(utils.TerraformResourceAuthorization)),
		ETAG:          utils.String(d.Get(utils.TerraformResourceEtag)),
	}

	if _, err := defaultConfig.UpdatePipelineGroup(cfg); err != nil {
		return diag.Errorf("updating pipeline group '%s' errored with: %v", cfg.Name, err)
	}

	return resourcePipelineGroupRead(ctx, d, meta)
}

func resourcePipelineGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if id := d.Id(); len(id) == 0 {
		return diag.Errorf("resource with the ID '%s' not found", id)
	}

	profileID := utils.String(d.Get(utils.TerraformResourceName))

	err := defaultConfig.DeletePipelineGroup(profileID)
	if err != nil {
		return diag.Errorf("deleting pipeline group errored with: %v", err)
	}

	d.SetId("")

	return nil
}

func getPipelineGroupAuthorizationConfig(authConfig interface{}) gocd.PipelineGroupAuthorizationConfig {
	var flattenedView, flattenedAdmins, flattenedOperate map[string]interface{}

	var authorisationConfig gocd.PipelineGroupAuthorizationConfig
	flattenedAuthConfig := authConfig.(*schema.Set).List()[0].(map[string]interface{})

	if len(flattenedAuthConfig[utils.TerraformResourceView].(*schema.Set).List()) > 0 {
		flattenedView = flattenedAuthConfig[utils.TerraformResourceView].(*schema.Set).List()[0].(map[string]interface{})

		authorisationConfig.View = gocd.AuthorizationConfig{
			Roles: utils.GetSlice(flattenedView[utils.TerraformResourceRoles].([]interface{})),
			Users: utils.GetSlice(flattenedView[utils.TerraformResourceUsers].([]interface{})),
		}
	}

	if len(flattenedAuthConfig[utils.TerraformResourceOperate].(*schema.Set).List()) > 0 {
		flattenedOperate = flattenedAuthConfig[utils.TerraformResourceOperate].(*schema.Set).List()[0].(map[string]interface{})

		authorisationConfig.Operate = gocd.AuthorizationConfig{
			Roles: utils.GetSlice(flattenedOperate[utils.TerraformResourceRoles].([]interface{})),
			Users: utils.GetSlice(flattenedOperate[utils.TerraformResourceUsers].([]interface{})),
		}
	}

	if len(flattenedAuthConfig[utils.TerraformResourceAdmins].(*schema.Set).List()) > 0 {
		flattenedAdmins = flattenedAuthConfig[utils.TerraformResourceAdmins].(*schema.Set).List()[0].(map[string]interface{})

		authorisationConfig.Admins = gocd.AuthorizationConfig{
			Roles: utils.GetSlice(flattenedAdmins[utils.TerraformResourceRoles].([]interface{})),
			Users: utils.GetSlice(flattenedAdmins[utils.TerraformResourceUsers].([]interface{})),
		}
	}

	return authorisationConfig
}
