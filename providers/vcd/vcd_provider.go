// Copyright 2018 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vcd

import (
	"os"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
)

type VcdProvider struct { //nolint
	terraformutils.Provider
	token                string
	user                 string
	password             string
	org                  string
	vdc                  string
	url                  string
	allow_unverified_ssl bool
}

func (p VcdProvider) GetResourceConnections() map[string]map[string][]string {
	return map[string]map[string][]string{
		"grafana_dashboard": {
			"grafana_folder": []string{"folder", "id"},
		},
	}
}

func (p VcdProvider) GetProviderData(arg ...string) map[string]interface{} {
	return map[string]interface{}{
		"provider": map[string]interface{}{
			"vcd": map[string]interface{}{
				"token":                p.token,
				"user":                 p.user,
				"password":             p.password,
				"org":                  p.org,
				"vdc":                  p.vdc,
				"url":                  p.url,
				"allow_unverified_ssl": p.allow_unverified_ssl,
			},
		},
	}
}

func (p *VcdProvider) GetConfig() cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"user":                 cty.StringVal(p.user),
		"token":                cty.StringVal(p.token),
		"password":             cty.StringVal(p.password),
		"org":                  cty.StringVal(p.org),
		"vdc":                  cty.StringVal(p.vdc),
		"url":                  cty.StringVal(p.url),
		"allow_unverified_ssl": cty.BoolVal(p.allow_unverified_ssl),
	})
}

func (p *VcdProvider) Init(args []string) error {
	p.token = os.Getenv("VCD_TOKEN")

	p.user = os.Getenv("VCD_USER")
	if p.token == "" && p.user == "" {
		return errors.New("Grafana API authentication must be set through `GRAFANA_AUTH` env var, either as an API token or as username:password for HTTP basic auth")
	}

	p.password = os.Getenv("VCD_PASSWORD")
	if p.token == "" && p.password == "" {
		return errors.New("Grafana API URL must be set through `GRAFANA_URL` env var")
	}

	p.org = os.Getenv("VCD_ORG")
	if p.org == "" {
		return errors.New("Grafana API URL must be set through `GRAFANA_URL` env var")
	}

	p.vdc = os.Getenv("VCD_VDC")
	if p.vdc == "" {
		return errors.New("Grafana API URL must be set through `GRAFANA_URL` env var")
	}

	p.url = os.Getenv("VCD_URL")
	if p.url == "" {
		return errors.New("Grafana API URL must be set through `GRAFANA_URL` env var")
	}

	if os.Getenv("HTTPS_INSECURE_SKIP_VERIFY") == "1" {
		p.allow_unverified_ssl = true
	}

	return nil
}

func (p *VcdProvider) GetName() string {
	return "vcd"
}

func (p *VcdProvider) InitService(serviceName string, verbose bool) error {
	var isSupported bool
	if _, isSupported = p.GetSupportedService()[serviceName]; !isSupported {
		return errors.New(p.GetName() + ": " + serviceName + " not supported service")
	}

	p.Service = p.GetSupportedService()[serviceName]
	p.Service.SetName(serviceName)
	p.Service.SetVerbose(verbose)
	p.Service.SetProviderName(p.GetName())
	p.Service.SetArgs(map[string]interface{}{
		"token":                p.token,
		"user":                 p.user,
		"password":             p.password,
		"org":                  p.org,
		"vdc":                  p.vdc,
		"url":                  p.url,
		"allow_unverified_ssl": p.allow_unverified_ssl,
	})
	return nil
}

func (p *VcdProvider) GetSupportedService() map[string]terraformutils.ServiceGenerator {
	return map[string]terraformutils.ServiceGenerator{
		"vcd_vapp": &VAppGenerator{},
		// "grafana_folder": &FolderGenerator{},
	}
}
