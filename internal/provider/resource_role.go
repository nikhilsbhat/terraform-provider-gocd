package provider

import (
	"context"
	"fmt"
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
		Importer: &schema.ResourceImporter{
			StateContext: resourceRoleImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The name of the role.",
			},
			"system_admin": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Computed:    false,
				Description: "Enable if the role should be set as admin",
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

	if err = updateAdmin(defaultConfig, d); err != nil {
		return diag.Errorf("%v", err)
	}

	d.SetId(id)

	return resourceRoleRead(ctx, d, meta)
}

func resourceRoleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	name := d.Id()
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
		!d.HasChange(utils.TerraformResourceUsers) &&
		!d.HasChange(utils.TerraformResourceSystemAdmin) {
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

	if err = updateAdmin(defaultConfig, d); err != nil {
		return diag.Errorf("%v", err)
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

func resourceRoleImport(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	defaultConfig := meta.(gocd.GoCd)

	roleName := utils.String(d.Id())
	response, err := defaultConfig.GetRole(roleName)
	if err != nil {
		return nil, fmt.Errorf("getting pipeline group %s errored with: %w", roleName, err)
	}

	if err = d.Set(utils.TerraformResourceName, roleName); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceName)
	}

	if err = d.Set(utils.TerraformResourceType, response.Type); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceType)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	if len(response.Policy) > 0 {
		policies := make([]interface{}, len(response.Policy))
		for index, policy := range response.Policy {
			policies[index] = map[string]interface{}{}
			for policyKey, policyValue := range policy {
				policies[index].(map[string]interface{})[policyKey] = policyValue
			}
		}

		if err = d.Set(utils.TerraformResourcePolicy, policies); err != nil {
			return nil, fmt.Errorf(settingAttrErrorTmp, utils.TerraformResourcePolicy, err)
		}
	}

	roleType := strings.ToLower(response.Type)
	switch roleType {
	case "plugin":
		if err = d.Set(utils.TerraformResourceAuthConfigID, response.Attributes.AuthConfigID); err != nil {
			return nil, fmt.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceAuthConfigID)
		}

		flattenedProperties, err := utils.MapSlice(response.Attributes.Properties)
		if err != nil {
			d.SetId("")

			return nil, fmt.Errorf("errored while flattening properties variable obtained: %w", err)
		}

		if err = d.Set(utils.TerraformResourceProperties, flattenedProperties); err != nil {
			return nil, fmt.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceProperties)
		}
	case "gocd":
		if err = d.Set(utils.TerraformResourceUsers, response.Attributes.Users); err != nil {
			return nil, fmt.Errorf(settingAttrErrorTmp, utils.TerraformResourceUsers, err)
		}
	default:
		return nil, fmt.Errorf("unknown role type '%s'", roleType)
	}

	return []*schema.ResourceData{d}, nil
}

// Ensures the role is added as a system admin in GoCD.
func updateAdmin(defaultConfig gocd.GoCd, d *schema.ResourceData) error {
	resourceName := utils.String(d.Get(utils.TerraformResourceName))
	isAdmin := utils.Bool(d.Get(utils.TerraformResourceSystemAdmin))

	admins, err := defaultConfig.GetSystemAdmins()
	if err != nil {
		return fmt.Errorf("fetching system admins errored with %w", err)
	}

	isAlreadyAdmin := utils.Contains(admins.Roles, resourceName)

	if isAlreadyAdmin == isAdmin {
		return nil
	}

	addNRemove := gocd.AddRemoves{}
	if isAdmin {
		log.Printf("Adding role '%s' to system admins", resourceName)
		addNRemove.Add = []string{resourceName}
	} else {
		log.Printf("Removing role '%s' from system admins", resourceName)
		addNRemove.Remove = []string{resourceName}
	}

	operationOptions := gocd.Operations{Roles: addNRemove}

	if _, err = defaultConfig.UpdateSystemAdminsBulk(operationOptions); err != nil {
		return fmt.Errorf("updating system admins bulk errored with %w", err)
	}

	return nil
}
