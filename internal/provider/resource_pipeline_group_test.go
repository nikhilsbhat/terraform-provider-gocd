//nolint:testpackage
package provider

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func TestFlattenPipelineGroupAuthorizationConfigPreservesConfiguredEmptyLists(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourcePipelineGroup().Schema, map[string]any{
		utils.TerraformResourceName: "AI",
		utils.TerraformResourceAuthorization: []any{
			map[string]any{
				utils.TerraformResourceView: []any{
					map[string]any{
						utils.TerraformResourceRoles: []any{"devops", "ai"},
					},
				},
				utils.TerraformResourceOperate: []any{
					map[string]any{
						utils.TerraformResourceRoles: []any{"devops", "ai"},
					},
				},
				utils.TerraformResourceAdmins: []any{
					map[string]any{
						utils.TerraformResourceUsers: []any{},
					},
				},
			},
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

	got := flattenPipelineGroupAuthorizationConfig(pipelineGroup, resourceData)
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
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected flattened authorization\nwant: %#v\n got: %#v", want, got)
	}
}

func TestRawAuthConfigFieldWasSet(t *testing.T) {
	authConfig := []cty.Value{
		cty.ObjectVal(map[string]cty.Value{
			utils.TerraformResourceView: cty.TupleVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					utils.TerraformResourceRoles: cty.TupleVal([]cty.Value{cty.StringVal("devops"), cty.StringVal("ai")}),
					utils.TerraformResourceUsers: cty.NullVal(cty.DynamicPseudoType),
				}),
			}),
			utils.TerraformResourceOperate: cty.TupleVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					utils.TerraformResourceRoles: cty.TupleVal([]cty.Value{cty.StringVal("devops"), cty.StringVal("ai")}),
					utils.TerraformResourceUsers: cty.NullVal(cty.DynamicPseudoType),
				}),
			}),
			utils.TerraformResourceAdmins: cty.TupleVal([]cty.Value{
				cty.ObjectVal(map[string]cty.Value{
					utils.TerraformResourceRoles: cty.NullVal(cty.DynamicPseudoType),
					utils.TerraformResourceUsers: cty.EmptyTupleVal,
				}),
			}),
		}),
	}

	if rawAuthConfigFieldWasSet(authConfig, utils.TerraformResourceView, utils.TerraformResourceUsers) {
		t.Fatal("view.users should be treated as omitted")
	}

	if !rawAuthConfigFieldWasSet(authConfig, utils.TerraformResourceView, utils.TerraformResourceRoles) {
		t.Fatal("view.roles should be treated as configured")
	}

	if !rawAuthConfigFieldWasSet(authConfig, utils.TerraformResourceAdmins, utils.TerraformResourceUsers) {
		t.Fatal("admins.users should be treated as explicitly configured")
	}
}

func TestPipelineGroupPipelinesAreSetBased(t *testing.T) {
	resource := resourcePipelineGroup()
	pipelinesSchema := resource.Schema[utils.TerraformResourcePipelines]

	if pipelinesSchema.Type != schema.TypeSet {
		t.Fatalf("pipeline group pipelines should be set based, got %v", pipelinesSchema.Type)
	}

	first := schema.NewSet(schema.HashString, []any{"helm-images", "helm-drift"})
	second := schema.NewSet(schema.HashString, []any{"helm-drift", "helm-images"})

	if !first.Equal(second) {
		t.Fatal("pipeline ordering should not change pipeline group state")
	}

	got := getPipelines(first)
	want := []gocd.Pipeline{
		{Name: "helm-drift"},
		{Name: "helm-images"},
	}

	if !pipelineNamesEqual(got, want) {
		t.Fatalf("unexpected pipelines\nwant: %#v\n got: %#v", want, got)
	}
}

func pipelineNamesEqual(left, right []gocd.Pipeline) bool {
	if len(left) != len(right) {
		return false
	}

	leftNames := make(map[string]struct{}, len(left))
	for _, pipeline := range left {
		leftNames[pipeline.Name] = struct{}{}
	}

	for _, pipeline := range right {
		if _, ok := leftNames[pipeline.Name]; !ok {
			return false
		}
	}

	return true
}
