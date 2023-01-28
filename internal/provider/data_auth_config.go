package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func datasourceAuthConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceAuthConfigRead,
		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The identifier of the elastic agent profile.",
			},
			"plugin_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Required:    false,
				Description: "The plugin identifier of the cluster profile.",
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
				Description: "the list of configuration properties that represent the configuration of this profile.",
				Elem:        propertiesSchemaData(),
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Etag used to track the authorisation configuration.",
			},
		},
	}
}

func datasourceAuthConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	profileID := utils.String(d.Get(utils.TerraformResourceProfileID))

	response, err := defaultConfig.GetAuthConfig(profileID)
	if err != nil {
		return diag.Errorf("getting authorization configuration %s errored with: %v", profileID, err)
	}

	if err = d.Set(utils.TerraformPluginID, response.PluginID); err != nil {
		return diag.Errorf("setting '%s' errored with %v", err, utils.TerraformPluginID)
	}

	if err = d.Set(utils.TerraformResourceAllowKnownUser, response.AllowOnlyKnownUsers); err != nil {
		return diag.Errorf("setting '%s' errored with %v", err, utils.TerraformResourceAllowKnownUser)
	}

	flattenedProperties, err := utils.MapSlice(response.Properties)
	if err != nil {
		d.SetId("")

		return diag.Errorf("errored while flattening Properties obtained: %v", err)
	}

	if err = d.Set(utils.TerraformProperties, flattenedProperties); err != nil {
		return diag.Errorf("setting '%s' errored with %v", err, utils.TerraformProperties)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf("setting etag errored with %v", err)
	}

	d.SetId(id)

	return nil
}
