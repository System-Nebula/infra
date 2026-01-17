package network

import (
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
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
				CompartmentID: "compartment-123",
				CidrBlock:     "10.0.0.0/16",
				DisplayName:   "test-vcn",
				Subnets: []struct {
					Name      string `yaml:"name"`
					CidrBlock string `yaml:"cidr_block"`
				}{
					{
						Name:      "my-subnet",
						CidrBlock: "10.0.1.0/24",
					},
				},
				SecurityLists: []struct {
					DisplayName string `yaml:"display_name"`
					Protocol    string `yaml:"protocol"`
					Description string `yaml:"description"`
					Destination string `yaml:"destination"`
					Source      string `yaml:"source"`
					Stateless   bool   `yaml:"stateless"`
					TCPOptions  []struct {
						MinPort int `yaml:"min_port"`
						MaxPort int `yaml:"max_port"`
					} `yaml:"tcp_options"`
				}{},
			},
			subnetIndex:   0,
			vcnID:         "vcn-123",
			expectedError: false,
		},
		{
			name: "Subnet with different CIDR",
			netCfg: NetCfg{
				CompartmentID: "compartment-456",
				CidrBlock:     "192.168.0.0/16",
				DisplayName:   "another-vcn",
				Subnets: []struct {
					Name      string `yaml:"name"`
					CidrBlock string `yaml:"cidr_block"`
				}{
					{
						Name:      "another-subnet",
						CidrBlock: "192.168.1.0/24",
					},
				},
				SecurityLists: []struct {
					DisplayName string `yaml:"display_name"`
					Protocol    string `yaml:"protocol"`
					Description string `yaml:"description"`
					Destination string `yaml:"destination"`
					Source      string `yaml:"source"`
					Stateless   bool   `yaml:"stateless"`
					TCPOptions  []struct {
						MinPort int `yaml:"min_port"`
						MaxPort int `yaml:"max_port"`
					} `yaml:"tcp_options"`
				}{},
			},
			subnetIndex:   0,
			vcnID:         "vcn-456",
			expectedError: false,
		},
		{
			name: "Multiple subnets - create second subnet",
			netCfg: NetCfg{
				CompartmentID: "compartment-789",
				CidrBlock:     "172.16.0.0/16",
				DisplayName:   "multi-subnet-vcn",
				Subnets: []struct {
					Name      string `yaml:"name"`
					CidrBlock string `yaml:"cidr_block"`
				}{
					{
						Name:      "first-subnet",
						CidrBlock: "172.16.1.0/24",
					},
					{
						Name:      "second-subnet",
						CidrBlock: "172.16.2.0/24",
					},
				},
				SecurityLists: []struct {
					DisplayName string `yaml:"display_name"`
					Protocol    string `yaml:"protocol"`
					Description string `yaml:"description"`
					Destination string `yaml:"destination"`
					Source      string `yaml:"source"`
					Stateless   bool   `yaml:"stateless"`
					TCPOptions  []struct {
						MinPort int `yaml:"min_port"`
						MaxPort int `yaml:"max_port"`
					} `yaml:"tcp_options"`
				}{},
			},
			subnetIndex:   1,
			vcnID:         "vcn-789",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pulumi.RunErr(func(ctx *pulumi.Context) error {
				subnet, err := tt.netCfg.CreateSubnet(ctx, tt.subnetIndex, tt.vcnID)
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
		CompartmentID: "",
		CidrBlock:     "",
		DisplayName:   "",
		Subnets: []struct {
			Name      string `yaml:"name"`
			CidrBlock string `yaml:"cidr_block"`
		}{},
		SecurityLists: []struct {
			DisplayName string `yaml:"display_name"`
			Protocol    string `yaml:"protocol"`
			Description string `yaml:"description"`
			Destination string `yaml:"destination"`
			Source      string `yaml:"source"`
			Stateless   bool   `yaml:"stateless"`
			TCPOptions  []struct {
				MinPort int `yaml:"min_port"`
				MaxPort int `yaml:"max_port"`
			} `yaml:"tcp_options"`
		}{},
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		subnet, err := netCfg.CreateSubnet(ctx, 0, "vcn-123")
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
		CompartmentID: "compartment-123",
		CidrBlock:     "10.0.0.0/16",
		DisplayName:   "test-vcn",
		Subnets: []struct {
			Name      string `yaml:"name"`
			CidrBlock string `yaml:"cidr_block"`
		}{
			{
				Name:      "test-subnet",
				CidrBlock: "10.0.1.0/24",
			},
		},
		SecurityLists: []struct {
			DisplayName string `yaml:"display_name"`
			Protocol    string `yaml:"protocol"`
			Description string `yaml:"description"`
			Destination string `yaml:"destination"`
			Source      string `yaml:"source"`
			Stateless   bool   `yaml:"stateless"`
			TCPOptions  []struct {
				MinPort int `yaml:"min_port"`
				MaxPort int `yaml:"max_port"`
			} `yaml:"tcp_options"`
		}{},
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
				subnet, err := netCfg.CreateSubnet(ctx, tt.subnetIndex, "vcn-123")
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
		CompartmentID: "compartment-123",
		CidrBlock:     "10.0.0.0/16",
		DisplayName:   "test-vcn",
		Subnets: []struct {
			Name      string `yaml:"name"`
			CidrBlock string `yaml:"cidr_block"`
		}{
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
		SecurityLists: []struct {
			DisplayName string `yaml:"display_name"`
			Protocol    string `yaml:"protocol"`
			Description string `yaml:"description"`
			Destination string `yaml:"destination"`
			Source      string `yaml:"source"`
			Stateless   bool   `yaml:"stateless"`
			TCPOptions  []struct {
				MinPort int `yaml:"min_port"`
				MaxPort int `yaml:"max_port"`
			} `yaml:"tcp_options"`
		}{},
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		subnets, err := netCfg.CreateAllSubnets(ctx, "vcn-123")
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
		CompartmentID: "compartment-123",
		CidrBlock:     "10.0.0.0/16",
		DisplayName:   "test-vcn",
		Subnets: []struct {
			Name      string `yaml:"name"`
			CidrBlock string `yaml:"cidr_block"`
		}{},
		SecurityLists: []struct {
			DisplayName string `yaml:"display_name"`
			Protocol    string `yaml:"protocol"`
			Description string `yaml:"description"`
			Destination string `yaml:"destination"`
			Source      string `yaml:"source"`
			Stateless   bool   `yaml:"stateless"`
			TCPOptions  []struct {
				MinPort int `yaml:"min_port"`
				MaxPort int `yaml:"max_port"`
			} `yaml:"tcp_options"`
		}{},
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		subnets, err := netCfg.CreateAllSubnets(ctx, "vcn-123")
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
