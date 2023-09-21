package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func dataSourcePipelineGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourcePipelineGroupRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "Name of the pipeline group to be retrieved.",
			},
			"pipelines": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of pipelines those are part of this pipeline group.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"authorization": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "The authorization configuration for the pipeline group.",
				Elem:        authConfigSchema(),
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Etag used to track the pipeline group.",
			},
		},
	}
}

func datasourcePipelineGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()

	if len(id) == 0 {
		resourceID := utils.String(d.Get(utils.TerraformResourcePipelineGroupID))
		id = resourceID
	}

	groupID := utils.String(d.Get(utils.TerraformResourcePipelineGroupID))

	response, err := defaultConfig.GetPipelineGroup(groupID)
	if err != nil {
		return diag.Errorf("getting pipeline group %s errored with: %v", groupID, err)
	}

	if err = d.Set(utils.TerraformResourcePipelines, flattenPipelines(response.Pipelines)); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourcePipelines)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	flattenedAuthorization, err := flattenAuthorization(response.Authorization)
	if err != nil {
		d.SetId("")

		return diag.Errorf("errored while flattening Authorization obtained: %v", err)
	}

	if err = d.Set(utils.TerraformResourceAuthorization, flattenedAuthorization); err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceAuthorization)
	}

	d.SetId(id)

	return nil
}

func flattenAuthorization(authorization gocd.PipelineGroupAuthorizationConfig) ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{
			"admins":  getUsersNRoles(authorization.Admins),
			"operate": getUsersNRoles(authorization.Operate),
			"view":    getUsersNRoles(authorization.View),
		},
	}, nil
}

func getUsersNRoles(authorization gocd.AuthorizationConfig) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"users": authorization.Users,
			"roles": authorization.Roles,
		},
	}
}
