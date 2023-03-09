package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func dataSourceSecretConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecretConfigRead,
		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The identifier of the secret config.",
			},
			"plugin_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Required:    false,
				Description: "The identifier of the plugin to which current secret config belongs.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Required:    false,
				Description: "The description for this secret config.",
			},
			"properties": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "The list of configuration properties that represent the configuration of this secret config.",
				Elem:        propertiesSchemaData(),
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Description: "The list of rules, which allows restricting the usage of the secret config. " +
					"Referring to the secret config from other parts of configuration is denied by default, " +
					"an explicit rule should be added to allow a specific resource to refer the secret config.",
				Elem: rulesSchema(),
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Etag used to track the secret config",
			},
		},
	}
}

func dataSourceSecretConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()

	if len(id) == 0 {
		resourceID := utils.String(d.Get(utils.TerraformResourceProfileID))
		id = resourceID
	}

	profileID := utils.String(d.Get(utils.TerraformResourceProfileID))

	response, err := defaultConfig.GetSecretConfig(profileID)
	if err != nil {
		return diag.Errorf("getting secret config %s errored with: %v", profileID, err)
	}

	if err = d.Set(utils.TerraformResourcePluginID, response.PluginID); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourcePluginID, err)
	}

	if err = d.Set(utils.TerraformResourceDescription, response.Description); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceDescription, err)
	}

	flattenedProperties, err := utils.MapSlice(response.Properties)
	if err != nil {
		d.SetId("")

		return diag.Errorf("errored while flattening Properties obtained: %v", err)
	}

	if err = d.Set(utils.TerraformResourceProperties, flattenedProperties); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceProperties, err)
	}

	if err = d.Set(utils.TerraformResourceRules, response.Rules); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceRules, err)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	d.SetId(id)

	return nil
}
