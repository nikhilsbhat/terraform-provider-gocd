// ----------------------------------------------------------------------------
//
//	***     TERRAGEN GENERATED CODE    ***    TERRAGEN GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//	This file was auto generated by Terragen.
//	This autogenerated code has to be enhanced further to make it fully working terraform-provider.
//
//	Get more information on how terragen works.
//	https://github.com/nikhilsbhat/terragen
//
// ----------------------------------------------------------------------------
//
//nolint:gocritic
package gocd

import (
	"context"
	"log"

	"github.com/google/go-cmp/cmp"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePlugins() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePluginsCreate,
		ReadContext:   resourcePluginsRead,
		DeleteContext: resourcePluginsDelete,
		UpdateContext: resourcePluginsUpdate,
		Schema: map[string]*schema.Schema{
			"plugin_id": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "id of the GoCD plugin to which the settings to be applied",
			},
			"plugin_configurations": {
				Type:        schema.TypeSet,
				Required:    true,
				Computed:    false,
				Description: "list of configurations to be applied to GoCD plugin",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Computed:    false,
							Description: "the name of the property key.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    false,
							ForceNew:    true,
							Description: "The value of the property",
						},
						"encrypted_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    false,
							ForceNew:    true,
							Description: "The encrypted value of the property",
						},
					},
				},
			},
			"etag": {
				Type:        schema.TypeString,
				Required:    false,
				Computed:    true,
				ForceNew:    false,
				Description: "etag used to track the plugin settings",
			},
		},
	}
}

func resourcePluginsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !d.IsNewResource() {
		return nil
	}

	id := d.Id()

	if len(id) == 0 {
		newID, err := utils.GetRandomID()
		if err != nil {
			d.SetId("")

			return diag.Errorf("errored while fetching randomID %v", err)
		}
		id = newID
	}

	pluginSettings := gocd.PluginSettings{
		ID:            utils.String(d.Get(utils.TerraformPluginID)),
		Configuration: getPluginConfiguration(d.Get(utils.TerraformResourcePluginConfiguration)),
	}

	_, err := defaultConfig.CreatePluginSettings(pluginSettings)
	if err != nil {
		return diag.Errorf("applying plugin setting errored with %v", err)
	}

	d.SetId(id)

	return resourcePluginsRead(ctx, d, meta)
}

func resourcePluginsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	response, err := defaultConfig.GetPluginSettings(utils.String(d.Get(utils.TerraformPluginID)))
	if err != nil {
		return diag.Errorf("getting plugin configuration errored with: %v", err)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf("setting etag errored with %v", err)
	}

	return nil
}

func resourcePluginsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if d.HasChange(utils.TerraformResourcePluginConfiguration) {
		oldCfg, newCfg := d.GetChange(utils.TerraformResourcePluginConfiguration)
		if !cmp.Equal(oldCfg, newCfg) {
			pluginSettings := gocd.PluginSettings{
				ID:            utils.String(d.Get(utils.TerraformPluginID)),
				Configuration: getPluginConfiguration(newCfg),
				ETAG:          utils.String(d.Get(utils.TerraformResourceEtag)),
			}

			_, err := defaultConfig.UpdatePluginSettings(pluginSettings)
			if err != nil {
				return diag.Errorf("updating plugin configuration errored with: %v", err)
			}

			return resourcePluginsRead(ctx, d, meta)
		}
	}

	log.Printf("nothing to update so skipping")

	return nil
}

func resourcePluginsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if id := d.Id(); len(id) == 0 {
		return diag.Errorf("resource with the ID %s not found", id)
	}

	pluginSettings := gocd.PluginSettings{
		ID:            utils.String(d.Get(utils.TerraformPluginID)),
		Configuration: []gocd.PluginConfiguration{},
		ETAG:          utils.String(d.Get(utils.TerraformResourceEtag)),
	}

	_, err := defaultConfig.UpdatePluginSettings(pluginSettings)
	if err != nil {
		return diag.Errorf("updating plugin configuration errored with: %v", err)
	}

	d.SetId("")

	return nil
}

func getPluginConfiguration(configs interface{}) []gocd.PluginConfiguration {
	pluginsConfigurations := make([]gocd.PluginConfiguration, 0)
	for _, config := range configs.(*schema.Set).List() {
		v := config.(map[string]interface{})
		pluginsConfigurations = append(pluginsConfigurations, gocd.PluginConfiguration{
			Key:            utils.String(v[utils.TerraformResourceKey]),
			Value:          utils.String(v[utils.TerraformResourceValue]),
			EncryptedValue: utils.String(v[utils.TerraformResourceENCValue]),
			IsSecure:       utils.Bool(v[utils.TerraformResourceIsSecure]),
		})
	}

	return pluginsConfigurations
}
