package provider

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func resourceRole() *schema.Resource {
	propertySchema := propertiesSchemaResource()
	propertySchema.Required = false
	propertySchema.Optional = true

	return &schema.Resource{
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		DeleteContext: resourceRoleDelete,
		UpdateContext: resourceRoleUpdate,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The name of the role.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "Type of the role. Use GoCD to create core role and plugin to create plugin role.",
			},
			"policy": {
				Type:        schema.TypeList,
				Required:    true,
				Computed:    false,
				ForceNew:    false,
				Description: "Policy is fine-grained permissions attached to the users belonging to the current role.",
				Elem:        &schema.Schema{Type: schema.TypeMap},
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    false,
				Optional:    true,
				ForceNew:    false,
				Description: "The list of users belongs to the role.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"auth_config_id": {
				Type:        schema.TypeString,
				Computed:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "The authorization configuration identifier.",
			},
			"properties": propertySchema,
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Etag used to track the role",
			},
		},
	}
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !d.IsNewResource() {
		return nil
	}

	id := d.Id()

	if len(id) == 0 {
		name := utils.String(d.Get(utils.TerraformResourceName))
		id = name
	}

	roleCfg := gocd.Role{
		Name: id,
		Type: utils.String(d.Get(utils.TerraformResourceType)),
	}

	policy, err := flattenMapSlice(d.Get(utils.TerraformResourcePolicy))
	if err != nil {
		return diag.Errorf("flattening policy errored with %v", err)
	}

	roleCfg.Policy = policy

	roleType := strings.ToLower(roleCfg.Type)
	switch roleType {
	case "plugin":
		roleCfg.Attributes.AuthConfigID = utils.String(d.Get(utils.TerraformResourceAuthConfigID))
		roleCfg.Attributes.Properties = getPluginConfiguration(d.Get(utils.TerraformResourceProperties))
	case "gocd":
		roleCfg.Attributes = gocd.RoleAttribute{Users: utils.GetSlice(d.Get(utils.TerraformResourceUsers).([]interface{}))}
	default:
		return diag.Errorf("unknown role type '%s'", roleType)
	}

	if _, err = defaultConfig.CreateRole(roleCfg); err != nil {
		return diag.Errorf("creating role '%s' of type '%s' errored with %v", id, roleCfg.Type, err)
	}

	d.SetId(id)

	return resourceRoleRead(ctx, d, meta)
}

func resourceRoleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	name := utils.String(d.Get(utils.TerraformResourceName))
	response, err := defaultConfig.GetRole(name)
	if err != nil {
		return diag.Errorf("fetching role %s errored with: %v", name, err)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	return nil
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !d.HasChange(utils.TerraformResourceProperties) &&
		!d.HasChange(utils.TerraformResourcePolicy) &&
		!d.HasChange(utils.TerraformResourceUsers) {
		log.Printf("nothing to update so skipping")

		return nil
	}

	roleCfg := gocd.Role{
		Name: utils.String(d.Get(utils.TerraformResourceName)),
		Type: utils.String(d.Get(utils.TerraformResourceType)),
		ETAG: utils.String(d.Get(utils.TerraformResourceEtag)),
	}

	policy, err := flattenMapSlice(d.Get(utils.TerraformResourcePolicy))
	if err != nil {
		return diag.Errorf("flattening policy errored with %v", err)
	}

	roleCfg.Policy = policy

	roleType := strings.ToLower(roleCfg.Type)
	switch roleType {
	case "plugin":
		roleCfg.Attributes.AuthConfigID = utils.String(d.Get(utils.TerraformResourceAuthConfigID))
		roleCfg.Attributes.Properties = getPluginConfiguration(d.Get(utils.TerraformResourceProperties))
	case "gocd":
		roleAttr := gocd.RoleAttribute{Users: utils.GetSlice(d.Get(utils.TerraformResourceUsers).([]interface{}))}
		roleCfg.Attributes = roleAttr
	default:
		return diag.Errorf("unknown role type '%s'", roleType)
	}

	if _, err = defaultConfig.UpdateRole(roleCfg); err != nil {
		return diag.Errorf("updating role '%s' of type '%s' errored with %v", roleCfg.Name, roleCfg.Type, err)
	}

	return resourceRoleRead(ctx, d, meta)
}

func resourceRoleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()
	if len(d.Id()) == 0 {
		return diag.Errorf("resource with the ID '%s' not found", id)
	}

	name := utils.String(d.Get(utils.TerraformResourceName))

	err := defaultConfig.DeleteRole(name)
	if err != nil {
		return diag.Errorf("deleting role '%s' errored with: %v", name, err)
	}

	d.SetId("")

	return nil
}
