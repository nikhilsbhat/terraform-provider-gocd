package client

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
)

func GetGoCDClient(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	clientCfg := struct {
		url         string
		username    string
		password    string
		bearerToken string
		loglevel    string
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

	gocdAuth := gocd.Auth{
		UserName:    clientCfg.username,
		Password:    clientCfg.password,
		BearerToken: clientCfg.bearerToken,
	}

	goCDClient := gocd.NewClient(clientCfg.url, gocdAuth, clientCfg.loglevel, clientCfg.ca)

	return goCDClient, nil
}
