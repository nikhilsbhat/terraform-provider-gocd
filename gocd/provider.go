package gocd

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/client"
)

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
				Description: "base url of GoCD server, which this terraform provider can interact with",
			},
			"ca_file": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Computed:    false,
				DefaultFunc: schema.EnvDefaultFunc("GOCD_CAFILE_CONTENT", "some_ca_context"),
				Description: "CA file contents, to be used while connecting to GoCD server when CA based auth is enabled",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Computed:    false,
				DefaultFunc: schema.EnvDefaultFunc("GOCD_USERNAME", "username"),
				Description: "username to be used while connecting with GoCD",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Computed:    false,
				DefaultFunc: schema.EnvDefaultFunc("GOCD_PASSWORD", "password"),
				Description: "password to be used while connecting with GoCD",
			},
			"loglevel": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Computed:    false,
				DefaultFunc: schema.EnvDefaultFunc("GOCD_PASSWORD", "password"),
				Description: "loglevel to be set for the api calls made to GoCD",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"gocd_plugins": resourcePlugins(),
			// "gocd_auth":    resourceAuth(),
			"gocd_cluster_profile":       clusterProfile(),
			"gocd_elastic_agent_profile": elasticAgentProfile(),
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ConfigureContextFunc: client.GetGoCDClient,
	}
}
