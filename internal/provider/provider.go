package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/client"
)

func init() { //nolint:gochecknoinits
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Computed:    false,
				DefaultFunc: schema.EnvDefaultFunc("GOCD_BASE_URL", "www.gocd.com"),
				Description: "base url of GoCD server, with which this terraform provider will with (https://gocd.myself.com/go)",
			},
			"ca_file": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    false,
				DefaultFunc: schema.EnvDefaultFunc("GOCD_CAFILE_CONTENT", nil),
				Description: "CA file contents, to be used while connecting to GoCD server when CA based auth is enabled",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Computed:    false,
				DefaultFunc: schema.EnvDefaultFunc("GOCD_USERNAME", nil),
				Description: "username to be used while connecting with GoCD",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    false,
				DefaultFunc: schema.EnvDefaultFunc("GOCD_PASSWORD", nil),
				Description: "password to be used while connecting with GoCD",
			},
			"auth_token": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      false,
				DefaultFunc:   schema.EnvDefaultFunc("GOCD_AUTH_TOKEN", nil),
				ConflictsWith: []string{"password"},
				Description: "bearer-token to be used while connecting with GoCD (API: https://api.gocd.org/current/#access-tokens, " +
					"UI: https://docs.gocd.org/current/configuration/access_tokens.html) cannot co-exist with password based auth.",
			},
			"loglevel": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Computed:    false,
				DefaultFunc: schema.EnvDefaultFunc("GOCD_LOGLEVEL", "info"),
				Description: "loglevel to be set for the api calls made to GoCD",
			},
			"skip_check": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Computed:    false,
				DefaultFunc: schema.EnvDefaultFunc("GOCD_SKIP_CHECK", "false"),
				Description: "setting this to false will skip a validation done during client creation, this helps by avoiding " +
					"errors being thrown from all resource/data block defined",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"gocd_plugin_setting":        resourcePluginsSetting(),
			"gocd_auth_config":           resourceAuthConfig(),
			"gocd_cluster_profile":       resourceClusterProfile(),
			"gocd_elastic_agent_profile": resourceElasticAgentProfile(),
			"gocd_config_repository":     resourceConfigRepository(),
			"gocd_environment":           resourceEnvironment(),
			"gocd_encrypt_value":         resourceEncryptValue(),
			"gocd_secret_config":         resourceSecretConfig(),
			"gocd_backup_config":         resourceBackupConfig(),
			"gocd_backup_schedule":       resourceBackupSchedule(),
			"gocd_agent":                 resourceAgentConfig(),
			"gocd_pipeline":              resourcePipeline(),
			"gocd_artifact_store":        resourceArtifactStore(),
			"gocd_role":                  resourceRole(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"gocd_plugin_setting":        dataSourcePluginsSetting(),
			"gocd_auth_config":           dataSourceAuthConfig(),
			"gocd_cluster_profile":       dataSourceClusterProfile(),
			"gocd_elastic_agent_profile": dataSourceElasticAgentProfile(),
			"gocd_config_repository":     dataSourceConfigRepository(),
			"gocd_environment":           dataSourceEnvironment(),
			"gocd_secret_config":         dataSourceSecretConfig(),
			"gocd_plugin_info":           dataSourcePluginInfo(),
			"gocd_agent":                 dataSourceAgentConfig(),
			"gocd_pipeline":              dataSourcePipeline(),
			"gocd_artifact_store":        dataSourceArtifactStore(),
			"gocd_role":                  dataSourceRole(),
		},

		ConfigureContextFunc: client.GetGoCDClient,
	}
}
