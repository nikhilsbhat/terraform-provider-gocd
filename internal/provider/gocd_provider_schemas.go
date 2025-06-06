package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var settingAttrErrorTmp = "setting '%s' errored with '%v'"

func configRepoSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"profile_id": {
			Type:        schema.TypeString,
			Required:    true,
			Computed:    false,
			ForceNew:    true,
			Description: "The identifier of the config repository.",
		},
		"plugin_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Required:    false,
			Description: "The name of the config repo plugin.",
		},
		"material": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "The material to be used by the config repo.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "The type of a material. Can be one of git, svn, hg, p4, tfs.",
					},
					"fingerprint": {
						Type:        schema.TypeString,
						Optional:    true,
						Computed:    true,
						Description: "The fingerprint of the material.",
					},
					"attributes": {
						Type:        schema.TypeList,
						Computed:    true,
						Description: "The attributes for each material type.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"url": {
									Type:        schema.TypeString,
									Optional:    true,
									Computed:    true,
									Description: "The URL of the subversion repository.",
								},
								"username": {
									Type:        schema.TypeString,
									Optional:    true,
									Computed:    true,
									Description: "The user account for the remote repository.",
								},
								"password": {
									Type:        schema.TypeString,
									Optional:    true,
									Computed:    true,
									Description: "The password for the specified user.",
								},
								"encrypted_password": {
									Type:        schema.TypeString,
									Optional:    true,
									Computed:    true,
									Description: "The encrypted password for the specified user.",
								},
								"branch": {
									Type:        schema.TypeString,
									Optional:    true,
									Computed:    true,
									Description: "The mercurial branch to build.",
								},
								"auto_update": {
									Type:        schema.TypeBool,
									Optional:    true,
									Computed:    true,
									Description: "Whether to poll for new changes or not.",
								},
							},
						},
					},
				},
			},
		},
		"configuration": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "The list of properties (key-value pairs) to be provided for using the plugin.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:        schema.TypeString,
						Optional:    true,
						Computed:    true,
						Description: "The name of the property key.",
					},
					"value": {
						Type:        schema.TypeString,
						Optional:    true,
						Computed:    true,
						Description: "The value of the property",
					},
					"encrypted_value": {
						Type:        schema.TypeString,
						Optional:    true,
						Computed:    true,
						Description: "The encrypted value of the property",
					},
					"is_secure": {
						Type:     schema.TypeBool,
						Optional: true,
						Computed: true,
						Description: "Specify whether the given property is secure or not. If true and encrypted_value is not specified, " +
							"GoCD will store the value in encrypted format.",
					},
				},
			},
		},
		"rules": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "The list of rules, which allows restricting the entities that the config repo can refer to.",
			Elem:        rulesSchema(),
		},
		"etag": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Etag used to track the config repository.",
		},
	}
}

func environmentsSchemaResource() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Computed:    false,
		Optional:    true,
		Description: "The list of environment variables that will be passed to all tasks (commands) that are part of this environment.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    false,
					Description: "The name of the environment variable.",
				},
				"value": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    false,
					Description: "The value of the environment variable. You MUST specify one of value or encrypted_value.",
				},
				"encrypted_value": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    false,
					Description: "The encrypted value of the environment variable. You MUST specify one of value or encrypted_value.",
				},
				"secure": {
					Type:        schema.TypeBool,
					Optional:    true,
					Computed:    false,
					Description: "Whether environment variable is secure or not. When set to true, encrypts the value if one is specified. The default value is false.",
				},
			},
		},
	}
}

func environmentsSchemaData() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		Description: "The list of environment variables that will be passed to all tasks (commands) that are part of this environment.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    false,
					Description: "The name of the environment variable.",
				},
				"value": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    false,
					Description: "The value of the environment variable. You MUST specify one of value or encrypted_value.",
				},
				"encrypted_value": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    false,
					Description: "The encrypted value of the environment variable. You MUST specify one of value or encrypted_value.",
				},
				"secure": {
					Type:        schema.TypeBool,
					Optional:    true,
					Computed:    false,
					Description: "Whether environment variable is secure or not. When set to true, encrypts the value if one is specified. The default value is false.",
				},
			},
		},
	}
}

func propertiesSchemaResource() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Required:    true,
		Computed:    false,
		Description: "The list of configuration properties that represent the configuration of the profile.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:        schema.TypeString,
					Required:    true,
					Computed:    false,
					ForceNew:    false,
					Description: "the name of the property key.",
				},
				"value": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    false,
					ForceNew:    false,
					Description: "The value of the property",
				},
				"encrypted_value": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    false,
					ForceNew:    false,
					Description: "The encrypted value of the property",
				},
				"is_secure": {
					Type:     schema.TypeBool,
					Optional: true,
					Computed: false,
					ForceNew: false,
					Description: "Specify whether the given property is secure or not. If true and encrypted_value is not specified, " +
						"GoCD will store the value in encrypted format.",
				},
			},
		},
	}
}

func propertiesSchemaData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    false,
				Description: "the name of the property key.",
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    false,
				Description: "The value of the property",
			},
			"encrypted_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    false,
				Description: "The encrypted value of the property",
			},
			"is_secure": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: false,
				ForceNew: false,
				Description: "Specify whether the given property is secure or not. If true and encrypted_value is not specified, " +
					"GoCD will store the value in encrypted format.",
			},
		},
	}
}

func attributesSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Computed:    false,
		ForceNew:    false,
		MaxItems:    1,
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
					Description: "The domain name for TFS authentication credentials.",
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
					Default:     false,
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
	}
}

func materialSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Computed:    false,
		MaxItems:    1,
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
				"attributes": attributesSchema(),
			},
		},
	}
}

func rulesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The type of entity that the rule is applied on.",
			},
			"directive": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The type of rule which can be either allow or deny.",
			},
			"action": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The action that is being controlled via this rule. Only refer is supported as of now.",
			},
			"resource": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The actual entity on which the rule is applied. Resource should be the name of the entity or a wildcard which matches one or more entities.",
			},
		},
	}
}

func authConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"view": {
				Type:        schema.TypeSet,
				Computed:    true,
				Optional:    true,
				Required:    false,
				Description: "The list of users and roles with view permission for this pipeline group",
				Elem:        usersANdRolesSchema(),
			},
			"operate": {
				Type:        schema.TypeSet,
				Computed:    true,
				Optional:    true,
				Required:    false,
				Description: "The list of users and roles with view permission for this pipeline group",
				Elem:        usersANdRolesSchema(),
			},
			"admins": {
				Type:        schema.TypeSet,
				Computed:    true,
				Optional:    true,
				Required:    false,
				Description: "The list of users and roles with view permission for this pipeline group",
				Elem:        usersANdRolesSchema(),
			},
		},
	}
}

func usersANdRolesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Required:    false,
				Description: "List of users present in GoCD.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"roles": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Required:    false,
				Description: "List of roles present in GoCD.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func retrySchemas() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Optional:    true,
		Computed:    false,
		Description: "Retry configs to be set for the API calls made forG GoCD server.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"count": {
					Type:        schema.TypeInt,
					Optional:    true,
					Computed:    false,
					ForceNew:    false,
					Description: "Number of times to retry in case of API failures.",
				},
				"wait_time": {
					Type:        schema.TypeInt,
					Optional:    true,
					Computed:    false,
					Required:    false,
					Description: "Time interval to wait between subsequent calls",
				},
			},
		},
	}
}
