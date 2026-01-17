package network

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"testing"
)

type mocks int

func TestCreateVCN(t *testing.T) {
	tests := []struct {
		name          string
		netCfg        NetCfg
		vcnName       string
		expectedError bool
	}{
		{
			name: "Valid VCN creation",
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
				}{
					{
						Type:        "ingress",
						Protocol:    "tcp",
						Description: "Allow SSH",
						Source:      "0.0.0.0/0",
						Stateless:   false,
						TCPOptions: []struct {
							MinPort int `yaml:"min_port"`
							MaxPort int `yaml:"max_port"`
						}{
							{MinPort: 22, MaxPort: 22},
						},
					},
				},
			},
			vcnName:       "my-vcn",
			expectedError: false,
		},
		{
			name: "VCN with different CIDR",
			netCfg: NetCfg{
				CompartmentID: "compartment-456",
				CidrBlock:     "192.168.0.0/24",
				DisplayName:   "another-vcn",
				Subnets: []struct {
					Name      string `yaml:"name"`
					CidrBlock string `yaml:"cidr_block"`
				}{},
				SecurityLists: []struct {
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
				}{},
			},
			vcnName:       "another-vcn",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pulumi.RunErr(func(ctx *pulumi.Context) error {
				vcn, err := tt.netCfg.CreateVCN(ctx, tt.vcnName)
				if err != nil {
					return err
				}

				if vcn == nil {
					t.Error("Expected VCN to be created, but got nil")
				}

				return nil
			}, pulumi.WithMocks("project", "stack", mocks(0)))

			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCreateVCNWithEmptyConfig(t *testing.T) {
	netCfg := NetCfg{
		CompartmentID: "",
		CidrBlock:     "",
		DisplayName:   "",
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		vcn, err := netCfg.CreateVCN(ctx, "test-vcn")
		if err != nil {
			return err
		}
		if vcn == nil {
			t.Error("Expected VCN to be created, but got nil")
		}
		return nil
	}, pulumi.WithMocks("project", "stack", mocks(0)))

	if err != nil {
		t.Errorf("Unexpected error with empty config: %v", err)
	}
}

func TestCreateVCNWithInvalidCidr(t *testing.T) {
	netCfg := NetCfg{
		CompartmentID: "compartment-123",
		CidrBlock:     "invalid-cidr",
		DisplayName:   "test-vcn",
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		vcn, err := netCfg.CreateVCN(ctx, "test-vcn")
		if err != nil {
			return err
		}
		if vcn == nil {
			t.Error("Expected VCN to be created, but got nil")
		}
		return nil
	}, pulumi.WithMocks("project", "stack", mocks(0)))

	if err != nil {
		t.Errorf("Unexpected error with invalid CIDR: %v", err)
	}
}
