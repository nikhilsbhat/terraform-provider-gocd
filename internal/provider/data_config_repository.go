package provider

import (
	"context"

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

func dataSourceConfigRepositoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()

	if len(id) == 0 {
		newID, err := utils.GetRandomID()
		if err != nil {
			d.SetId("")

			return diag.Errorf("errored while fetching randomID %v", err)
		}
		id = newID
	}

	profileID := utils.String(d.Get(utils.TerraformResourceProfileID))
	response, err := defaultConfig.GetConfigRepo(profileID)
	if err != nil {
		return diag.Errorf("getting config repo %s errored with: %v", profileID, err)
	}

	if err = d.Set(utils.TerraformPluginID, response.PluginID); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformPluginID, err)
	}

	if err = d.Set(utils.TerraformResourceMaterial, flattenMaterial(response.Material)); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceMaterial, err)
	}

	flattenedConfiguration, err := utils.MapSlice(response.Configuration)
	if err != nil {
		d.SetId("")

		return diag.Errorf("errored while flattening Configuration obtained: %v", err)
	}

	if err = d.Set(utils.TerraformResourceConfiguration, flattenedConfiguration); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourcePluginConfiguration, err)
	}

	if err = d.Set(utils.TerraformResourceRules, response.Rules); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceRules, err)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	d.SetId(id)

	return nil
}

func flattenMaterial(material gocd.Material) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type":        material.Type,
			"fingerprint": material.Fingerprint,
			"attributes":  flattenAttributes(material.Attributes),
		},
	}
}

func flattenAttributes(attribute gocd.Attribute) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"url":                attribute.URL,
			"username":           attribute.Username,
			"password":           attribute.Password,
			"encrypted_password": attribute.EncryptedPassword,
			"branch":             attribute.Branch,
			"auto_update":        attribute.AutoUpdate,
		},
	}
}
