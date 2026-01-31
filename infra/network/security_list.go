package network

import (
	"github.com/pulumi/pulumi-oci/sdk/v3/go/oci/core"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// CreateACL Constructor function that creates the Ingress/Egress Security Lists required
func (n *NetCfg) CreateACL(ctx *pulumi.Context, vcnID string) ([]*core.SecurityList, error) {
	var seclists []*core.SecurityList
	for _, v := range n.SecurityLists {
		var egressRules core.SecurityListEgressSecurityRuleArray
		var ingressRules core.SecurityListIngressSecurityRuleArray

		if v.Destination != "" && v.Destination != "null" {
			egressRule := core.SecurityListEgressSecurityRuleArgs{
				Protocol:    pulumi.String(v.Protocol),
				Destination: pulumi.String(v.Destination),
				Description: pulumi.String(v.Description),
				Stateless:   pulumi.Bool(v.Stateless),
			}

			if len(v.TCPOptions) > 0 {
				for _, tcp := range v.TCPOptions {
					egressRule.TcpOptions = &core.SecurityListEgressSecurityRuleTcpOptionsArgs{
						Min: pulumi.Int(tcp.MinPort),
						Max: pulumi.Int(tcp.MaxPort),
					}
				}
			}

			egressRules = append(egressRules, egressRule)
		}

		if v.Source != "" && v.Source != "null" {
			ingressRule := core.SecurityListIngressSecurityRuleArgs{
				Protocol:    pulumi.String(v.Protocol),
				Source:      pulumi.String(v.Source),
				Description: pulumi.String(v.Description),
				Stateless:   pulumi.Bool(v.Stateless),
			}

			if len(v.TCPOptions) > 0 {
				for _, tcp := range v.TCPOptions {
					ingressRule.TcpOptions = &core.SecurityListIngressSecurityRuleTcpOptionsArgs{
						Min: pulumi.Int(tcp.MinPort),
						Max: pulumi.Int(tcp.MaxPort),
					}
				}
			}

			ingressRules = append(ingressRules, ingressRule)
		}

		sec, err := core.NewSecurityList(ctx, v.DisplayName, &core.SecurityListArgs{
			CompartmentId:        pulumi.String(n.CompartmentID),
			VcnId:                pulumi.String(vcnID),
			DisplayName:          pulumi.String(v.DisplayName),
			EgressSecurityRules:  egressRules,
			IngressSecurityRules: ingressRules,
		})
		if err != nil {
			return nil, err
		}
		seclists = append(seclists, sec)
	}
	return seclists, nil
}

// CreateACLMap creates security lists and returns a map of display names to security list resources
// This enables security lists to be easily referenced by their display names when attaching to subnets
func (n *NetCfg) CreateACLMap(ctx *pulumi.Context, vcnID string) (map[string]*core.SecurityList, error) {
	secListMap := make(map[string]*core.SecurityList)

	for _, v := range n.SecurityLists {
		var egressRules core.SecurityListEgressSecurityRuleArray
		var ingressRules core.SecurityListIngressSecurityRuleArray

		if v.Destination != "" && v.Destination != "null" {
			egressRule := core.SecurityListEgressSecurityRuleArgs{
				Protocol:    pulumi.String(v.Protocol),
				Destination: pulumi.String(v.Destination),
				Description: pulumi.String(v.Description),
				Stateless:   pulumi.Bool(v.Stateless),
			}

			if len(v.TCPOptions) > 0 {
				for _, tcp := range v.TCPOptions {
					egressRule.TcpOptions = &core.SecurityListEgressSecurityRuleTcpOptionsArgs{
						Min: pulumi.Int(tcp.MinPort),
						Max: pulumi.Int(tcp.MaxPort),
					}
				}
			}

			egressRules = append(egressRules, egressRule)
		}

		if v.Source != "" && v.Source != "null" {
			ingressRule := core.SecurityListIngressSecurityRuleArgs{
				Protocol:    pulumi.String(v.Protocol),
				Source:      pulumi.String(v.Source),
				Description: pulumi.String(v.Description),
				Stateless:   pulumi.Bool(v.Stateless),
			}

			if len(v.TCPOptions) > 0 {
				for _, tcp := range v.TCPOptions {
					ingressRule.TcpOptions = &core.SecurityListIngressSecurityRuleTcpOptionsArgs{
						Min: pulumi.Int(tcp.MinPort),
						Max: pulumi.Int(tcp.MaxPort),
					}
				}
			}

			ingressRules = append(ingressRules, ingressRule)
		}

		sec, err := core.NewSecurityList(ctx, v.DisplayName, &core.SecurityListArgs{
			CompartmentId:        pulumi.String(n.CompartmentID),
			VcnId:                pulumi.String(vcnID),
			DisplayName:          pulumi.String(v.DisplayName),
			EgressSecurityRules:  egressRules,
			IngressSecurityRules: ingressRules,
		})
		if err != nil {
			return nil, err
		}
		secListMap[v.DisplayName] = sec
	}

	return secListMap, nil
}
