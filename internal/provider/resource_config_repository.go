package provider

import (
	"context"
	"fmt"
	"log"
	"reflect"

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
				Description: "The identifier of the config repository.",
			},
			"plugin_id": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The name of the config repo plugin.",
			},
			"material": materialSchema(),
			"configuration": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    false,
				Description: "The list of configuration properties that represent the configuration of config repositories.",
				Elem:        propertiesSchemaResource().Elem,
			},
			"rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    false,
				ForceNew:    false,
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
				Description: "Etag used to track the config repository.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceConfigRepoImport,
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
		resourceID := utils.String(d.Get(utils.TerraformResourceProfileID))
		id = resourceID
	}

	rules, err := flattenMapSlice(d.Get(utils.TerraformResourceRules))
	if err != nil {
		return diag.Errorf("reading rules errored with %v", err)
	}

	material, err := getMaterials(d.Get(utils.TerraformResourceMaterial))
	if err != nil {
		return diag.Errorf("failed to parse material: %v", err)
	}

	cfg := gocd.ConfigRepo{
		ID:            utils.String(d.Get(utils.TerraformResourceProfileID)),
		PluginID:      utils.String(d.Get(utils.TerraformResourcePluginID)),
		Configuration: getPluginConfiguration(d.Get(utils.TerraformResourceConfiguration)),
		Rules:         rules,
		Material:      material,
	}

	if err = defaultConfig.CreateConfigRepo(cfg); err != nil {
		return diag.Errorf("creating config repo %s errored with %v", cfg.ID, err)
	}

	d.SetId(id)

	return dataSourceConfigRepositoryRead(ctx, d, meta)
}

func resourceConfigRepoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	profileID := utils.String(d.Get(utils.TerraformResourceProfileID))
	response, err := defaultConfig.GetConfigRepo(profileID)
	if err != nil {
		return diag.Errorf("getting config repo %s errored with: %v", profileID, err)
	}

	if err = d.Set(utils.TerraformResourcePluginID, response.PluginID); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourcePluginID, err)
	}

	flattened := flattenMaterial(response.Material)

	if err = d.Set("material", flattened); err != nil {
		return diag.Errorf("setting material errored with: %v", err)
	}

	flattenedConfiguration, err := utils.MapSlice(response.Configuration)
	if err != nil {
		return diag.Errorf("errored while flattening Configuration obtained: %v", err)
	}

	if err = d.Set(utils.TerraformResourceConfiguration, flattenedConfiguration); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceConfiguration, err)
	}

	if err = d.Set(utils.TerraformResourceRules, response.Rules); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceRules, err)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	return nil
}

func resourceConfigRepoUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !d.HasChange(utils.TerraformResourceMaterial) &&
		!d.HasChange(utils.TerraformResourceRules) &&
		!d.HasChange(utils.TerraformResourceConfiguration) {
		return nil
	}

	rules, err := flattenMapSlice(d.Get(utils.TerraformResourceRules))
	if err != nil {
		return diag.Errorf("reading rules errored with %v", err)
	}

	material, err := getMaterials(d.Get(utils.TerraformResourceMaterial))
	if err != nil {
		return diag.Errorf("failed to parse material: %v", err)
	}

	cfg := gocd.ConfigRepo{
		ID:            utils.String(d.Get(utils.TerraformResourceProfileID)),
		PluginID:      utils.String(d.Get(utils.TerraformResourcePluginID)),
		Rules:         rules,
		Material:      material,
		Configuration: getPluginConfiguration(d.Get(utils.TerraformResourceConfiguration)),
		ETAG:          utils.String(d.Get(utils.TerraformResourceEtag)),
	}

	if _, err = defaultConfig.UpdateConfigRepo(cfg); err != nil {
		return diag.Errorf("updating config repo %s errored with: %v", cfg.ID, err)
	}

	return dataSourceConfigRepositoryRead(ctx, d, meta)
}

func resourceConfigRepoDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceConfigRepoImport(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	defaultConfig := meta.(gocd.GoCd)

	profileID := utils.String(d.Id())
	response, err := defaultConfig.GetConfigRepo(profileID)
	if err != nil {
		return nil, fmt.Errorf("getting config repo %s errored with: %w", profileID, err)
	}

	if err = d.Set(utils.TerraformResourceProfileID, profileID); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceStoreID)
	}

	if err = d.Set(utils.TerraformResourcePluginID, response.PluginID); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, utils.TerraformResourcePluginID, err)
	}

	flattened := flattenMaterial(response.Material)

	if err = d.Set("material", flattened); err != nil {
		return nil, fmt.Errorf("setting material errored with: %w", err)
	}

	flattenedConfiguration, err := utils.MapSlice(response.Configuration)
	if err != nil {
		return nil, fmt.Errorf("errored while flattening Configuration obtained: %w", err)
	}

	if err = d.Set(utils.TerraformResourceConfiguration, flattenedConfiguration); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, utils.TerraformResourceConfiguration, err)
	}

	if err = d.Set(utils.TerraformResourceRules, response.Rules); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, utils.TerraformResourceRules, err)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return nil, fmt.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	return []*schema.ResourceData{d}, nil
}

func flattenMapSlice(configs interface{}) ([]map[string]string, error) {
	var rules []map[string]string
	if err := mapstructure.Decode(configs, &rules); err != nil {
		return nil, err
	}

	return rules, nil
}

func getMaterials(configs interface{}) (gocd.Material, error) {
	if set, ok := configs.(*schema.Set); ok {
		configs = set.List()
	}

	materialList, ok := configs.([]interface{})
	if !ok {
		return gocd.Material{}, fmt.Errorf("expected []interface{} or *schema.Set, got %T", configs)
	}

	if len(materialList) == 0 {
		return gocd.Material{}, nil
	}

	flattenedMaterial, ok := materialList[0].(map[string]interface{})
	if !ok {
		return gocd.Material{}, fmt.Errorf("expected map[string]interface{} for material, got %T", materialList[0])
	}

	attrRaw := flattenedMaterial["attributes"]
	if attrRaw == nil {
		return gocd.Material{}, nil
	}

	var attrList []interface{}
	if set, ok := attrRaw.(*schema.Set); ok {
		attrList = set.List()
	} else if list, ok := attrRaw.([]interface{}); ok {
		attrList = list
	} else {
		return gocd.Material{}, fmt.Errorf("expected []interface{} or *schema.Set for attributes, got %T", attrRaw)
	}

	var flattenedAttr map[string]interface{}
	if len(attrList) > 0 {
		flattenedAttr, ok = attrList[0].(map[string]interface{})
		if !ok {
			return gocd.Material{}, fmt.Errorf("expected map[string]interface{} for attributes[0], got %T", attrList[0])
		}
	}

	material := gocd.Material{
		Type:        utils.String(flattenedMaterial["type"]),
		Fingerprint: utils.String(flattenedMaterial["fingerprint"]),
	}

	if flattenedAttr != nil {
		log.Printf("TerraformResourceAutoUpdate: %v", flattenedAttr[utils.TerraformResourceAutoUpdate])
		log.Printf("TerraformResourceAutoUpdate Bool: %v", utils.Bool(flattenedAttr[utils.TerraformResourceAutoUpdate]))
		log.Printf("TerraformResourceAutoUpdate Type: %T", utils.Bool(flattenedAttr[utils.TerraformResourceAutoUpdate]))

		material.Attributes = gocd.Attribute{
			URL:                 utils.String(flattenedAttr[utils.TerraformResourceURL]),
			Username:            utils.String(flattenedAttr[utils.TerraformResourceUserName]),
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
		}
	}

	return material, nil
}

func flattenMaterial(material gocd.Material) []interface{} {
	if reflect.DeepEqual(material, gocd.Material{}) {
		return nil
	}

	result := map[string]interface{}{
		"type":        material.Type,
		"fingerprint": material.Fingerprint,
	}

	attrs := make(map[string]interface{})
	if !reflect.DeepEqual(material.Attributes, gocd.Attribute{}) {
		for _, field := range []struct {
			name string
			val  string
		}{
			{"url", material.Attributes.URL},
			{"username", material.Attributes.Username},
			{"encrypted_password", material.Attributes.EncryptedPassword},
			{"branch", material.Attributes.Branch},
			{"view", material.Attributes.View},
			{"port", material.Attributes.Port},
			{"project_path", material.Attributes.ProjectPath},
			{"domain", material.Attributes.Domain},
			{"ref", material.Attributes.Ref},
			{"name", material.Attributes.Name},
			{"stage", material.Attributes.Stage},
			{"pipeline", material.Attributes.Pipeline},
			{"destination", material.Attributes.Destination},
		} {
			if field.val != "" {
				attrs[field.name] = field.val
			}
		}

		for _, field := range []struct {
			name string
			val  bool
		}{
			{"auto_update", material.Attributes.AutoUpdate},
			{"check_externals", material.Attributes.CheckExternals},
			{"use_tickets", material.Attributes.UseTickets},
			{"ignore_for_scheduling", material.Attributes.IgnoreForScheduling},
			{"invert_filter", material.Attributes.InvertFilter},
		} {
			attrs[field.name] = field.val
		}
	}

	if len(attrs) > 0 {
		result["attributes"] = []interface{}{attrs}
	}

	return []interface{}{result}
}
