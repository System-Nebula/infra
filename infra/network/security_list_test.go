package network

import (
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
				CompartmentID: "compartment-123",
				CidrBlock:     "10.0.0.0/16",
				DisplayName:   "test-vcn",
				Subnets: []struct {
					Name      string `yaml:"name"`
					CidrBlock string `yaml:"cidr_block"`
				}{
					{Name: "subnet1", CidrBlock: "10.0.1.0/24"},
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
				}{
					{
						DisplayName: "public-ingress",
						Protocol:    "6",
						Description: "Allow HTTP/HTTPS/SSH access",
						Source:      "0.0.0.0/0",
						Stateless:   false,
						TCPOptions: []struct {
							MinPort int `yaml:"min_port"`
							MaxPort int `yaml:"max_port"`
						}{
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
						TCPOptions: []struct {
							MinPort int `yaml:"min_port"`
							MaxPort int `yaml:"max_port"`
						}{},
					},
				},
			},
			vcnID:         "vcn-123",
			expectedError: false,
		},
		{
			name: "Empty security lists",
			netCfg: NetCfg{
				CompartmentID: "compartment-456",
				CidrBlock:     "192.168.0.0/16",
				DisplayName:   "empty-vcn",
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
		CompartmentID: "compartment-123",
		CidrBlock:     "10.0.0.0/16",
		DisplayName:   "egress-only-vcn",
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
		}{
			{
				DisplayName: "egress-only",
				Protocol:    "6",
				Description: "Allow outbound HTTP traffic",
				Destination: "0.0.0.0/0",
				Stateless:   false,
				TCPOptions: []struct {
					MinPort int `yaml:"min_port"`
					MaxPort int `yaml:"max_port"`
				}{
					{MinPort: 80, MaxPort: 80},
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
		CompartmentID: "compartment-123",
		CidrBlock:     "10.0.0.0/16",
		DisplayName:   "ingress-only-vcn",
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
		}{
			{
				DisplayName: "ingress-only",
				Protocol:    "6",
				Description: "Allow inbound SSH",
				Source:      "10.0.1.0/24",
				Stateless:   false,
				TCPOptions: []struct {
					MinPort int `yaml:"min_port"`
					MaxPort int `yaml:"max_port"`
				}{
					{MinPort: 22, MaxPort: 22},
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
		CompartmentID: "compartment-123",
		CidrBlock:     "10.0.0.0/16",
		DisplayName:   "stateless-vcn",
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
		}{
			{
				DisplayName: "stateless-ingress",
				Protocol:    "6",
				Description: "Stateless inbound traffic",
				Source:      "0.0.0.0/0",
				Stateless:   true,
				TCPOptions: []struct {
					MinPort int `yaml:"min_port"`
					MaxPort int `yaml:"max_port"`
				}{
					{MinPort: 80, MaxPort: 80},
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
