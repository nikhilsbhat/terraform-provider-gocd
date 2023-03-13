package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func resourceAgentConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgentConfigCreate,
		ReadContext:   resourceAgentConfigRead,
		DeleteContext: resourceAgentConfigDelete,
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
				Computed:    false,
				ForceNew:    true,
				Description: "The hostname of the agent.",
			},
			"agent_config_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "Whether an agent is enabled or not. Can be one of `Enabled`, `Disabled`.",
			},
			"resources": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The set of resources that this agent is tagged with (if agent is not an elastic agent).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"environments": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The set of environments that this agent belongs to.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				Description: "The IP address of the agent.",
			},
			"operating_system": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
				Description: "The operating system as reported by the agent.",
			},
		},
	}
}

func resourceAgentConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !d.IsNewResource() {
		return nil
	}

	id := d.Id()

	if len(id) == 0 {
		resourceID := utils.String(d.Get(utils.TerraformResourceUUID))
		id = resourceID
	}

	cfg := gocd.Agent{
		ID:           id,
		Name:         utils.String(d.Get(utils.TerraformResourceHostname)),
		Environments: utils.GetSlice(d.Get(utils.TerraformResourceEnvironments).([]interface{})),
		Resources:    utils.GetSlice(d.Get(utils.TerraformResourceResources).([]interface{})),
		ConfigState:  utils.String(d.Get(utils.TerraformResourceAgentConfigState)),
	}

	if err := defaultConfig.UpdateAgent(cfg); err != nil {
		return diag.Errorf("updating agent '%s' errored with %v", id, err)
	}

	d.SetId(id)

	return resourceAgentConfigRead(ctx, d, meta)
}

func resourceAgentConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	response, err := defaultConfig.GetAgent(d.Id())
	if err != nil {
		return diag.Errorf("fetching information of agent '%s' errored with %v", d.Id(), err)
	}

	if err = d.Set(utils.TerraformResourceIPAddress, response.IPAddress); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceIPAddress, err)
	}

	if err = d.Set(utils.TerraformResourceOperatingSystem, response.OS); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceOperatingSystem, err)
	}

	return nil
}

func resourceAgentConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Id()
	if len(d.Id()) == 0 {
		return diag.Errorf("resource with the ID '%s' not found", id)
	}

	d.SetId("")

	return nil
}
