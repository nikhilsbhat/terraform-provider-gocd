package gocd

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func propertiesSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Required:    true,
		Computed:    false,
		Description: "the list of configuration properties that represent the configuration of this profile.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:        schema.TypeString,
					Required:    true,
					Computed:    false,
					Description: "the name of the property key.",
				},
				"value": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    false,
					ForceNew:    true,
					Description: "The value of the property",
				},
				"encrypted_value": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    false,
					ForceNew:    true,
					Description: "The encrypted value of the property",
				},
			},
		},
	}
}
