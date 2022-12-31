package gocd

import (
	"context"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func encryptValue() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEncryptValueCreate,
		ReadContext:   resourceEncryptValueRead,
		DeleteContext: resourceEncryptValueDelete,
		Schema: map[string]*schema.Schema{
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "Plain text value to encrypt.",
			},
			"encrypted_value": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Required:    false,
				Computed:    true,
				ForceNew:    false,
				Description: "Encrypted value of plain text.",
			},
		},
	}
}

func resourceEncryptValueCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !d.IsNewResource() {
		return nil
	}

	id := d.Id()

	if len(id) == 0 {
		newID, err := utils.GetRandomID()
		if err != nil {
			d.SetId("")

			return diag.Errorf("errored while fetching randomID %v", err)
		}
		id = newID
	}

	encryptedValue, err := defaultConfig.EncryptText(utils.String(d.Get(utils.TerraformValue)))
	if err != nil {
		return diag.Errorf("encrypting value errored with %v", err)
	}

	if err = d.Set(utils.TerraformEncryptedValue, encryptedValue.EncryptedValue); err != nil {
		return diag.Errorf("setting etag errored with %v", err)
	}

	d.SetId(id)

	return nil
}

func resourceEncryptValueRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceEncryptValueDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Id()
	if len(d.Id()) == 0 {
		return diag.Errorf("resource with the ID '%s' not found", id)
	}

	d.SetId("")

	return nil
}
