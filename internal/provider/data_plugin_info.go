package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func resourcePluginInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourcePluginInfoRead,
		Schema: map[string]*schema.Schema{
			"plugin_id": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The unique plugin identifier.",
			},
			"plugin_file_location": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The location where the plugin is installed.",
			},
			"bundled_plugin": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Indicates whether the plugin is bundled with GoCD.",
			},
			"status": {
				Type:        schema.TypeSet,
				Computed:    true,
				Optional:    true,
				Description: "The status of the plugin.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Description: "Status of the plugin. Can be one of active, invalid.",
							Computed:    true,
						},
					},
				},
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Etag used to track the plugin information",
			},
		},
	}
}

func datasourcePluginInfoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()

	if len(id) == 0 {
		pluginID := utils.String(d.Get(utils.TerraformPluginID))
		id = pluginID
	}

	pluginID := utils.String(d.Get(utils.TerraformPluginID))
	response, err := defaultConfig.GetPluginInfo(pluginID)
	if err != nil {
		return diag.Errorf("getting plugin information of '%s' errored with: %v", pluginID, err)
	}

	if err = d.Set(utils.TerraformPluginLocation, response.PluginFileLocation); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformPluginLocation)
	}

	if err = d.Set(utils.TerraformPluginBundled, response.BundledPlugin); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformPluginBundled)
	}

	if err = d.Set(utils.TerraformPluginStatus, flattenStatus(response)); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformPluginStatus)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	d.SetId(id)

	return nil
}

func flattenStatus(plugin gocd.Plugin) []map[string]string {
	return []map[string]string{
		{
			"status": plugin.Status.State,
		},
	}
}
