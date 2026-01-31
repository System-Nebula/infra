package network

import (
	"github.com/pulumi/pulumi-oci/sdk/v3/go/oci/core"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"infra/config"
	"testing"
)

type SubnetMocks int

func (SubnetMocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return args.Name + "_id", args.Inputs, nil
}

func (SubnetMocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return args.Args, nil
}

func TestCreateSubnet(t *testing.T) {
	tests := []struct {
		name          string
		netCfg        NetCfg
		subnetIndex   int
		vcnID         string
		expectedError bool
	}{
		{
			name: "Valid subnet creation",
			netCfg: NetCfg{
				NetworkConfig: config.NetworkConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-123",
					},
					CidrBlock:   "10.0.0.0/16",
					DisplayName: "test-vcn",
					Subnets: []config.SubnetConfig{
						{
							Name:      "my-subnet",
							CidrBlock: "10.0.1.0/24",
						},
					},
					SecurityLists: []config.SecurityListConfig{},
				},
			},
			subnetIndex:   0,
			vcnID:         "vcn-123",
			expectedError: false,
		},
		{
			name: "Subnet with different CIDR",
			netCfg: NetCfg{
				NetworkConfig: config.NetworkConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-456",
					},
					CidrBlock:   "192.168.0.0/16",
					DisplayName: "another-vcn",
					Subnets: []config.SubnetConfig{
						{
							Name:      "another-subnet",
							CidrBlock: "192.168.1.0/24",
						},
					},
					SecurityLists: []config.SecurityListConfig{},
				},
			},
			subnetIndex:   0,
			vcnID:         "vcn-456",
			expectedError: false,
		},
		{
			name: "Multiple subnets - create second subnet",
			netCfg: NetCfg{
				NetworkConfig: config.NetworkConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-789",
					},
					CidrBlock:   "172.16.0.0/16",
					DisplayName: "multi-subnet-vcn",
					Subnets: []config.SubnetConfig{
						{
							Name:      "first-subnet",
							CidrBlock: "172.16.1.0/24",
						},
						{
							Name:      "second-subnet",
							CidrBlock: "172.16.2.0/24",
						},
					},
					SecurityLists: []config.SecurityListConfig{},
				},
			},
			subnetIndex:   1,
			vcnID:         "vcn-789",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pulumi.RunErr(func(ctx *pulumi.Context) error {
				seclists := []string{"seclist-1", "seclist-2"}
				subnet, err := tt.netCfg.CreateSubnet(ctx, tt.subnetIndex, tt.vcnID, seclists)
				if err != nil {
					return err
				}

				if subnet == nil {
					t.Error("Expected subnet to be created, but got nil")
				}

				return nil
			}, pulumi.WithMocks("project", "stack", SubnetMocks(0)))

			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCreateSubnetWithEmptyConfig(t *testing.T) {
	netCfg := NetCfg{
		NetworkConfig: config.NetworkConfig{
			BaseConfig: config.BaseConfig{
				CompartmentID: "",
			},
			CidrBlock:     "",
			DisplayName:   "",
			Subnets:       []config.SubnetConfig{},
			SecurityLists: []config.SecurityListConfig{},
		},
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		seclists := []string{"seclist-1", "seclist-2"}
		subnet, err := netCfg.CreateSubnet(ctx, 0, "vcn-123", seclists)
		if err != nil {
			return err
		}
		if subnet == nil {
			t.Error("Expected subnet to be created, but got nil")
		}
		return nil
	}, pulumi.WithMocks("project", "stack", SubnetMocks(0)))

	if err == nil {
		t.Errorf("Expected error with empty subnets config but got none")
	}
}

func TestCreateSubnetWithInvalidIndex(t *testing.T) {
	netCfg := NetCfg{
		NetworkConfig: config.NetworkConfig{
			BaseConfig: config.BaseConfig{
				CompartmentID: "compartment-123",
			},
			CidrBlock:   "10.0.0.0/16",
			DisplayName: "test-vcn",
			Subnets: []config.SubnetConfig{
				{
					Name:      "test-subnet",
					CidrBlock: "10.0.1.0/24",
				},
			},
			SecurityLists: []config.SecurityListConfig{},
		},
	}

	tests := []struct {
		name        string
		subnetIndex int
	}{
		{"Negative index", -1},
		{"Index out of bounds", 1},
		{"Index way out of bounds", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pulumi.RunErr(func(ctx *pulumi.Context) error {
				seclists := []string{"seclist-1", "seclist-2"}
				subnet, err := netCfg.CreateSubnet(ctx, tt.subnetIndex, "vcn-123", seclists)
				if err != nil {
					return err
				}
				if subnet == nil {
					t.Error("Expected subnet to be created, but got nil")
				}
				return nil
			}, pulumi.WithMocks("project", "stack", SubnetMocks(0)))

			if err == nil {
				t.Errorf("Expected error with invalid index %d but got none", tt.subnetIndex)
			}
		})
	}
}

func TestCreateAllSubnets(t *testing.T) {
	netCfg := NetCfg{
		NetworkConfig: config.NetworkConfig{
			BaseConfig: config.BaseConfig{
				CompartmentID: "compartment-123",
			},
			CidrBlock:   "10.0.0.0/16",
			DisplayName: "test-vcn",
			Subnets: []config.SubnetConfig{
				{
					Name:      "subnet-1",
					CidrBlock: "10.0.1.0/24",
				},
				{
					Name:      "subnet-2",
					CidrBlock: "10.0.2.0/24",
				},
				{
					Name:      "subnet-3",
					CidrBlock: "10.0.3.0/24",
				},
			},
			SecurityLists: []config.SecurityListConfig{},
		},
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		seclists := []string{"seclist-1", "seclist-2"}
		subnets, err := netCfg.CreateAllSubnets(ctx, "vcn-123", seclists)
		if err != nil {
			return err
		}

		if len(subnets) != 3 {
			t.Errorf("Expected 3 subnets, but got %d", len(subnets))
		}

		for i, subnet := range subnets {
			if subnet == nil {
				t.Errorf("Expected subnet %d to be created, but got nil", i)
			}
		}

		return nil
	}, pulumi.WithMocks("project", "stack", SubnetMocks(0)))

	if err != nil {
		t.Errorf("Unexpected error creating all subnets: %v", err)
	}
}

func TestCreateAllSubnetsWithEmptySlice(t *testing.T) {
	netCfg := NetCfg{
		NetworkConfig: config.NetworkConfig{
			BaseConfig: config.BaseConfig{
				CompartmentID: "compartment-123",
			},
			CidrBlock:     "10.0.0.0/16",
			DisplayName:   "test-vcn",
			Subnets:       []config.SubnetConfig{},
			SecurityLists: []config.SecurityListConfig{},
		},
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		seclists := []string{"seclist-1", "seclist-2"}
		subnets, err := netCfg.CreateAllSubnets(ctx, "vcn-123", seclists)
		if err != nil {
			return err
		}

		if len(subnets) != 0 {
			t.Errorf("Expected 0 subnets, but got %d", len(subnets))
		}

		return nil
	}, pulumi.WithMocks("project", "stack", SubnetMocks(0)))

	if err != nil {
		t.Errorf("Unexpected error with empty subnets slice: %v", err)
	}
}

func TestBuildSubnetSecurityListMap(t *testing.T) {
	tests := []struct {
		name          string
		netCfg        NetCfg
		subnetName    string
		expectedCount int
		expectedLists []string
	}{
		{
			name: "Security lists with subnet_name mapping",
			netCfg: NetCfg{
				NetworkConfig: config.NetworkConfig{
					SecurityLists: []config.SecurityListConfig{
						{
							DisplayName: "public-ingress",
							SubnetName:  "public-subnet",
						},
						{
							DisplayName: "public-egress",
							SubnetName:  "public-subnet",
						},
						{
							DisplayName: "private-ingress",
							SubnetName:  "private-subnet",
						},
					},
				},
			},
			subnetName:    "public-subnet",
			expectedCount: 2,
			expectedLists: []string{"public-ingress", "public-egress"},
		},
		{
			name: "Security lists without subnet_name",
			netCfg: NetCfg{
				NetworkConfig: config.NetworkConfig{
					SecurityLists: []config.SecurityListConfig{
						{
							DisplayName: "default-list",
							SubnetName:  "",
						},
					},
				},
			},
			subnetName:    "any-subnet",
			expectedCount: 0,
			expectedLists: nil,
		},
		{
			name: "No security lists",
			netCfg: NetCfg{
				NetworkConfig: config.NetworkConfig{
					SecurityLists: []config.SecurityListConfig{},
				},
			},
			subnetName:    "any-subnet",
			expectedCount: 0,
			expectedLists: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.netCfg.BuildSubnetSecurityListMap()
			secLists, exists := result[tt.subnetName]

			if !exists && len(tt.expectedLists) > 0 {
				t.Errorf("Expected subnet %s to have security lists, but got none", tt.subnetName)
			}

			if exists && len(secLists) != tt.expectedCount {
				t.Errorf("Expected %d security lists for subnet %s, but got %d", tt.expectedCount, tt.subnetName, len(secLists))
			}

			if exists && len(tt.expectedLists) > 0 {
				for _, expected := range tt.expectedLists {
					found := false
					for _, actual := range secLists {
						if actual == expected {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected security list %s not found in result", expected)
					}
				}
			}
		})
	}
}

func TestCreateSubnetWithSecurityLists(t *testing.T) {
	tests := []struct {
		name            string
		netCfg          NetCfg
		subnetIndex     int
		vcnID           string
		securityListMap map[string]*core.SecurityList
		expectedError   bool
	}{
		{
			name: "Subnet with attached security lists",
			netCfg: NetCfg{
				NetworkConfig: config.NetworkConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-123",
					},
					CidrBlock:   "10.0.0.0/16",
					DisplayName: "test-vcn",
					Subnets: []config.SubnetConfig{
						{
							Name:      "public-subnet",
							CidrBlock: "10.0.1.0/24",
						},
					},
					SecurityLists: []config.SecurityListConfig{
						{
							DisplayName: "public-ingress",
							SubnetName:  "public-subnet",
						},
					},
				},
			},
			subnetIndex:   0,
			vcnID:         "vcn-123",
			expectedError: false,
		},
		{
			name: "Subnet without security lists (no mapping)",
			netCfg: NetCfg{
				NetworkConfig: config.NetworkConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-123",
					},
					CidrBlock:   "10.0.0.0/16",
					DisplayName: "test-vcn",
					Subnets: []config.SubnetConfig{
						{
							Name:      "orphan-subnet",
							CidrBlock: "10.0.9.0/24",
						},
					},
					SecurityLists: []config.SecurityListConfig{
						{
							DisplayName: "public-ingress",
							SubnetName:  "public-subnet",
						},
					},
				},
			},
			subnetIndex:   0,
			vcnID:         "vcn-123",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pulumi.RunErr(func(ctx *pulumi.Context) error {
				subnet, err := tt.netCfg.CreateSubnetWithSecurityLists(ctx, tt.subnetIndex, tt.vcnID, tt.securityListMap)
				if err != nil {
					return err
				}

				if subnet == nil {
					t.Error("Expected subnet to be created, but got nil")
				}

				return nil
			}, pulumi.WithMocks("project", "stack", SubnetMocks(0)))

			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCreateAllSubnetsWithSecurityLists(t *testing.T) {
	netCfg := NetCfg{
		NetworkConfig: config.NetworkConfig{
			BaseConfig: config.BaseConfig{
				CompartmentID: "compartment-123",
			},
			CidrBlock:   "10.0.0.0/16",
			DisplayName: "test-vcn",
			Subnets: []config.SubnetConfig{
				{
					Name:      "public-subnet",
					CidrBlock: "10.0.1.0/24",
				},
				{
					Name:      "private-subnet",
					CidrBlock: "10.0.2.0/24",
				},
				{
					Name:      "database-subnet",
					CidrBlock: "10.0.3.0/24",
				},
			},
			SecurityLists: []config.SecurityListConfig{
				{
					DisplayName: "public-ingress",
					SubnetName:  "public-subnet",
				},
				{
					DisplayName: "public-egress",
					SubnetName:  "public-subnet",
				},
				{
					DisplayName: "private-ingress",
					SubnetName:  "private-subnet",
				},
			},
		},
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		subnets, err := netCfg.CreateAllSubnetsWithSecurityLists(ctx, "vcn-123", map[string]*core.SecurityList{})
		if err != nil {
			return err
		}

		if len(subnets) != 3 {
			t.Errorf("Expected 3 subnets, but got %d", len(subnets))
		}

		for i, subnet := range subnets {
			if subnet == nil {
				t.Errorf("Expected subnet %d to be created, but got nil", i)
			}
		}

		return nil
	}, pulumi.WithMocks("project", "stack", SubnetMocks(0)))

	if err != nil {
		t.Errorf("Unexpected error creating all subnets with security lists: %v", err)
	}
}
