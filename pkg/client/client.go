package client

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
	goErr "github.com/nikhilsbhat/gocd-sdk-go/pkg/errors"
	"github.com/nikhilsbhat/terraform-provider-gocd/pkg/utils"
)

func GetGoCDClient(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	clientCfg := struct {
		url         string
		username    string
		password    string
		bearerToken string
		loglevel    string
		skipCheck   bool
		ca          []byte
	}{}

	if baseURL := d.Get("base_url").(string); len(baseURL) == 0 {
		diag.Errorf("'base_url' was not set")
	} else {
		clientCfg.url = baseURL
	}

	if username := d.Get("username").(string); len(username) == 0 {
		diag.Errorf("'username' was not set")
	} else {
		clientCfg.username = username
	}

	if password := d.Get("password").(string); len(password) == 0 {
		diag.Errorf("'password' was not set")
	} else {
		clientCfg.password = password
	}

	if authToken, ok := d.GetOk("auth_token"); !ok {
		diag.Errorf("'auth_token' was not set")
	} else {
		clientCfg.bearerToken = authToken.(string)
	}

	if skipCheck, ok := d.GetOk("skip_check"); !ok {
		diag.Errorf("'skip_check' was not set")
	} else {
		clientCfg.skipCheck = skipCheck.(bool)
	}

	if caFileContent := d.Get("ca_file").(string); len(caFileContent) == 0 {
		diag.Errorf("'ca_file' was not set")
	} else {
		clientCfg.ca = []byte(caFileContent)
	}

	if loglevel := d.Get("loglevel").(string); len(loglevel) == 0 {
		clientCfg.loglevel = "info"
	} else {
		clientCfg.loglevel = loglevel
	}

	goCDAuth := gocd.Auth{
		UserName:    clientCfg.username,
		Password:    clientCfg.password,
		BearerToken: clientCfg.bearerToken,
	}

	goCDClient := gocd.NewClient(clientCfg.url, goCDAuth, clientCfg.loglevel, clientCfg.ca)

	retryConfigs := getRetryConfig(d.Get(utils.TerraformResourceRetries))
	if retryConfigs.count != 0 {
		log.Printf("setting API retry count to %d:\n", retryConfigs.count)
		goCDClient.SetRetryCount(retryConfigs.count)
	}

	if retryConfigs.waitTime != 0 {
		log.Printf("setting API retry wait time to %d:\n", retryConfigs.waitTime)
		goCDClient.SetRetryWaitTime(retryConfigs.waitTime)
	}

	if !clientCfg.skipCheck {
		if _, err := goCDClient.GetServerHealth(); err != nil {
			if !errors.Is(err, goErr.MarshalError{}) {
				return nil, diag.Errorf("errored while connecting to server\nerror: %v\nkindly re-check the baseURL and authorization config before rerunning plan again", err)
			}
		}
	}

	return goCDClient, nil
}

type retryConfig struct {
	count    int
	waitTime int
}

func getRetryConfig(retryConfigs interface{}) retryConfig {
	retrySet := retryConfigs.(*schema.Set).List()
	if len(retrySet) == 0 {
		return retryConfig{
			count:    0,
			waitTime: 0,
		}
	}

	flattenedRetryConfigs := retrySet[0].(map[string]interface{})

	return retryConfig{
		count:    flattenedRetryConfigs[utils.TerraformResourceCount].(int),
		waitTime: flattenedRetryConfigs[utils.TerraformResourceWaitTime].(int),
	}
}
