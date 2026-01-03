package network

import (
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"testing"
)

type mocks int

func (mocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return args.Name + "_id", args.Inputs, nil
}

func (mocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return args.Args, nil
}

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
