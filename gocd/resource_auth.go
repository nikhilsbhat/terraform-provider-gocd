// ----------------------------------------------------------------------------
//
//	***     TERRAGEN GENERATED CODE    ***    TERRAGEN GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//	This file was auto generated by Terragen.
//	This autogenerated code has to be enhanced further to make it fully working terraform-provider.
//
//	Get more information on how terragen works.
//	https://github.com/nikhilsbhat/terragen
//
// ----------------------------------------------------------------------------
package gocd

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAuth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthCreate,
		ReadContext:   resourceAuthRead,
		DeleteContext: resourceAuthDelete,
		UpdateContext: resourceAuthUpdate,
		Schema:        map[string]*schema.Schema{},
	}
}

func resourceAuthCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Your code goes here
	return nil
}

func resourceAuthRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Your code goes here
	return nil
}

func resourceAuthDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Your code goes here
	return nil
}

func resourceAuthUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Your code goes here
	return nil
}