package provider

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func dataSourceConfigRepository() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConfigRepositoryRead,
		Schema:      configRepoSchema(),
	}
}

func dataSourceConfigRepositoryRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()

	if len(id) == 0 {
		resourceID := utils.String(d.Get(utils.TerraformResourceProfileID))
		id = resourceID
	}

	profileID := utils.String(d.Get(utils.TerraformResourceProfileID))

	response, err := defaultConfig.GetConfigRepo(profileID)
	if err != nil {
		return diag.Errorf("getting config repo %s errored with: %v", profileID, err)
	}

	err = d.Set(utils.TerraformResourcePluginID, response.PluginID)
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourcePluginID, err)
	}

	flattened := flattenMaterialRead(response.Material)

	err = d.Set("material", flattened)
	if err != nil {
		return diag.Errorf("setting material errored with: %v", err)
	}

	flattenedConfiguration, err := utils.MapSlice(response.Configuration)
	if err != nil {
		return diag.Errorf("errored while flattening Configuration obtained: %v", err)
	}

	err = d.Set(utils.TerraformResourceConfiguration, flattenedConfiguration)
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceConfiguration, err)
	}

	err = d.Set(utils.TerraformResourceRules, response.Rules)
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceRules, err)
	}

	err = d.Set(utils.TerraformResourceEtag, response.ETAG)
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	d.SetId(id)

	return nil
}

func flattenMaterialRead(material gocd.Material) []any {
	if reflect.DeepEqual(material, gocd.Material{}) {
		return nil
	}

	materialMap := map[string]any{
		"type":        material.Type,
		"fingerprint": material.Fingerprint,
	}

	if !reflect.DeepEqual(material.Attributes, gocd.Attribute{}) {
		attrsMap := map[string]any{
			"url":                   material.Attributes.URL,
			"username":              material.Attributes.Username,
			"password":              material.Attributes.Password,
			"encrypted_password":    material.Attributes.EncryptedPassword,
			"branch":                material.Attributes.Branch,
			"auto_update":           material.Attributes.AutoUpdate,
			"check_externals":       material.Attributes.CheckExternals,
			"use_tickets":           material.Attributes.UseTickets,
			"view":                  material.Attributes.View,
			"port":                  material.Attributes.Port,
			"project_path":          material.Attributes.ProjectPath,
			"domain":                material.Attributes.Domain,
			"ref":                   material.Attributes.Ref,
			"name":                  material.Attributes.Name,
			"stage":                 material.Attributes.Stage,
			"pipeline":              material.Attributes.Pipeline,
			"ignore_for_scheduling": material.Attributes.IgnoreForScheduling,
			"destination":           material.Attributes.Destination,
			"invert_filter":         material.Attributes.InvertFilter,
		}

		for k, v := range attrsMap {
			if v == nil || v == "" || v == false {
				delete(attrsMap, k)
			}
		}

		materialMap["attributes"] = []any{attrsMap}
	}

	return []any{materialMap}
}
