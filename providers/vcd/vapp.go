package vcd

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

type VAppGenerator struct {
	VcdService
}

func (v *VAppGenerator) InitResources() error {
	vdc, err := v.buildVDC()
	if err != nil {
		return fmt.Errorf("unable to build vcd client: %v", err)
	}

	err = v.createVappResources(vdc)
	if err != nil {
		return err
	}

	return nil
}

func (g *VAppGenerator) createVappResources(vdc *govcd.Vdc) error {
	vapps := vdc.GetVappList()

	for _, vapp := range vapps {
		// search result doesn't include slug, so need to look up dashboard.

		resource := terraformutils.NewResource(
			vapp.ID,
			vapp.Name,
			"vdc_vapp",
			"vdc",
			map[string]string{},
			[]string{},
			map[string]interface{}{},
		)
		resource.DataFiles = map[string][]byte{}

		g.Resources = append(g.Resources, resource)
	}

	return nil
}
