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
	"fmt"
	"net/url"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

type Config struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Org      string `json:"org"`
	Url      string `json:"url"`
	VDC      string `json:"vdc"`
	Insecure bool   `json:"insecure"`
	Token    string `json:"token"`
}

type VcdService struct { //nolint
	terraformutils.Service
}

func (s *VcdService) buildVDC() (*govcd.Vdc, error) {
	var c Config
	c.User = s.Args["user"].(string)
	c.Password = s.Args["password"].(string)
	c.Org = s.Args["org"].(string)
	c.Url = s.Args["url"].(string)
	c.VDC = s.Args["vdc"].(string)
	c.Insecure = s.Args["allow_unverified_ssl"].(bool)
	c.Token = s.Args["token"].(string)

	u, err := url.ParseRequestURI(c.Url)
	if err != nil {
		return nil, fmt.Errorf("unable to pass url: %s", err)
	}

	client := govcd.NewVCDClient(*u, c.Insecure)

	if s.Args["token"].(string) != "" {
		_ = client.SetToken(c.Org, govcd.AuthorizationHeader, c.Token)
	} else {
		resp, err := client.GetAuthResponse(c.User, c.Password, c.Org)
		if err != nil {
			return nil, fmt.Errorf("unable to authenticate: %s", err)
		}
		fmt.Printf("Token: %s\n", resp.Header[govcd.AuthorizationHeader])
	}

	org, err := client.GetOrgByName(c.Org)
	if err != nil {
		fmt.Printf("organization %s not found : %s\n", c.Org, err)
		return nil, err
	}

	vdc, err := org.GetVDCByName(c.VDC, false)
	if err != nil {
		fmt.Printf("VDC %s not found : %s\n", c.VDC, err)
		return nil, err
	}

	return vdc, nil
}
