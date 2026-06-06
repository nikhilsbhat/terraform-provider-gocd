package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
	"github.com/spf13/cast"
)

func dataSourceAgentConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceAgentRead,
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The identifier of this agent.",
			},
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				Description: "The hostname of the agent.",
			},
			"elastic_agent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				Description: "The elastic agent identifier of this agent. This attribute is only available if the agent is an elastic agent.",
			},
			"elastic_plugin_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				Description: "The identifier of the elastic agent plugin that manages this agent instance. This attribute is only available if the agent is an elastic agent.",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				Description: "The IP address of the agent.",
			},
			"sandbox": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				Description: "The path where the agent will perform its builds.",
			},
			"operating_system": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				Description: "The operating system as reported by the agent.",
			},
			"free_space": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				Description: "The amount of free space in bytes.",
			},
			"agent_config_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				Description: "Whether an agent is enabled or not. Can be one of `Pending`, `Enabled`, `Disabled`.",
			},
			"agent_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				Description: "The state an agent is in. Can be one of Idle, `Building`, `LostContact`, `Missing`, `Building`, `Unknown`.",
			},
			"agent_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				Description: "The version of the agent.",
			},
			"resources": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				Description: "The set of resources that this agent is tagged with (if agent is not an elastic agent).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"environments": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				Description: "The set of environments that this agent belongs to.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"build_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				Description: "If the agent is running a build, the state of the build on the agent. Can be one of Idle, `Building`, `Cancelled`, `Unknown`.",
			},
			"build_details": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				Description: "The build details provides information like pipeline, stage and job if the build_state of the agent is `Building`",
			},
		},
	}
}

func datasourceAgentRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()

	if len(id) == 0 {
		resourceID := utils.String(d.Get(utils.TerraformResourceUUID))
		id = resourceID
	}

	agentID := utils.String(d.Get(utils.TerraformResourceUUID))

	response, err := defaultConfig.GetAgent(agentID)
	if err != nil {
		return diag.Errorf("getting information of agent '%s' errored with: %v", agentID, err)
	}

	err = d.Set(utils.TerraformResourceHostname, response.Name)
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceHostname)
	}

	err = d.Set(utils.TerraformResourceElasticAgentAD, response.ElasticAgentID)
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceElasticAgentAD)
	}

	err = d.Set(utils.TerraformResourceElasticPluginAD, response.ElasticPluginID)
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceElasticPluginAD)
	}

	err = d.Set(utils.TerraformResourceIPAddress, response.IPAddress)
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceIPAddress)
	}

	err = d.Set(utils.TerraformResourceSandbox, response.Sandbox)
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceSandbox)
	}

	err = d.Set(utils.TerraformResourceOperatingSystem, response.OS)
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceOperatingSystem)
	}

	err = d.Set(utils.TerraformResourceFreeSpace, cast.ToFloat64(response.DiskSpaceAvailable))
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceFreeSpace)
	}

	err = d.Set(utils.TerraformResourceAgentConfigState, response.ConfigState)
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceAgentConfigState)
	}

	err = d.Set(utils.TerraformResourceAgentState, response.CurrentState)
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceAgentState)
	}

	err = d.Set(utils.TerraformResourceAgentVersion, response.Version)
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceAgentVersion)
	}

	err = d.Set(utils.TerraformResourceResources, response.Resources)
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceResources)
	}

	err = d.Set(utils.TerraformResourceEnvironments, flattenEnvironments(response.Environments))
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceEnvironments)
	}

	err = d.Set(utils.TerraformResourceBuildState, response.BuildState)
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceBuildState)
	}

	buildDetails, err := utils.Map(response.BuildDetails)
	if err != nil {
		return diag.Errorf("flattening '%s' errored with :%v", utils.TerraformResourceBuildDetails, err)
	}

	err = d.Set(utils.TerraformResourceBuildDetails, buildDetails)
	if err != nil {
		return diag.Errorf(settingAttrErrorTmp, err, utils.TerraformResourceBuildDetails)
	}

	d.SetId(id)

	return nil
}

func flattenEnvironments(envs any) []string {
	environments := envs.([]any)
	envList := make([]string, 0, len(environments))

	for _, environment := range environments {
		newEnvironment := environment.(map[string]any)
		envList = append(envList, newEnvironment["name"].(string))
	}

	return envList
}
