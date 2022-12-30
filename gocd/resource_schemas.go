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
				"is_secure": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: false,
					ForceNew: true,
					Description: "Specify whether the given property is secure or not. If true and encrypted_value is not specified, " +
						"GoCD will store the value in encrypted format.",
				},
			},
		},
	}
}

func materialSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Required:    true,
		Computed:    false,
		ForceNew:    true,
		Description: "The material to be used by the config repo.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:        schema.TypeString,
					Required:    true,
					Computed:    false,
					ForceNew:    false,
					Description: "The type of a material. Can be one of git, svn, hg, p4, tfs.",
				},
				"fingerprint": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    false,
					ForceNew:    false,
					Description: "The fingerprint of the material.",
				},
				"attributes": {
					Type:        schema.TypeSet,
					Required:    true,
					Computed:    false,
					ForceNew:    false,
					Description: "The attributes for each material type.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"url": {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "The URL of the subversion repository.",
							},
							"username": {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "The user account for the remote repository.",
							},
							"password": {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "The password for the specified user.",
							},
							"encrypted_password": {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "The encrypted password for the specified user.",
							},
							"branch": {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "The mercurial branch to build.",
							},
							"view": {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "The Perforce view.",
							},
							"port": {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "Perforce server connection to use ([transport:]host:port).",
							},
							"project_path": {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "The project path within the TFS collection.",
							},
							"domain": {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "\tThe domain name for TFS authentication credentials.",
							},
							"ref": {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "The unique package repository id.",
							},
							"name": {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "The name of this material.",
							},
							"stage": {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "The name of a stage which will trigger this pipeline once it is successful.",
							},
							"pipeline": {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "The name of a pipeline that this pipeline depends on.",
							},
							"destination": {
								Type:        schema.TypeString,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "The directory (relative to the pipeline directory) in which source code will be checked out.",
							},
							"auto_update": {
								Type:        schema.TypeBool,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "Whether to poll for new changes or not.",
							},
							"check_externals": {
								Type:        schema.TypeBool,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "Whether the changes o the externals will trigger the pipeline automatically or not.",
							},
							"use_tickets": {
								Type:        schema.TypeBool,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "Whether to work with the Perforce tickets or not.",
							},
							"ignore_for_scheduling": {
								Type:        schema.TypeBool,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "Whether the pipeline should be triggered when there are changes in this material.",
							},
							"invert_filter": {
								Type:        schema.TypeBool,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "Invert filter to enable whitelist.",
							},
							"filter": {
								Type:        schema.TypeSet,
								Optional:    true,
								Computed:    false,
								ForceNew:    false,
								Description: "The filter specifies files in changesets that should not trigger a pipeline automatically.",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"ignore": {
											Type:        schema.TypeList,
											Optional:    true,
											Computed:    false,
											ForceNew:    false,
											Description: "Invert filter to enable whitelist.",
											Elem:        &schema.Schema{Type: schema.TypeString},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
