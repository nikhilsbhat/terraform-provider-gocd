package client

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/nikhilsbhat/gocd-sdk-go"
)

type GoCDClient struct {
	gocd.GoCd

	templateClient *resty.Client
}

type PipelineTemplateClient interface {
	CreateTemplateRaw(config map[string]any) (gocd.Template, error)
	UpdateTemplateRaw(name, etag string, config map[string]any) (gocd.Template, error)
	SetRawRetryCount(count int)
	SetRawRetryWaitTime(count int)
}

func newGoCDClient(baseURL string, auth gocd.Auth, logLevel string, caContent []byte) *GoCDClient {
	return &GoCDClient{
		GoCd:           gocd.NewClient(baseURL, auth, logLevel, caContent),
		templateClient: newTemplateClient(baseURL, auth, caContent),
	}
}

func newTemplateClient(baseURL string, auth gocd.Auth, caContent []byte) *resty.Client {
	newClient := resty.New()
	newClient.SetBaseURL(baseURL)

	switch {
	case auth.NoAuth:
	case len(auth.BearerToken) != 0:
		newClient.SetAuthToken(auth.BearerToken)
	default:
		newClient.SetBasicAuth(auth.UserName, auth.Password)
	}

	if len(caContent) != 0 {
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(caContent)
		newClient.SetTLSClientConfig(&tls.Config{RootCAs: certPool}) //nolint:gosec
	} else {
		newClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) //nolint:gosec
	}

	return newClient
}

func (client *GoCDClient) SetRawRetryCount(count int) {
	client.templateClient.SetRetryCount(count)
}

func (client *GoCDClient) SetRawRetryWaitTime(count int) {
	client.templateClient.SetRetryWaitTime(time.Duration(count) * time.Second)
}

func (client *GoCDClient) CreateTemplateRaw(config map[string]any) (gocd.Template, error) {
	var template gocd.Template

	resp, err := client.templateClient.R().
		SetHeaders(map[string]string{
			"Accept":       gocd.HeaderVersionSeven,
			"Content-Type": gocd.ContentJSON,
		}).
		SetBody(config).
		Post(gocd.TemplateConfigEndpoint)
	if err != nil {
		return template, fmt.Errorf("create template '%s': %w", config["name"], err)
	}

	return decodeTemplateResponse(resp, template)
}

func (client *GoCDClient) UpdateTemplateRaw(name, etag string, config map[string]any) (gocd.Template, error) {
	var template gocd.Template

	resp, err := client.templateClient.R().
		SetHeaders(map[string]string{
			"Accept":       gocd.HeaderVersionSeven,
			"Content-Type": gocd.ContentJSON,
			"If-Match":     etag,
		}).
		SetBody(config).
		Put(filepath.Join(gocd.TemplateConfigEndpoint, name))
	if err != nil {
		return template, fmt.Errorf("update template '%s': %w", name, err)
	}

	return decodeTemplateResponse(resp, template)
}

func decodeTemplateResponse(resp *resty.Response, template gocd.Template) (gocd.Template, error) {
	if resp.StatusCode() != http.StatusOK {
		return template, fmt.Errorf("got %d from GoCD while making %s call for %s\nwith BODY:%s",
			resp.StatusCode(), resp.Request.Method, resp.Request.URL, resp.String())
	}

	if err := json.Unmarshal(resp.Body(), &template); err != nil {
		return template, fmt.Errorf("decode template response: %w", err)
	}

	template.ETAG = resp.Header().Get("ETag")

	return template, nil
}
