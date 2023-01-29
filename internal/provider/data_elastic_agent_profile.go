package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func dataSourceElasticAgentProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceElasticAgentProfileRead,
		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The identifier of the elastic agent profile.",
			},
			"cluster_profile_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Required:    false,
				Description: "The identifier of the cluster profile to which current elastic agent profile belongs.",
			},
			"properties": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "The list of configuration properties that represent the configuration of this profile.",
				Elem:        propertiesSchemaData(),
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Etag used to track the elastic agent profile.",
			},
		},
	}
}

func datasourceElasticAgentProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	response, err := defaultConfig.GetElasticAgentProfile(profileID)
	if err != nil {
		return diag.Errorf("getting elastic agent profile %s errored with: %v", profileID, err)
	}

	if err = d.Set(utils.TerraformResourceClusterProfileID, response.ClusterProfileID); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceClusterProfileID)
	}

	flattenedProperties, err := utils.MapSlice(response.Properties)
	if err != nil {
		d.SetId("")

		return diag.Errorf("errored while flattening Properties obtained: %v", err)
	}

	if err = d.Set(utils.TerraformProperties, flattenedProperties); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformProperties)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf("setting etag errored with %v", err)
	}

	d.SetId(id)

	return nil
}
