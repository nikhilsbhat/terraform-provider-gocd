package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func resourceBackupConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBackupConfigCreate,
		ReadContext:   resourceBackupConfigRead,
		UpdateContext: resourceBackupConfigUpdate,
		DeleteContext: resourceBackupConfigDelete,
		Schema: map[string]*schema.Schema{
			"schedule": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				Description: "The backup schedule. See the quartz documentation for syntax and examples.",
			},
			"post_backup_script": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				Description: "The script that will be executed once the backup finishes. See the gocd documentation for details.",
			},
			"email_on_success": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    false,
				Description: "If set to true, an email will be sent when backup completes successfully.",
			},
			"email_on_failure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    false,
				Description: "If set to true, an email will be sent when backup fails.",
			},
		},
	}
}

func resourceBackupConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !d.IsNewResource() {
		return nil
	}

	id := d.Id()

	if len(id) == 0 {
		resourceID := utils.String(d.Get(utils.TerraformResourceSchedule))
		id = resourceID
	}

	cfg := gocd.BackupConfig{
		Schedule:         utils.String(d.Get(utils.TerraformResourceSchedule)),
		PostBackupScript: utils.String(d.Get(utils.TerraformResourcePostBackupScript)),
		EmailOnSuccess:   utils.Bool(d.Get(utils.TerraformResourceEmailOnSuccess)),
		EmailOnFailure:   utils.Bool(d.Get(utils.TerraformResourceEmailOnFailure)),
	}

	if err := defaultConfig.CreateOrUpdateBackupConfig(cfg); err != nil {
		return diag.Errorf("creating backup configuration errored with %v", err)
	}

	d.SetId(id)

	return resourceBackupConfigRead(ctx, d, meta)
}

func resourceBackupConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	response, err := defaultConfig.GetBackupConfig()
	if err != nil {
		return diag.Errorf("getting backup configuration errored with: %v", err)
	}

	if err = d.Set(utils.TerraformResourceSchedule, response.Schedule); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceSchedule, err)
	}

	if err = d.Set(utils.TerraformResourcePostBackupScript, response.PostBackupScript); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourcePostBackupScript, err)
	}

	if err = d.Set(utils.TerraformResourceEmailOnSuccess, response.EmailOnSuccess); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEmailOnSuccess, err)
	}

	if err = d.Set(utils.TerraformResourceEmailOnFailure, response.EmailOnFailure); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEmailOnFailure, err)
	}

	return nil
}

func resourceBackupConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !d.HasChanges(
		utils.TerraformResourceSchedule,
		utils.TerraformResourcePostBackupScript,
		utils.TerraformResourceEmailOnSuccess,
		utils.TerraformResourceEmailOnFailure,
	) {
		log.Printf("nothing to update so skipping")

		return nil
	}

	cfg := gocd.BackupConfig{
		Schedule:         utils.String(d.Get(utils.TerraformResourceSchedule)),
		PostBackupScript: utils.String(d.Get(utils.TerraformResourcePostBackupScript)),
		EmailOnSuccess:   utils.Bool(d.Get(utils.TerraformResourceEmailOnSuccess)),
		EmailOnFailure:   utils.Bool(d.Get(utils.TerraformResourceEmailOnFailure)),
	}

	if err := defaultConfig.CreateOrUpdateBackupConfig(cfg); err != nil {
		return diag.Errorf("updating backup configuration errored with %v", err)
	}

	return resourceBackupConfigRead(ctx, d, meta)
}

func resourceBackupConfigDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()
	if len(d.Id()) == 0 {
		return diag.Errorf("resource with the ID '%s' not found", id)
	}

	err := defaultConfig.DeleteBackupConfig()
	if err != nil {
		return diag.Errorf("deleting backup configuration errored with: %v", err)
	}

	d.SetId("")

	return nil
}
