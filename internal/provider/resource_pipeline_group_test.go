//nolint:testpackage
package provider

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func TestFlattenPipelineGroupAuthorizationConfigPreservesConfiguredEmptyLists(t *testing.T) {
	configuredAuth := schema.NewSet(schema.HashResource(authConfigSchema()), []any{
		map[string]any{
			utils.TerraformResourceView: schema.NewSet(schema.HashResource(usersANdRolesSchema()), []any{
				map[string]any{
					utils.TerraformResourceRoles: []any{"devops", "ai"},
				},
			}),
			utils.TerraformResourceOperate: schema.NewSet(schema.HashResource(usersANdRolesSchema()), []any{
				map[string]any{
					utils.TerraformResourceRoles: []any{"devops", "ai"},
				},
			}),
			utils.TerraformResourceAdmins: schema.NewSet(schema.HashResource(usersANdRolesSchema()), []any{
				map[string]any{
					utils.TerraformResourceUsers: []any{},
				},
			}),
		},
	})

	pipelineGroup := gocd.PipelineGroup{
		Authorization: gocd.PipelineGroupAuthorizationConfig{
			View: gocd.AuthorizationConfig{
				Roles: []string{"devops", "ai"},
			},
			Operate: gocd.AuthorizationConfig{
				Roles: []string{"devops", "ai"},
			},
		},
	}

	got := flattenPipelineGroupAuthorizationConfig(pipelineGroup, configuredAuth)
	want := []any{
		map[string]any{
			utils.TerraformResourceView: []any{
				map[string]any{
					utils.TerraformResourceRoles: []string{"devops", "ai"},
				},
			},
			utils.TerraformResourceOperate: []any{
				map[string]any{
					utils.TerraformResourceRoles: []string{"devops", "ai"},
				},
			},
			utils.TerraformResourceAdmins: []any{
				map[string]any{
					utils.TerraformResourceUsers: []string(nil),
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected flattened authorization\nwant: %#v\n got: %#v", want, got)
	}
}
