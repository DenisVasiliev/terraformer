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
package cmd

import (
	"github.com/DenisVasiliev/terraformer/providers/vcd"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/spf13/cobra"
)

func newCmdVcdImporter(options ImportOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vcd",
		Short: "Import current state to Terraform configuration from VMWare VCD",
		Long:  "Import current state to Terraform configuration from VMWare VCD",
		RunE: func(cmd *cobra.Command, args []string) error {
			provider := newVcdProvider()
			err := Import(provider, options, []string{})
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.AddCommand(listCmd(newVcdProvider()))
	baseProviderFlags(cmd.PersistentFlags(), &options, "vm", "vm=vm1")
	return cmd
}

func newVcdProvider() terraformutils.ProviderGenerator {
	return &vcd.VcdProvider{}
}
