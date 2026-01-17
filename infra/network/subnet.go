package network

import (
	"fmt"
	"github.com/pulumi/pulumi-oci/sdk/v3/go/oci/core"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// CreateSubnet creates a subnet within oci using subnet configuration from the Subnets slice
func (n *NetCfg) CreateSubnet(ctx *pulumi.Context, subnetIndex int, vcnID string, seclists []string) (*core.Subnet, error) {
	if subnetIndex < 0 || subnetIndex >= len(n.Subnets) {
		return nil, fmt.Errorf("subnet index %d out of range", subnetIndex)
	}

	subnet := n.Subnets[subnetIndex]
	return core.NewSubnet(ctx, subnet.Name, &core.SubnetArgs{
		CompartmentId:   pulumi.String(n.CompartmentID),
		CidrBlock:       pulumi.String(subnet.CidrBlock),
		DisplayName:     pulumi.String(subnet.Name),
		VcnId:           pulumi.String(vcnID),
		SecurityListIds: pulumi.ToStringArray(seclists),
	})
}

// CreateAllSubnets creates all subnets defined in the Subnets slice
func (n *NetCfg) CreateAllSubnets(ctx *pulumi.Context, vcnID string, seclists []string) ([]*core.Subnet, error) {
	var subnets []*core.Subnet

	for i := range n.Subnets {
		subnet, err := n.CreateSubnet(ctx, i, vcnID, seclists)
		if err != nil {
			return nil, err
		}
		subnets = append(subnets, subnet)
	}

	return subnets, nil
}
