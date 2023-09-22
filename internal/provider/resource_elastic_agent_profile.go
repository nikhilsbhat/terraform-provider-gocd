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

func resourceElasticAgentProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceElasticAgentProfileCreate,
		ReadContext:   resourceElasticAgentProfileRead,
		DeleteContext: resourceElasticAgentProfileDelete,
		UpdateContext: resourceElasticAgentProfileUpdate,
		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "the identifier of the elastic agent profile.",
			},
			"cluster_profile_id": {
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
				Description: "etag used to track the elastic agent profile configurations",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceElasticAgentProfileImport,
		},
	}
}

func resourceElasticAgentProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if d.IsNewResource() {
		id := d.Id()

		if len(id) == 0 {
			resourceID := utils.String(d.Get(utils.TerraformResourceProfileID))
			id = resourceID
		}

		cfg := gocd.CommonConfig{
			ID:               utils.String(d.Get(utils.TerraformResourceProfileID)),
			ClusterProfileID: utils.String(d.Get(utils.TerraformResourceClusterProfileID)),
			Properties:       getPluginConfiguration(d.Get(utils.TerraformResourceProperties)),
		}

		_, err := defaultConfig.CreateElasticAgentProfile(cfg)
		if err != nil {
			return diag.Errorf("creating elastic agent profile %s for cluster profile %s errored with %v", cfg.ID, cfg.ClusterProfileID, err)
		}

		d.SetId(id)

		return resourceElasticAgentProfileRead(ctx, d, meta)
	}

	return nil
}

func resourceElasticAgentProfileRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	profileID := utils.String(d.Get(utils.TerraformResourceProfileID))
	response, err := defaultConfig.GetElasticAgentProfile(profileID)
	if err != nil {
		return diag.Errorf("getting elastic agent profile configuration %s errored with: %v", profileID, err)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	return nil
}

func resourceElasticAgentProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !d.HasChange(utils.TerraformResourceProperties) {
		log.Printf("nothing to update so skipping")

		return nil
	}

	cfg := gocd.CommonConfig{
		ID:               utils.String(d.Get(utils.TerraformResourceProfileID)),
		ClusterProfileID: utils.String(d.Get(utils.TerraformResourceClusterProfileID)),
		Properties:       getPluginConfiguration(d.Get(utils.TerraformResourceProperties)),
		ETAG:             utils.String(d.Get(utils.TerraformResourceEtag)),
	}

	_, err := defaultConfig.UpdateElasticAgentProfile(cfg)
	if err != nil {
		return diag.Errorf("updating elastic agent profile %s errored with: %v", cfg.ID, err)
	}

	return resourceElasticAgentProfileRead(ctx, d, meta)
}

func resourceElasticAgentProfileDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if id := d.Id(); len(id) == 0 {
		return diag.Errorf("resource with the ID '%s' not found", id)
	}

	profileID := utils.String(d.Get(utils.TerraformResourceProfileID))

	err := defaultConfig.DeleteElasticAgentProfile(profileID)
	if err != nil {
		return diag.Errorf("deleting elastic agent profile %s errored with: %v", profileID, err)
	}

	d.SetId("")

	return nil
}

func resourceElasticAgentProfileImport(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	defaultConfig := meta.(gocd.GoCd)

	profileID := utils.String(d.Id())
	response, err := defaultConfig.GetElasticAgentProfile(profileID)
	if err != nil {
		return nil, fmt.Errorf("getting elastic agent profile configuration %s errored with: %w", profileID, err)
	}

	if err = d.Set(utils.TerraformResourceProfileID, profileID); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceStoreID)
	}

	if err = d.Set(utils.TerraformResourceClusterProfileID, response.ClusterProfileID); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceClusterProfileID)
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
