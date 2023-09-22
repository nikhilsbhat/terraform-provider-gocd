package provider

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-cli/pkg/render"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func resourcePipeline() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelineCreate,
		ReadContext:   resourcePipelineRead,
		UpdateContext: resourcePipelineUpdate,
		DeleteContext: resourcePipelineDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The name of the pipeline to be created (this should be the same that would be passed under `config`).",
			},
			"group": {
				Type:        schema.TypeString,
				Optional:    false,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "Name of the pipeline group that this pipeline should be part of.",
			},
			// "no_validate": {
			//	Type:        schema.TypeBool,
			//	Optional:    true,
			//	Computed:    false,
			//	ForceNew:    true,
			//	Description: "Enabling this would disable the pipeline syntax validation",
			// },
			// "plugin_config": {
			//	Type:        schema.TypeList,
			//	Optional:    true,
			//	Computed:    false,
			//	ForceNew:    true,
			//	MaxItems:    1,
			//	Description: "Configurations of the plugin against which the pipeline has to be validated.",
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"version": {
			//				Type:          schema.TypeString,
			//				Optional:      true,
			//				Computed:      false,
			//				ForceNew:      false,
			//				Description:   "Version of the plugin that needs to be downloaded for validating pipeline syntax",
			//				ConflictsWith: []string{"plugin_config.0.path", "plugin_config.0.url"},
			//			},
			//			"path": {
			//				Type:          schema.TypeString,
			//				Optional:      true,
			//				Computed:      false,
			//				ForceNew:      false,
			//				Description:   "Local path to plugin if in case the plugin jar is present locally",
			//				ConflictsWith: []string{"plugin_config.0.version", "plugin_config.0.url"},
			//			},
			//			"url": {
			//				Type:          schema.TypeString,
			//				Optional:      true,
			//				Computed:      false,
			//				ForceNew:      false,
			//				Description:   "URL to download the plugin which should be used to validate the pipeline config +
			//				(ex: https://github.com/tomzo/gocd-yaml-config-plugin/releases/download/0.14.1/yaml-config-plugin-0.14.1.jar)",
			//				ConflictsWith: []string{"plugin_config.0.path", "plugin_config.0.version"},
			//				ValidateFunc: validation.All(
			//					validation.IsURLWithHTTPorHTTPS,
			//				),
			//			},
			//		},
			//	},
			// },
			"config": {
				Type:        schema.TypeString,
				Optional:    false,
				Required:    true,
				Computed:    false,
				ForceNew:    false,
				Description: "The config of the pipeline to be created (it can take in yaml/json data based on the attribute set).",
			},
			"pause_on_creation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "Enabling this would have the pipeline paused on creation",
			},
			"pause_reason": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "Reason for pausing the pipeline on start",
			},
			"yaml": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Would be set to true when pipeline config declared under `config` is of type yaml.",
			},
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Etag used to track the pipeline config",
			},
		},
	}
}

func resourcePipelineCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !d.IsNewResource() {
		return nil
	}

	id := d.Id()

	if len(id) == 0 {
		resourceID := utils.String(d.Get(utils.TerraformResourceName))
		id = resourceID
	}

	pipelineCfg := gocd.PipelineConfig{
		Name:          id,
		Group:         utils.String(d.Get(utils.TerraformResourceGroup)),
		PausePipeline: utils.Bool(d.Get(utils.TerraformResourcePauseOnCreation)),
		PauseReason:   utils.String(d.Get(utils.TerraformResourcePauseReason)),
	}

	obj := render.Object(utils.String(d.Get(utils.TerraformResourceConfig)))
	logger := logrus.New()
	obj.CheckFileType(logger)

	var configMap map[string]interface{}
	switch objType := obj.CheckFileType(logger); objType {
	case render.FileTypeJSON:
		if err := json.Unmarshal([]byte(obj.String()), &configMap); err != nil {
			return diag.Errorf("decoding pipeline config errored with: %v", err)
		}
		pipelineCfg.Config = configMap
	case render.FileTypeYAML:
		if err := yaml.Unmarshal([]byte(obj.String()), &configMap); err != nil {
			return diag.Errorf("decoding pipeline config errored with: %v", err)
		}
		if err := d.Set(utils.TerraformResourceYAML, true); err != nil {
			return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceYAML, err)
		}

		pipelineCfg.Config = configMap
	default:
		return diag.Errorf("pipeline config type is unknown")
	}

	if pipelineCfg.Config["name"] != id {
		return diag.Errorf("pipeline name passed under attribute and pipeline config are not same, make sure to pass the same values, "+
			"current values: 'attribute:%s config:%s'", id, pipelineCfg.Config["name"].(string))
	}

	if _, err := defaultConfig.CreatePipeline(pipelineCfg); err != nil {
		return diag.Errorf("creating pipeline '%s' errored with: %v", id, err)
	}

	d.SetId(id)

	return resourcePipelineRead(ctx, d, meta)
}

func resourcePipelineRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	name := utils.String(d.Get(utils.TerraformResourceName))
	response, err := defaultConfig.GetPipelineConfig(name)
	if err != nil {
		return diag.Errorf("getting pipeline config %s errored with: %v", name, err)
	}

	if err = d.Set(utils.TerraformResourceEtag, response.ETAG); err != nil {
		return diag.Errorf(settingAttrErrorTmp, utils.TerraformResourceEtag, err)
	}

	return nil
}

func resourcePipelineUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	if !d.HasChanges(utils.TerraformResourceConfig) {
		log.Printf("nothing to update so skipping")

		return nil
	}

	pipelineCfg := gocd.PipelineConfig{
		Name: utils.String(d.Get(utils.TerraformResourceName)),
		ETAG: utils.String(d.Get(utils.TerraformResourceEtag)),
	}

	obj := render.Object(utils.String(d.Get(utils.TerraformResourceConfig)))
	logger := logrus.New()
	obj.CheckFileType(logger)

	var configMap map[string]interface{}
	switch objType := obj.CheckFileType(logger); objType {
	case render.FileTypeJSON:
		if err := json.Unmarshal([]byte(obj.String()), &configMap); err != nil {
			return diag.Errorf("decoding pipeline config errored with: %v", err)
		}
		pipelineCfg.Config = configMap
	case render.FileTypeYAML:
		if err := yaml.Unmarshal([]byte(obj.String()), &configMap); err != nil {
			return diag.Errorf("decoding pipeline config errored with: %v", err)
		}
		pipelineCfg.Config = configMap
	default:
		return diag.Errorf("pipeline config type is unknown")
	}
	pipelineCfg.Config["group"] = utils.String(d.Get(utils.TerraformResourceGroup))

	if _, err := defaultConfig.UpdatePipelineConfig(pipelineCfg); err != nil {
		return diag.Errorf("updating pipeline '%s' errored with: %v", pipelineCfg.Name, err)
	}

	return resourcePipelineRead(ctx, d, meta)
}

func resourcePipelineDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaultConfig := meta.(gocd.GoCd)

	id := d.Id()
	if len(d.Id()) == 0 {
		return diag.Errorf("resource with the ID '%s' not found", id)
	}

	name := utils.String(d.Get(utils.TerraformResourceName))

	err := defaultConfig.DeletePipeline(name)
	if err != nil {
		return diag.Errorf("deleting pipeline %s errored with: %v", name, err)
	}

	d.SetId("")

	return nil
}
