package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func dataSourceRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceRoleRead,
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
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				Description: "Type of the role. Use GoCD to create core role and plugin to create plugin role.",
			},
			"policy": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				Description: "Policy is fine-grained permissions attached to the users belonging to the current role.",
				Elem:        &schema.Schema{Type: schema.TypeMap},
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    false,
				Description: "The list of users belongs to the role.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"auth_config_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    false,
				Description: "The authorization configuration identifier.",
			},
			"properties": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ForceNew:    false,
				Description: "Attributes are used to describes the configuration for gocd role or plugin role.",
				Elem:        propertiesSchemaData(),
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Etag used to track the role",
			},
		},
	}
}

func datasourceRoleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()

	if len(id) == 0 {
		name := utils.String(d.Get(utils.TerraformResourceName))
		id = name
	}

	response, err := defaultConfig.GetRole(id)
	if err != nil {
		return diag.Errorf("getting role '%s' errored with: %v", id, err)
	}

	if err = d.Set(utils.TerraformResourceType, response.Type); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceType, err)
	}

	flattenedPolicy, err := utils.MapSlice(response.Policy)
	if err != nil {
		return diag.Errorf("errored while flattening Policy from the role obtained: %v", err)
	}

	if err = d.Set(utils.TerraformResourcePolicy, flattenedPolicy); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourcePolicy, err)
	}

	roleType := strings.ToLower(utils.String(d.Get(utils.TerraformResourceType)))
	switch roleType {
	case "plugin":
		if err = d.Set(utils.TerraformResourceAuthConfigID, response.Attributes.AuthConfigID); err != nil {
			return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceAuthConfigID, err)
		}

		flattenedProperties, err := utils.MapSlice(response.Attributes.Properties)
		if err != nil {
			d.SetId("")

			return diag.Errorf("errored while flattening properties of the role '%s' obtained with: %v", id, err)
		}

		if err = d.Set(utils.TerraformResourceProperties, flattenedProperties); err != nil {
			return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceProperties, err)
		}
	case "gocd":
		if err = d.Set(utils.TerraformResourceUsers, response.Attributes.Users); err != nil {
			return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceUsers, err)
		}
	default:
		return diag.Errorf("unknown role type '%s'", roleType)
	}

	d.SetId(id)

	return nil
}
