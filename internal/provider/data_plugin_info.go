package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func dataSourcePluginInfo() *schema.Resource {
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
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Status of the plugin. Can be one of active, invalid.",
			},
			"about": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Additional details about the plugin.",
			},
			"extensions": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: pluginAttributesSchema(),
				},
				Description: "A list of extension information pertaining to the list of extensions the plugin implements.",
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

func pluginAttributesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Default:     nil,
			Description: "The type of the plugin extension.",
		},
		"display_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Default:     nil,
			Description: "The descriptive name of the plugin.",
		},
		"auth_config_settings": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: pluginSettingAttributeSchema(),
			},
			Default:     nil,
			Description: "The list of properties that can be used to configure auth configs.",
		},
		"artifact_config_settings": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: pluginSettingAttributeSchema(),
			},
			Default:     nil,
			Description: "The publish artifact configs and the list of properties that can be used to configure publish artifact config.",
		},
		"elastic_agent_profile_settings": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: pluginSettingAttributeSchema(),
			},
			Default:     nil,
			Description: "The elastic agent profile and the list of properties required to be configured.",
		},
		"fetch_artifact_settings": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: pluginSettingAttributeSchema(),
			},
			Default:     nil,
			Description: "The fetch artifact config and the list of properties that can be used to configure fetch artifact config.",
		},
		"cluster_profile_settings": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: pluginSettingAttributeSchema(),
			},
			Default:     nil,
			Description: "The cluster profile and the list of properties required to be configured. Present in case of plugin supports defining cluster profiles.",
		},
		"plugin_settings": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: pluginSettingAttributeSchema(),
			},
			Default:     nil,
			Description: "The plugin and the list of properties required to be configured.",
		},
		"package_settings": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: pluginSettingAttributeSchema(),
			},
			Default:     nil,
			Description: "The list of properties that can be used to configure a package material",
		},
		"repository_settings": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: pluginSettingAttributeSchema(),
			},
			Default:     nil,
			Description: "The list of properties that can be used to configure package repositories.",
		},
		"scm_settings": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: pluginSettingAttributeSchema(),
			},
			Default:     nil,
			Description: "List of properties that can be used to configure the pluggable scm material.",
		},
		"store_config_settings": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: pluginSettingAttributeSchema(),
			},
			Default:     nil,
			Description: "The list of properties that can be used to configure the artifact store.",
		},
		"secret_config_settings": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: pluginSettingAttributeSchema(),
			},
			Default:     nil,
			Description: "List of properties that can be used to configure the secret configs.",
		},
		"role_settings": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: pluginSettingAttributeSchema(),
			},
			Default:     nil,
			Description: "The list of properties that can be used to configure role configs.",
		},
		"task_settings": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: pluginSettingAttributeSchema(),
			},
			Default:     nil,
			Description: "The list of properties that can be used to configure the pluggable task.",
		},
	}
}

func pluginSettingAttributeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"configurations": {
			Type:     schema.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: pluginConfigurationSchema(),
			},
			Default:     nil,
			Description: "List of configuration required to configure the plugin settings.",
		},
	}
}

func pluginConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Default:     nil,
			Description: "The name of the property key.",
		},
		"value": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Default:     nil,
			Description: "The value of the property",
		},
		"encrypted_value": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Default:     nil,
			Description: "The encrypted value of the property.",
		},
		"is_secure": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Default:     nil,
			Description: "Specify whether the given property is secure or not.",
		},
		"metadata": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
			Default:  nil,
			Elem: &schema.Schema{
				Type: schema.TypeBool,
			},
			Description: "Metadata for the configuration property.",
		},
	}
}

func datasourcePluginInfoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()

	if len(id) == 0 {
		pluginID := utils.String(d.Get(utils.TerraformResourcePluginID))
		id = pluginID
	}

	pluginID := utils.String(d.Get(utils.TerraformResourcePluginID))
	response, err := defaultConfig.GetPluginInfo(pluginID)
	if err != nil {
		return diag.Errorf("getting plugin information of '%s' errored with: %v", pluginID, err)
	}

	if err = d.Set(utils.TerraformResourcePluginLocation, response.PluginFileLocation); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourcePluginLocation)
	}

	if err = d.Set(utils.TerraformResourcePluginBundled, response.BundledPlugin); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourcePluginBundled)
	}

	if err = d.Set(utils.TerraformResourcePluginStatus, response.Status.State); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourcePluginStatus)
	}

	if err = d.Set(utils.TerraformResourceExtensions, flattenPluginExtensions(response.Extensions)); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceExtensions)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	d.SetId(id)

	return nil
}

func flattenPluginExtensions(extensions []gocd.PluginAttributes) []map[string]interface{} {
	pluginExtensions := make([]map[string]interface{}, 0)

	for _, extension := range extensions {
		pluginExtension := make(map[string]interface{})

		if len(extension.Type) != 0 {
			pluginExtension["type"] = extension.Type
		}

		if len(extension.DisplayName) != 0 {
			pluginExtension["display_name"] = extension.DisplayName
		}

		if extension.AuthConfigSettings != nil {
			pluginExtension["auth_config_settings"] = flattenPluginSettingAttributeSchema(extension.AuthConfigSettings)
		}

		if extension.ArtifactConfigSettings != nil {
			pluginExtension["artifact_config_settings"] = flattenPluginSettingAttributeSchema(extension.ArtifactConfigSettings)
		}

		if extension.ElasticAgentProfileSettings != nil {
			pluginExtension["elastic_agent_profile_settings"] = flattenPluginSettingAttributeSchema(extension.ElasticAgentProfileSettings)
		}

		if extension.FetchArtifactSettings != nil {
			pluginExtension["fetch_artifact_settings"] = flattenPluginSettingAttributeSchema(extension.FetchArtifactSettings)
		}

		if extension.ClusterProfileSettings != nil {
			pluginExtension["cluster_profile_settings"] = flattenPluginSettingAttributeSchema(extension.ClusterProfileSettings)
		}

		if extension.PluginSettings != nil {
			pluginExtension["plugin_settings"] = flattenPluginSettingAttributeSchema(extension.PluginSettings)
		}

		if extension.PackageSettings != nil {
			pluginExtension["package_settings"] = flattenPluginSettingAttributeSchema(extension.PackageSettings)
		}

		if extension.RepositorySettings != nil {
			pluginExtension["repository_settings"] = flattenPluginSettingAttributeSchema(extension.RepositorySettings)
		}

		if extension.ScmSettings != nil {
			pluginExtension["scm_settings"] = flattenPluginSettingAttributeSchema(extension.ScmSettings)
		}

		if extension.StoreConfigSettings != nil {
			pluginExtension["store_config_settings"] = flattenPluginSettingAttributeSchema(extension.StoreConfigSettings)
		}

		if extension.SecretConfigSettings != nil {
			pluginExtension["secret_config_settings"] = flattenPluginSettingAttributeSchema(extension.SecretConfigSettings)
		}

		if extension.RoleSettings != nil {
			pluginExtension["role_settings"] = flattenPluginSettingAttributeSchema(extension.RoleSettings)
		}

		if extension.TaskSettings != nil {
			pluginExtension["task_settings"] = flattenPluginSettingAttributeSchema(extension.TaskSettings)
		}

		pluginExtensions = append(pluginExtensions, pluginExtension)
	}

	return pluginExtensions
}

func flattenPluginSettingAttributeSchema(extensions *gocd.PluginSettingAttribute) []map[string]interface{} {
	if extensions == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"configurations": flattenPluginConfigurations(extensions.Configurations),
		},
	}
}

func flattenPluginConfigurations(configurations []gocd.PluginConfiguration) []map[string]interface{} {
	pluginConfigurations := make([]map[string]interface{}, 0)

	for _, configuration := range configurations {
		pluginConfiguration := map[string]interface{}{
			"key":             configuration.Key,
			"value":           configuration.Value,
			"encrypted_value": configuration.EncryptedValue,
			"is_secure":       configuration.IsSecure,
		}

		if configuration.Metadata != nil {
			pluginConfiguration["metadata"] = flattenMetadata(configuration.Metadata)
		}

		pluginConfigurations = append(pluginConfigurations, pluginConfiguration)
	}

	return pluginConfigurations
}

func flattenMetadata(meta map[string]interface{}) map[string]interface{} {
	metadata := make(map[string]interface{})

	if _, ok := meta["part_of_identity"]; ok {
		metadata["part_of_identity"] = meta["part_of_identity"]
	}

	if _, ok := metadata["required"]; ok {
		metadata["required"] = meta["required"]
	}

	if _, ok := metadata["secure"]; ok {
		metadata["secure"] = meta["secure"]
	}

	return metadata
}
