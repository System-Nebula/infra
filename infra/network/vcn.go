// Package network provides all the network related stacks, vcn,subnet,seclist,nlb
package network

import (
	"github.com/pulumi/pulumi-oci/sdk/v3/go/oci/core"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"gopkg.in/yaml.v3"
	"os"
)

// NetCfg Struct contains all the needed fields to setup a VCN using pulumi
type NetCfg struct {
	CompartmentID string `yaml:"compartment_id"`
	CidrBlock     string `yaml:"cidr_block"`
	DisplayName   string `yaml:"display_name"`
	Subnets       []struct {
		Name      string `yaml:"name"`
		CidrBlock string `yaml:"cidr_block"`
	} `yaml:"subnets"`
	SecurityLists []struct {
		Type        string `yaml:"type"`
		Protocol    string `yaml:"protocol"`
		Description string `yaml:"description"`
		Destination string `yaml:"destination"`
		Source      string `yaml:"source"`
		Stateless   bool   `yaml:"stateless"`
		TCPOptions  []struct {
			MinPort int `yaml:"min_port"`
			MaxPort int `yaml:"max_port"`
		} `yaml:"tcp_options"`
	} `yaml:"security_lists"`
}

// CreateVCN creates a vcn within oci
func (n *NetCfg) CreateVCN(ctx *pulumi.Context, name string) (*core.Vcn, error) {

	return core.NewVcn(ctx, name, &core.VcnArgs{
		CompartmentId: pulumi.String(n.CompartmentID),
		CidrBlock:     pulumi.String(n.CidrBlock),
		DisplayName:   pulumi.String(n.DisplayName),
	})

}

// LoadFromYaml reads a YAML file and populates the NetCfg struct
func (n *NetCfg) LoadFromYaml(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, n)
}
