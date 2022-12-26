package client

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nikhilsbhat/gocd-sdk-go"
)

func GetGoCDClient(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	clientCfg := struct {
		url      string
		username string
		password string
		loglevel string
		ca       []byte
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

	goCDClient := gocd.NewClient(clientCfg.url, clientCfg.username, clientCfg.password, clientCfg.loglevel, clientCfg.ca)
	return goCDClient, nil
}
