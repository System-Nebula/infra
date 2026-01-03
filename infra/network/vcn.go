// Package network provides all the network related stacks, vcn,subnet,seclist,nlb
package network

import (
	"github.com/pulumi/pulumi-oci/sdk/v3/go/oci/core"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// NetCfg Struct contains all the needed fields to setup a VCN using pulumi
type NetCfg struct {
	CompartmentID string
	CidrBlock     string
	DisplayName   string
}

// CreateVCN creates a vcn within oci
func (n *NetCfg) CreateVCN(ctx *pulumi.Context, name string) (*core.Vcn, error) {

	return core.NewVcn(ctx, name, &core.VcnArgs{
		CompartmentId: pulumi.String(n.CompartmentID),
		CidrBlock:     pulumi.String(n.CidrBlock),
		DisplayName:   pulumi.String(n.DisplayName),
	})

}
