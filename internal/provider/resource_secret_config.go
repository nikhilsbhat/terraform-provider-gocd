package provider

import (
	"context"
	"log"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func resourceSecretConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecretConfigCreate,
		ReadContext:   resourceSecretConfigRead,
		DeleteContext: resourceSecretConfigDelete,
		UpdateContext: resourceSecretConfigUpdate,
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
				Optional:    true,
				Required:    false,
				ForceNew:    true,
				Description: "The identifier of the plugin to which current secret config belongs.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Required:    false,
				ForceNew:    true,
				Description: "The description for this secret config.",
			},
			"properties": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "The list of configuration properties that represent the configuration of this secret config.",
				Elem:        propertiesSchemaData(),
			},
			"rules": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Description: "The list of rules, which allows restricting the usage of the secret config. " +
					"Referring to the secret config from other parts of configuration is denied by default, " +
					"an explicit rule should be added to allow a specific resource to refer the secret config.",
				Elem: &schema.Schema{
					Type:        schema.TypeMap,
					Required:    false,
					Computed:    false,
					ForceNew:    false,
					Description: "Rule, which allows restricting the entities that the secret config can refer to.",
				},
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

func resourceSecretConfigCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !data.IsNewResource() {
		return nil
	}

	id := data.Id()

	if len(id) == 0 {
		resourceID := utils.String(data.Get(utils.TerraformResourceProfileID))
		id = resourceID
	}

	rules, err := flattenMapSlice(data.Get(utils.TerraformResourceRules))
	if err != nil {
		return diag.Errorf("reading rules errored with %v", err)
	}

	cfg := gocd.CommonConfig{
		ID:          utils.String(data.Get(utils.TerraformResourceProfileID)),
		PluginID:    utils.String(data.Get(utils.TerraformResourcePluginID)),
		Description: utils.String(data.Get(utils.TerraformResourceDescription)),
		Properties:  getPluginConfiguration(data.Get(utils.TerraformResourceProperties)),
		Rules:       rules,
	}

	if _, err = defaultConfig.CreateSecretConfig(cfg); err != nil {
		return diag.Errorf("creating secret config %s errored with %v", cfg.ID, err)
	}

	data.SetId(id)

	return resourceSecretConfigRead(ctx, data, meta)
}

func resourceSecretConfigUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if data.HasChange(utils.TerraformResourceProperties) ||
		data.HasChange(utils.TerraformResourceRules) {
		oldCfg, newCfg := data.GetChange(utils.TerraformResourceProperties)

		if cmp.Equal(oldCfg, newCfg) {
			return nil
		}

		rules, err := flattenMapSlice(data.Get(utils.TerraformResourceRules))
		if err != nil {
			return diag.Errorf("reading '%s' errored with %v", utils.TerraformResourceRules, err)
		}

		properties := getPluginConfiguration(data.Get(utils.TerraformResourceProperties))
		if err != nil {
			return diag.Errorf("reading '%s' errored with %v", utils.TerraformResourceProperties, err)
		}

		cfg := gocd.CommonConfig{
			ID:          utils.String(data.Get(utils.TerraformResourceProfileID)),
			PluginID:    utils.String(data.Get(utils.TerraformResourcePluginID)),
			Description: utils.String(data.Get(utils.TerraformResourceDescription)),
			Properties:  properties,
			Rules:       rules,
		}

		_, err = defaultConfig.UpdateSecretConfig(cfg)
		if err != nil {
			return diag.Errorf("updating secret config %s errored with: %v", cfg.ID, err)
		}

		return resourceSecretConfigRead(ctx, data, meta)
	}

	log.Printf("nothing to update so skipping")

	return nil
}

func resourceSecretConfigRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	profileID := utils.String(data.Get(utils.TerraformResourceProfileID))
	response, err := defaultConfig.GetSecretConfig(profileID)
	if err != nil {
		return diag.Errorf("getting secret config %s errored with: %v", profileID, err)
	}

	if err = data.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	return nil
}

func resourceSecretConfigDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if id := data.Id(); len(id) == 0 {
		return diag.Errorf("resource with the ID '%s' not found", id)
	}

	profileID := utils.String(data.Get(utils.TerraformResourceProfileID))

	err := defaultConfig.DeleteSecretConfig(profileID)
	if err != nil {
		return diag.Errorf("deleting secret config errored with: %v", err)
	}

	data.SetId("")

	return nil
}
