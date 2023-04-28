package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func resourceArtifactStore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArtifactStoreCreate,
		ReadContext:   resourceArtifactStoreRead,
		DeleteContext: resourceArtifactStoreDelete,
		UpdateContext: resourceArtifactStoreUpdate,
		Schema: map[string]*schema.Schema{
			"store_id": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The identifier of the artifact store.",
			},
			"plugin_id": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The plugin identifier of the artifact plugin.",
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
	}
}

func resourceArtifactStoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !d.IsNewResource() {
		return nil
	}

	id := d.Id()

	if len(id) == 0 {
		resourceID := utils.String(d.Get(utils.TerraformResourceStoreID))
		id = resourceID
	}

	cfg := gocd.CommonConfig{
		ID:         id,
		PluginID:   utils.String(d.Get(utils.TerraformResourcePluginID)),
		Properties: getPluginConfiguration(d.Get(utils.TerraformResourceProperties)),
	}

	if _, err := defaultConfig.CreateArtifactStore(cfg); err != nil {
		return diag.Errorf("creating artifact store '%s' for plugin '%s' errored with %v", cfg.ID, cfg.PluginID, err)
	}

	d.SetId(id)

	return resourceArtifactStoreRead(ctx, d, meta)
}

func resourceArtifactStoreRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	storeID := utils.String(d.Get(utils.TerraformResourceStoreID))
	response, err := defaultConfig.GetArtifactStore(storeID)
	if err != nil {
		return diag.Errorf("getting artifact store configuration '%s' errored with: %v", storeID, err)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	return nil
}

func resourceArtifactStoreUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !d.HasChange(utils.TerraformResourceProperties) {
		log.Printf("nothing to update so skipping")

		return nil
	}

	cfg := gocd.CommonConfig{
		ID:         utils.String(d.Get(utils.TerraformResourceStoreID)),
		PluginID:   utils.String(d.Get(utils.TerraformResourcePluginID)),
		Properties: getPluginConfiguration(d.Get(utils.TerraformResourceProperties)),
		ETAG:       utils.String(d.Get(utils.TerraformResourceEtag)),
	}

	_, err := defaultConfig.UpdateArtifactStore(cfg)
	if err != nil {
		return diag.Errorf("updating artifact store config '%s' errored with: %v", cfg.ID, err)
	}

	return resourceArtifactStoreRead(ctx, d, meta)
}

func resourceArtifactStoreDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()
	if len(d.Id()) == 0 {
		return diag.Errorf("resource with the ID '%s' not found", id)
	}

	storeID := utils.String(d.Get(utils.TerraformResourceStoreID))

	err := defaultConfig.DeleteArtifactStore(storeID)
	if err != nil {
		return diag.Errorf("deleting artifact store '%s' errored with: %v", storeID, err)
	}

	d.SetId("")

	return nil
}
