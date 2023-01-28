package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/mapstructure"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func resourceConfigRepository() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConfigRepoCreate,
		ReadContext:   resourceConfigRepoRead,
		DeleteContext: resourceConfigRepoDelete,
		UpdateContext: resourceConfigRepoUpdate,
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
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The plugin identifier of the cluster profile.",
			},
			"material":      materialSchema(),
			"configuration": propertiesSchemaResource(),
			"rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The list of rules, which allows restricting the entities that the config repo can refer to.",
				Elem: &schema.Schema{
					Type:        schema.TypeMap,
					Required:    false,
					Computed:    false,
					ForceNew:    false,
					Description: "Rule, which allows restricting the entities that the config repo can refer to.",
				},
			},
			"etag": {
				Type:        schema.TypeString,
				Required:    false,
				Computed:    true,
				ForceNew:    false,
				Description: "Etag used to track the plugin settings",
			},
		},
	}
}

func resourceConfigRepoCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	rules, err := getRules(d.Get(utils.TerraformResourceRules))
	if err != nil {
		return diag.Errorf("reading rules errored with %v", err)
	}
	material := getMaterials(d.Get(utils.TerraformResourceMaterial))
	if err != nil {
		return diag.Errorf("reading material errored with %v", err)
	}

	cfg := gocd.ConfigRepo{
		ID:            utils.String(d.Get(utils.TerraformResourceProfileID)),
		PluginID:      utils.String(d.Get(utils.TerraformPluginID)),
		Configuration: getPluginConfiguration(d.Get(utils.TerraformResourceConfiguration)),
		Rules:         rules,
		Material:      material,
	}

	if err = defaultConfig.CreateConfigRepo(cfg); err != nil {
		return diag.Errorf("creating config repo %s errored with %v", cfg.ID, err)
	}

	d.SetId(id)

	return resourceConfigRepoRead(ctx, d, meta)
}

func resourceConfigRepoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	profileID := utils.String(d.Get(utils.TerraformResourceProfileID))
	response, err := defaultConfig.GetConfigRepo(profileID)
	if err != nil {
		return diag.Errorf("getting config repo %s errored with: %v", profileID, err)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf("setting etag errored with %v", err)
	}

	return nil
}

func resourceConfigRepoUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if d.HasChange(utils.TerraformResourceMaterial) ||
		d.HasChange(utils.TerraformResourceRules) {
		oldCfg, newCfg := d.GetChange(utils.TerraformProperties)

		if cmp.Equal(oldCfg, newCfg) {
			return nil
		}

		rules, err := getRules(d.Get(utils.TerraformResourceRules))
		if err != nil {
			return diag.Errorf("reading rules errored with %v", err)
		}

		material := getMaterials(d.Get(utils.TerraformResourceMaterial))
		if err != nil {
			return diag.Errorf("reading material errored with %v", err)
		}

		cfg := gocd.ConfigRepo{
			ID:       utils.String(d.Get(utils.TerraformResourceProfileID)),
			PluginID: utils.String(d.Get(utils.TerraformPluginID)),
			Rules:    rules,
			Material: material,
			ETAG:     utils.String(d.Get(utils.TerraformResourceEtag)),
		}

		_, err = defaultConfig.UpdateConfigRepo(cfg)
		if err != nil {
			return diag.Errorf("updating config repo %s errored with: %v", cfg.ID, err)
		}

		return resourceConfigRepoRead(ctx, d, meta)
	}

	log.Printf("nothing to update so skipping")

	return nil
}

func resourceConfigRepoDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if id := d.Id(); len(id) == 0 {
		return diag.Errorf("resource with the ID '%s' not found", id)
	}

	profileID := utils.String(d.Get(utils.TerraformResourceProfileID))

	err := defaultConfig.DeleteConfigRepo(profileID)
	if err != nil {
		return diag.Errorf("deleting config repo errored with: %v", err)
	}

	d.SetId("")

	return nil
}

func getRules(configs interface{}) ([]map[string]interface{}, error) {
	var rules []map[string]interface{}
	if err := mapstructure.Decode(configs, &rules); err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return rules, nil
}

func getMaterials(configs interface{}) gocd.Material {
	var material gocd.Material
	flattenedMaterial := configs.(*schema.Set).List()[0].(map[string]interface{})
	flattenedAttr := flattenedMaterial[utils.TerraformResourceAttr].(*schema.Set).List()[0].(map[string]interface{})
	material = gocd.Material{
		Type:        utils.String(flattenedMaterial[utils.TerraformResourceType]),
		Fingerprint: utils.String(flattenedMaterial[utils.TerraformResourceFgPrint]),
		Attributes: gocd.Attribute{
			URL:                 utils.String(flattenedAttr[utils.TerraformResourceURL]),
			Username:            utils.String(flattenedAttr[utils.TerraformResourceUserName]),
			Password:            utils.String(flattenedAttr[utils.TerraformResourcePassword]),
			EncryptedPassword:   utils.String(flattenedAttr[utils.TerraformResourceEncryptPassword]),
			Branch:              utils.String(flattenedAttr[utils.TerraformResourceBranch]),
			AutoUpdate:          utils.Bool(flattenedAttr[utils.TerraformResourceAutoUpdate]),
			CheckExternals:      utils.Bool(flattenedAttr[utils.TerraformResourceCheck]),
			UseTickets:          utils.Bool(flattenedAttr[utils.TerraformResourceUseTickets]),
			View:                utils.String(flattenedAttr[utils.TerraformResourceView]),
			Port:                utils.String(flattenedAttr[utils.TerraformResourcePort]),
			ProjectPath:         utils.String(flattenedAttr[utils.TerraformResourceProjectPath]),
			Domain:              utils.String(flattenedAttr[utils.TerraformResourceDomain]),
			Ref:                 utils.String(flattenedAttr[utils.TerraformResourceRef]),
			Name:                utils.String(flattenedAttr[utils.TerraformResourceName]),
			Stage:               utils.String(flattenedAttr[utils.TerraformResourceStage]),
			Pipeline:            utils.String(flattenedAttr[utils.TerraformResourcePipeline]),
			IgnoreForScheduling: utils.Bool(flattenedAttr[utils.TerraformResourceIgnoreForScheduling]),
			Destination:         utils.String(flattenedAttr[utils.TerraformResourceDestination]),
			InvertFilter:        utils.Bool(flattenedAttr[utils.TerraformResourceInvertFilter]),
		},
	}

	return material
}
