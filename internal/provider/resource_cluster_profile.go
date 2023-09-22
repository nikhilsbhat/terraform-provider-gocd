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
package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func resourceClusterProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterProfileCreate,
		ReadContext:   resourceClusterProfileRead,
		DeleteContext: resourceClusterProfileDelete,
		UpdateContext: resourceClusterProfileUpdate,
		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "the identifier of the cluster profile.",
			},
			"plugin_id": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "the plugin identifier of the cluster profile.",
			},
			"properties": propertiesSchemaResource(),
			"etag": {
				Type:        schema.TypeString,
				Required:    false,
				Computed:    true,
				ForceNew:    false,
				Description: "etag used to track the plugin settings",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceClusterProfileImport,
		},
	}
}

func resourceClusterProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !d.IsNewResource() {
		return nil
	}

	id := d.Id()

	if len(id) == 0 {
		resourceID := utils.String(d.Get(utils.TerraformResourceProfileID))
		id = resourceID
	}

	cfg := gocd.CommonConfig{
		ID:         utils.String(d.Get(utils.TerraformResourceProfileID)),
		PluginID:   utils.String(d.Get(utils.TerraformResourcePluginID)),
		Properties: getPluginConfiguration(d.Get(utils.TerraformResourceProperties)),
	}

	_, err := defaultConfig.CreateClusterProfile(cfg)
	if err != nil {
		return diag.Errorf("creating cluster profile %s setting for plugin %s errored with %v", cfg.ID, cfg.PluginID, err)
	}

	d.SetId(id)

	return resourceClusterProfileRead(ctx, d, meta)
}

func resourceClusterProfileRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	profileID := utils.String(d.Get(utils.TerraformResourceProfileID))
	response, err := defaultConfig.GetClusterProfile(profileID)
	if err != nil {
		return diag.Errorf("getting cluster profile configuration %s errored with: %v", profileID, err)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	return nil
}

func resourceClusterProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !d.HasChange(utils.TerraformResourceProperties) {
		log.Printf("nothing to update so skipping")

		return nil
	}

	cfg := gocd.CommonConfig{
		ID:         utils.String(d.Get(utils.TerraformResourceProfileID)),
		PluginID:   utils.String(d.Get(utils.TerraformResourcePluginID)),
		Properties: getPluginConfiguration(d.Get(utils.TerraformResourceProperties)),
		ETAG:       utils.String(d.Get(utils.TerraformResourceEtag)),
	}

	_, err := defaultConfig.UpdateClusterProfile(cfg)
	if err != nil {
		return diag.Errorf("updating cluster profile %s errored with: %v", cfg.ID, err)
	}

	return resourceClusterProfileRead(ctx, d, meta)
}

func resourceClusterProfileDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()
	if len(d.Id()) == 0 {
		return diag.Errorf("resource with the ID '%s' not found", id)
	}

	profileID := utils.String(d.Get(utils.TerraformResourceProfileID))

	err := defaultConfig.DeleteClusterProfile(profileID)
	if err != nil {
		return diag.Errorf("deleting cluster profile %s errored with: %v", profileID, err)
	}

	d.SetId("")

	return nil
}

func resourceClusterProfileImport(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	defaultConfig := meta.(gocd.GoCd)

	profileID := utils.String(d.Id())
	response, err := defaultConfig.GetClusterProfile(profileID)
	if err != nil {
		return nil, fmt.Errorf("getting cluster profile configuration %s errored with: %w", profileID, err)
	}

	if err = d.Set(utils.TerraformResourceProfileID, profileID); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceStoreID)
	}

	if err = d.Set(utils.TerraformResourcePluginID, response.PluginID); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, err, utils.TerraformResourcePluginID)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	flattenedProperties, err := utils.MapSlice(response.Properties)
	if err != nil {
		d.SetId("")

		return nil, fmt.Errorf("errored while flattening artifact store properties obtained: %w", err)
	}

	if err = d.Set(utils.TerraformResourceProperties, flattenedProperties); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceProperties)
	}

	return []*schema.ResourceData{d}, nil
}
