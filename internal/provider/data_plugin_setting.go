package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func dataSourcePluginsSetting() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluginsSettingRead,
		Schema: map[string]*schema.Schema{
			"plugin_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Computed:    false,
				Description: "The unique identifier of the plugin.",
			},
			"configuration": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "List of configuration required to configure the plugin settings.",
				Elem:        propertiesSchemaData(),
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Etag used to track the plugin settings.",
			},
		},
	}
}

func dataSourcePluginsSettingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()

	if len(id) == 0 {
		resourceID := utils.String(d.Get(utils.TerraformResourcePluginID))
		id = resourceID
	}

	pluginID := utils.String(d.Get(utils.TerraformResourcePluginID))

	response, err := defaultConfig.GetPluginSettings(pluginID)
	if err != nil {
		return diag.Errorf("getting cluster profile %s errored with: %v", pluginID, err)
	}

	flattenedConfiguration, err := utils.MapSlice(response.Configuration)
	if err != nil {
		d.SetId("")

		return diag.Errorf("errored while flattening Configuration obtained: %v", err)
	}

	if err = d.Set(utils.TerraformResourceConfiguration, flattenedConfiguration); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceConfiguration)
	}

	d.SetId(id)

	return nil
}
