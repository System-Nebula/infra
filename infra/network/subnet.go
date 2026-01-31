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

// BuildSubnetSecurityListMap creates a mapping of subnet names to security list display names
// This is based on the configuration, not on created resources
func (n *NetCfg) BuildSubnetSecurityListMap() map[string][]string {
	subnetSecLists := make(map[string][]string)

	for _, slConfig := range n.SecurityLists {
		if slConfig.SubnetName != "" {
			subnetSecLists[slConfig.SubnetName] = append(subnetSecLists[slConfig.SubnetName], slConfig.DisplayName)
		}
	}

	return subnetSecLists
}

// GetSecurityListIDsForSubnet returns the security list IDs that should be attached to a specific subnet
// based on the subnet name configuration
func (n *NetCfg) GetSecurityListIDsForSubnet(ctx *pulumi.Context, subnetName string, securityListMap map[string]*core.SecurityList) pulumi.StringArray {
	subnetSecListConfig := n.BuildSubnetSecurityListMap()
	secListNames, hasConfiguredSecLists := subnetSecListConfig[subnetName]

	if !hasConfiguredSecLists || len(secListNames) == 0 {
		// No security lists configured for this subnet
		return nil
	}

	// Find the security list resources by display name and collect their IDs
	secListIDs := make(pulumi.StringArray, 0, len(secListNames))
	for _, slName := range secListNames {
		if sl, exists := securityListMap[slName]; exists {
			secListIDs = append(secListIDs, sl.ID().ToStringOutput())
		}
	}

	if len(secListIDs) == 0 {
		return nil
	}

	return secListIDs
}

// CreateSubnetWithSecurityLists creates a subnet and attaches the appropriate security lists
// based on the subnet name configuration
func (n *NetCfg) CreateSubnetWithSecurityLists(ctx *pulumi.Context, subnetIndex int, vcnID string, securityListMap map[string]*core.SecurityList) (*core.Subnet, error) {
	if subnetIndex < 0 || subnetIndex >= len(n.Subnets) {
		return nil, fmt.Errorf("subnet index %d out of range", subnetIndex)
	}

	subnetConfig := n.Subnets[subnetIndex]
	secListIDs := n.GetSecurityListIDsForSubnet(ctx, subnetConfig.Name, securityListMap)

	return core.NewSubnet(ctx, subnetConfig.Name, &core.SubnetArgs{
		CompartmentId:   pulumi.String(n.CompartmentID),
		CidrBlock:       pulumi.String(subnetConfig.CidrBlock),
		DisplayName:     pulumi.String(subnetConfig.Name),
		VcnId:           pulumi.String(vcnID),
		SecurityListIds: secListIDs,
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

// CreateAllSubnetsWithSecurityLists creates all subnets with their respective security lists
// attached based on subnet name configuration
func (n *NetCfg) CreateAllSubnetsWithSecurityLists(ctx *pulumi.Context, vcnID string, securityListMap map[string]*core.SecurityList) ([]*core.Subnet, error) {
	var subnets []*core.Subnet

	for i := range n.Subnets {
		subnet, err := n.CreateSubnetWithSecurityLists(ctx, i, vcnID, securityListMap)
		if err != nil {
			return nil, fmt.Errorf("failed to create subnet %s: %w", n.Subnets[i].Name, err)
		}
		subnets = append(subnets, subnet)
	}

	return subnets, nil
}
