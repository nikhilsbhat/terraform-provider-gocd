package provider

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func resourceBackupSchedule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBackupScheduleCreate,
		ReadContext:   resourceBackupScheduleRead,
		DeleteContext: resourceBackupScheduleDelete,
		Schema: map[string]*schema.Schema{
			"schedule": {
				Type:        schema.TypeBool,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "Enable to trigger the backup, (would be unset post backup is successful).",
			},
			"retry": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    false,
				Default:     defaultRetry,
				ForceNew:    true,
				Description: "Number of times to retry to get the ID of latest successful backup taken.",
			},
			"delay": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    false,
				Default:     defaultDelay,
				ForceNew:    true,
				Description: "Time delay between each retries that would be made to get backup stats (in seconds ex: 5).",
			},
			"retry_after": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "This would be set to handle the backup scheduling internally.",
			},
			"backup_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Id of the backup that was taken successfully",
			},
		},
	}
}

var (
	defaultRetry = 30
	defaultDelay = 5
)

func resourceBackupScheduleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !d.IsNewResource() {
		return nil
	}

	var id string

	newID, err := utils.GetRandomID()
	if err != nil {
		d.SetId("")

		return diag.Errorf("errored while fetching randomID %v", err)
	}
	id = newID

	if !utils.Bool(d.Get(utils.TerraformResourceSchedule)) {
		return diag.Errorf("scheduling backup is disabled, set attribute '%s' to true to schedule a backup", utils.TerraformResourceSchedule)
	}

	response, err := defaultConfig.ScheduleBackup()
	if err != nil {
		return diag.Errorf("scheduling backup errored with: %v", err)
	}

	retryAfter, err := strconv.Atoi(response["RetryAfter"])
	if err != nil {
		return diag.Errorf("%v", err)
	}

	if err = d.Set(utils.TerraformResourceBackupID, response["BackUpID"]); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceBackupID, err)
	}

	if err = d.Set(utils.TerraformResourceRetryAfter, retryAfter); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceRetryAfter, err)
	}

	d.SetId(id)

	return resourceBackupScheduleRead(ctx, d, meta)
}

func resourceBackupScheduleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	retryAfter := d.Get(utils.TerraformResourceRetryAfter).(int)
	backupRetry := d.Get(utils.TerraformResourceRetry).(int)
	delayCount := d.Get(utils.TerraformResourceDelay).(int)
	backupID := utils.String(d.Get(utils.TerraformResourceBackupID))

	delay := time.Duration(delayCount) * time.Second

	time.Sleep(time.Duration(retryAfter) * time.Second)
	currentRetryCount := 0
	var latestBackupStatus string
	for {
		if currentRetryCount > backupRetry {
			return diag.Errorf("maximum retry count of '%d' crossed with current count '%d', still backup is not ready yet with status '%s'. Exiting",
				backupRetry, currentRetryCount, latestBackupStatus)
		}

		response, err := defaultConfig.GetBackup(backupID)
		if err != nil {
			return diag.Errorf("getting last configured backup ID errored with: %v", err)
		}

		retryRemaining := backupRetry - currentRetryCount
		if response.Status == "IN_PROGRESS" {
			log.Printf("the backup stats is still in IN_PROGRESS status, retrying... '%d' more to go", retryRemaining)
		}

		if response.Status == "COMPLETED" {
			if err = d.Set(utils.TerraformResourceSchedule, false); err != nil {
				return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceSchedule, err)
			}

			break
		}

		if response.Status != "COMPLETED" && response.Status != "IN_PROGRESS" {
			return diag.Errorf("looks like backup status is neither IN_PROGRESS nor COMPLETED rather it is %s", response.Status)
		}

		latestBackupStatus = response.Status
		time.Sleep(delay)
		currentRetryCount++
	}

	return nil
}

func resourceBackupScheduleDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	id := d.Id()
	if len(d.Id()) == 0 {
		return diag.Errorf("resource with the ID '%s' not found", id)
	}

	d.SetId("")

	return nil
}
