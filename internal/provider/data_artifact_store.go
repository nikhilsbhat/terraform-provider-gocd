package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func dataSourceArtifactStore() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceArtifactStoreRead,
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
				Computed:    true,
				Optional:    true,
				Description: "The plugin identifier of the artifact plugin.",
			},
			"properties": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "The list of configuration properties that represent the configuration of this artifact store.",
				Elem:        propertiesSchemaData(),
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Etag used to track an artifact store",
			},
		},
	}
}

func datasourceArtifactStoreRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()

	if len(id) == 0 {
		resourceID := utils.String(d.Get(utils.TerraformResourceStoreID))
		id = resourceID
	}

	response, err := defaultConfig.GetArtifactStore(id)
	if err != nil {
		return diag.Errorf("getting artifact store '%s' errored with: %v", id, err)
	}

	if err = d.Set(utils.TerraformResourcePluginID, response.PluginID); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourcePluginID)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	flattenedProperties, err := utils.MapSlice(response.Properties)
	if err != nil {
		d.SetId("")

		return diag.Errorf("errored while flattening artifact store properties obtained: %v", err)
	}

	if err = d.Set(utils.TerraformResourceProperties, flattenedProperties); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceProperties)
	}

	d.SetId(id)

	return nil
}
