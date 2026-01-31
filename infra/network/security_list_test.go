package network

import (
	"infra/config"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"testing"
)

type SecurityListMocks int

func (SecurityListMocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return args.Name + "_id", args.Inputs, nil
}

func (SecurityListMocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return args.Args, nil
}

func TestCreateACL(t *testing.T) {
	tests := []struct {
		name          string
		netCfg        NetCfg
		vcnID         string
		expectedError bool
	}{
		{
			name: "Valid security lists creation",
			netCfg: NetCfg{
				NetworkConfig: config.NetworkConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-123",
					},
					CidrBlock:   "10.0.0.0/16",
					DisplayName: "test-vcn",
					Subnets: []config.SubnetConfig{
						{Name: "subnet1", CidrBlock: "10.0.1.0/24"},
					},
					SecurityLists: []config.SecurityListConfig{
						{
							DisplayName: "public-ingress",
							Protocol:    "6",
							Description: "Allow HTTP/HTTPS/SSH access",
							Source:      "0.0.0.0/0",
							Stateless:   false,
							TCPOptions: []config.TCPOptionConfig{
								{MinPort: 22, MaxPort: 22},
								{MinPort: 80, MaxPort: 80},
								{MinPort: 443, MaxPort: 443},
							},
						},
						{
							DisplayName: "public-egress",
							Protocol:    "6",
							Description: "Allow all outbound traffic",
							Destination: "0.0.0.0/0",
							Stateless:   false,
							TCPOptions:  []config.TCPOptionConfig{},
						},
					},
				},
			},
			vcnID:         "vcn-123",
			expectedError: false,
		},
		{
			name: "Empty security lists",
			netCfg: NetCfg{
				NetworkConfig: config.NetworkConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-456",
					},
					CidrBlock:     "192.168.0.0/16",
					DisplayName:   "empty-vcn",
					Subnets:       []config.SubnetConfig{},
					SecurityLists: []config.SecurityListConfig{},
				},
			},
			vcnID:         "vcn-456",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pulumi.RunErr(func(ctx *pulumi.Context) error {
				seclists, err := tt.netCfg.CreateACL(ctx, tt.vcnID)
				if err != nil {
					return err
				}

				if len(seclists) != len(tt.netCfg.SecurityLists) {
					t.Errorf("Expected %d security lists, but got %d", len(tt.netCfg.SecurityLists), len(seclists))
				}

				for i, seclist := range seclists {
					if seclist == nil {
						t.Errorf("Expected security list %d to be created, but got nil", i)
					}
				}

				return nil
			}, pulumi.WithMocks("project", "stack", SecurityListMocks(0)))

			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCreateACLEgressOnly(t *testing.T) {
	netCfg := NetCfg{
		NetworkConfig: config.NetworkConfig{
			BaseConfig: config.BaseConfig{
				CompartmentID: "compartment-123",
			},
			CidrBlock:   "10.0.0.0/16",
			DisplayName: "egress-only-vcn",
			Subnets:     []config.SubnetConfig{},
			SecurityLists: []config.SecurityListConfig{
				{
					DisplayName: "egress-only",
					Protocol:    "6",
					Description: "Allow outbound HTTP traffic",
					Destination: "0.0.0.0/0",
					Stateless:   false,
					TCPOptions: []config.TCPOptionConfig{
						{MinPort: 80, MaxPort: 80},
					},
				},
			},
		},
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		seclists, err := netCfg.CreateACL(ctx, "vcn-123")
		if err != nil {
			return err
		}
		if len(seclists) != 1 {
			t.Errorf("Expected 1 security list, but got %d", len(seclists))
		}
		if seclists[0] == nil {
			t.Error("Expected security list to be created, but got nil")
		}
		return nil
	}, pulumi.WithMocks("project", "stack", SecurityListMocks(0)))

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestCreateACLIngressOnly(t *testing.T) {
	netCfg := NetCfg{
		NetworkConfig: config.NetworkConfig{
			BaseConfig: config.BaseConfig{
				CompartmentID: "compartment-123",
			},
			CidrBlock:   "10.0.0.0/16",
			DisplayName: "ingress-only-vcn",
			Subnets:     []config.SubnetConfig{},
			SecurityLists: []config.SecurityListConfig{
				{
					DisplayName: "ingress-only",
					Protocol:    "6",
					Description: "Allow inbound SSH",
					Source:      "10.0.1.0/24",
					Stateless:   false,
					TCPOptions: []config.TCPOptionConfig{
						{MinPort: 22, MaxPort: 22},
					},
				},
			},
		},
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		seclists, err := netCfg.CreateACL(ctx, "vcn-123")
		if err != nil {
			return err
		}
		if len(seclists) != 1 {
			t.Errorf("Expected 1 security list, but got %d", len(seclists))
		}
		if seclists[0] == nil {
			t.Error("Expected security list to be created, but got nil")
		}
		return nil
	}, pulumi.WithMocks("project", "stack", SecurityListMocks(0)))

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestCreateACLStatelessRules(t *testing.T) {
	netCfg := NetCfg{
		NetworkConfig: config.NetworkConfig{
			BaseConfig: config.BaseConfig{
				CompartmentID: "compartment-123",
			},
			CidrBlock:   "10.0.0.0/16",
			DisplayName: "stateless-vcn",
			Subnets:     []config.SubnetConfig{},
			SecurityLists: []config.SecurityListConfig{
				{
					DisplayName: "stateless-ingress",
					Protocol:    "6",
					Description: "Stateless inbound traffic",
					Source:      "0.0.0.0/0",
					Stateless:   true,
					TCPOptions: []config.TCPOptionConfig{
						{MinPort: 80, MaxPort: 80},
					},
				},
			},
		},
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		seclists, err := netCfg.CreateACL(ctx, "vcn-123")
		if err != nil {
			return err
		}
		if len(seclists) != 1 {
			t.Errorf("Expected 1 security list, but got %d", len(seclists))
		}
		if seclists[0] == nil {
			t.Error("Expected security list to be created, but got nil")
		}
		return nil
	}, pulumi.WithMocks("project", "stack", SecurityListMocks(0)))

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestCreateACLMap(t *testing.T) {
	netCfg := NetCfg{
		NetworkConfig: config.NetworkConfig{
			BaseConfig: config.BaseConfig{
				CompartmentID: "compartment-123",
			},
			CidrBlock:   "10.0.0.0/16",
			DisplayName: "test-vcn",
			Subnets: []config.SubnetConfig{
				{Name: "subnet1", CidrBlock: "10.0.1.0/24"},
			},
			SecurityLists: []config.SecurityListConfig{
				{
					DisplayName: "public-ingress",
					Protocol:    "6",
					Description: "Allow HTTP/HTTPS/SSH access",
					Source:      "0.0.0.0/0",
					SubnetName:  "public-subnet",
					Stateless:   false,
					TCPOptions: []config.TCPOptionConfig{
						{MinPort: 22, MaxPort: 22},
						{MinPort: 80, MaxPort: 80},
						{MinPort: 443, MaxPort: 443},
					},
				},
				{
					DisplayName: "public-egress",
					Protocol:    "6",
					Description: "Allow all outbound traffic",
					Destination: "0.0.0.0/0",
					SubnetName:  "public-subnet",
					Stateless:   false,
					TCPOptions:  []config.TCPOptionConfig{},
				},
				{
					DisplayName: "private-ingress",
					Protocol:    "6",
					Description: "Allow SSH from bastion",
					Source:      "10.0.1.0/24",
					SubnetName:  "private-subnet",
					Stateless:   false,
					TCPOptions: []config.TCPOptionConfig{
						{MinPort: 22, MaxPort: 22},
					},
				},
			},
		},
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		secListMap, err := netCfg.CreateACLMap(ctx, "vcn-123")
		if err != nil {
			return err
		}

		if len(secListMap) != 3 {
			t.Errorf("Expected 3 security lists in map, but got %d", len(secListMap))
		}

		// Verify that each security list exists in the map
		expectedLists := []string{"public-ingress", "public-egress", "private-ingress"}
		for _, name := range expectedLists {
			if _, exists := secListMap[name]; !exists {
				t.Errorf("Expected security list %s to be in map, but it was not found", name)
			}
		}

		for name, sl := range secListMap {
			if sl == nil {
				t.Errorf("Security list %s is nil", name)
			}
		}

		return nil
	}, pulumi.WithMocks("project", "stack", SecurityListMocks(0)))

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
