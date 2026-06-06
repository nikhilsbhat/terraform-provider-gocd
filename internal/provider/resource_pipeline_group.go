package provider

import (
	"context"
	"encoding/json"
	"fmt"
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
		Importer: &schema.ResourceImporter{
			StateContext: resourcePipelineGroupImport,
		},
	}
}

func resourcePipelineGroupCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

func resourcePipelineGroupRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	name := utils.String(d.Get(utils.TerraformResourceName))

	response, err := defaultConfig.GetPipelineGroup(name)
	if err != nil {
		return diag.Errorf("getting pipeline group '%s' errored with: %v", name, err)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	if err = d.Set(utils.TerraformResourcePipelines, flattenPipelines(response.Pipelines)); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourcePipelines)
	}

	if err = d.Set(utils.TerraformResourceAuthorization, flattenPipelineGroupAuthorizationConfig(response, d.Get(utils.TerraformResourceAuthorization))); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceAuthorization)
	}

	return nil
}

func resourcePipelineGroupUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

func resourcePipelineGroupDelete(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

func resourcePipelineGroupImport(_ context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	defaultConfig := meta.(gocd.GoCd)

	pipelineGroupName := utils.String(d.Id())

	response, err := defaultConfig.GetPipelineGroup(pipelineGroupName)
	if err != nil {
		return nil, fmt.Errorf("getting pipeline group %s errored with: %w", pipelineGroupName, err)
	}

	if err = d.Set(utils.TerraformResourceName, pipelineGroupName); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceName)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	if err = d.Set(utils.TerraformResourcePipelines, flattenPipelines(response.Pipelines)); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, err, utils.TerraformResourcePipelines)
	}

	flattenedAuthVar := flattenPipelineGroupAuthorizationConfig(response, d.Get(utils.TerraformResourceAuthorization))

	if err = d.Set(utils.TerraformResourceAuthorization, flattenedAuthVar); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceAuthorization)
	}

	return []*schema.ResourceData{d}, nil
}

func flattenPipelineGroupAuthorizationConfig(pipelineGroup gocd.PipelineGroup, configuredAuth any) []any {
	auth := pipelineGroup.Authorization
	configuredAuthMap := getConfiguredPipelineGroupAuthorization(configuredAuth)
	authConfig := make(map[string]any)

	addSection := func(key string, section gocd.AuthorizationConfig) {
		flattenedSection := flattenAuthorizationConfig(key, section, configuredAuthMap)
		if len(flattenedSection) > 0 {
			authConfig[key] = []any{flattenedSection}
		}
	}

	addSection(utils.TerraformResourceView, auth.View)
	addSection(utils.TerraformResourceOperate, auth.Operate)
	addSection(utils.TerraformResourceAdmins, auth.Admins)

	if len(authConfig) == 0 {
		return nil
	}

	return []any{authConfig}
}

func flattenAuthorizationConfig(sectionName string, auth gocd.AuthorizationConfig, configuredAuth map[string]any) map[string]any {
	authSection := make(map[string]any)

	if len(auth.Users) > 0 || authConfigFieldWasSet(configuredAuth, sectionName, utils.TerraformResourceUsers) {
		authSection[utils.TerraformResourceUsers] = auth.Users
	}

	if len(auth.Roles) > 0 || authConfigFieldWasSet(configuredAuth, sectionName, utils.TerraformResourceRoles) {
		authSection[utils.TerraformResourceRoles] = auth.Roles
	}

	return authSection
}

func getConfiguredPipelineGroupAuthorization(configuredAuth any) map[string]any {
	authSet, ok := configuredAuth.(*schema.Set)
	if !ok || authSet.Len() == 0 {
		return nil
	}

	configuredAuthMap, ok := authSet.List()[0].(map[string]any)
	if !ok {
		return nil
	}

	return configuredAuthMap
}

func authConfigFieldWasSet(configuredAuth map[string]any, sectionName, fieldName string) bool {
	if configuredAuth == nil {
		return false
	}

	sectionSet, ok := configuredAuth[sectionName].(*schema.Set)
	if !ok || sectionSet.Len() == 0 {
		return false
	}

	section, ok := sectionSet.List()[0].(map[string]any)
	if !ok {
		return false
	}

	_, ok = section[fieldName]

	return ok
}

func getPipelineGroupAuthorizationConfig(authConfig any) gocd.PipelineGroupAuthorizationConfig {
	var flattenedView, flattenedAdmins, flattenedOperate map[string]any

	var authorisationConfig gocd.PipelineGroupAuthorizationConfig

	authSet := authConfig.(*schema.Set)
	if authSet.Len() == 0 {
		return authorisationConfig
	}

	flattenedAuthConfig := authSet.List()[0].(map[string]any)

	if len(flattenedAuthConfig[utils.TerraformResourceView].(*schema.Set).List()) > 0 {
		flattenedView = flattenedAuthConfig[utils.TerraformResourceView].(*schema.Set).List()[0].(map[string]any)

		authorisationConfig.View = gocd.AuthorizationConfig{
			Roles: utils.GetSlice(flattenedView[utils.TerraformResourceRoles].([]any)),
			Users: utils.GetSlice(flattenedView[utils.TerraformResourceUsers].([]any)),
		}
	}

	if len(flattenedAuthConfig[utils.TerraformResourceOperate].(*schema.Set).List()) > 0 {
		flattenedOperate = flattenedAuthConfig[utils.TerraformResourceOperate].(*schema.Set).List()[0].(map[string]any)

		authorisationConfig.Operate = gocd.AuthorizationConfig{
			Roles: utils.GetSlice(flattenedOperate[utils.TerraformResourceRoles].([]any)),
			Users: utils.GetSlice(flattenedOperate[utils.TerraformResourceUsers].([]any)),
		}
	}

	if len(flattenedAuthConfig[utils.TerraformResourceAdmins].(*schema.Set).List()) > 0 {
		flattenedAdmins = flattenedAuthConfig[utils.TerraformResourceAdmins].(*schema.Set).List()[0].(map[string]any)

		authorisationConfig.Admins = gocd.AuthorizationConfig{
			Roles: utils.GetSlice(flattenedAdmins[utils.TerraformResourceRoles].([]any)),
			Users: utils.GetSlice(flattenedAdmins[utils.TerraformResourceUsers].([]any)),
		}
	}

	return authorisationConfig
}
