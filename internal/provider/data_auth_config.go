package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func dataSourceAuthConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceAuthConfigRead,
		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The identifier of the authorization configuration.",
			},
			"plugin_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Required:    false,
				Description: "The plugin identifier of the authorization plugin.",
			},
			"allow_only_known_users_to_login": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Required:    false,
				Description: "Allow only those users to login who have explicitly been added by an administrator.",
			},
			"properties": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "The list of configuration properties that represent the configuration of this authorization configuration.",
				Elem:        propertiesSchemaData(),
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Etag used to track the authorization configuration.",
			},
		},
	}
}

func datasourceAuthConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()

	if len(id) == 0 {
		resourceID := utils.String(d.Get(utils.TerraformResourceProfileID))
		id = resourceID
	}

	profileID := utils.String(d.Get(utils.TerraformResourceProfileID))

	response, err := defaultConfig.GetAuthConfig(profileID)
	if err != nil {
		return diag.Errorf("getting authorization configuration %s errored with: %v", profileID, err)
	}

	if err = d.Set(utils.TerraformResourcePluginID, response.PluginID); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourcePluginID)
	}

	if err = d.Set(utils.TerraformResourceAllowKnownUser, response.AllowOnlyKnownUsers); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceAllowKnownUser)
	}

	flattenedProperties, err := utils.MapSlice(response.Properties)
	if err != nil {
		d.SetId("")

		return diag.Errorf("errored while flattening Properties obtained: %v", err)
	}

	if err = d.Set(utils.TerraformResourceProperties, flattenedProperties); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceProperties)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	d.SetId(id)

	return nil
}
