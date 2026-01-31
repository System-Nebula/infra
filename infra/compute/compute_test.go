package compute

import (
	"infra/config"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"testing"
)

type ComputeMocks int

func (ComputeMocks) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	return args.Name + "_id", args.Inputs, nil
}

func (ComputeMocks) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return args.Args, nil
}

// Test helpers and constants

const (
	testCompartmentID = "compartment-123"
	testSubnetID      = "subnet-123"
	testImageOCID     = "ocid1.image.oc1..example"
	testSSHPublicKey  = "ssh-rsa AAAAB3NzaC1yc2E..."
)

// newTestComputeCfg creates a test ComputeCfg with the specified instances
func newTestComputeCfg(compartmentID string, instances []config.InstanceConfig) ComputeCfg {
	return ComputeCfg{
		ComputeConfig: config.ComputeConfig{
			BaseConfig: config.BaseConfig{
				CompartmentID: compartmentID,
			},
			Instances: instances,
		},
	}
}

// newTestInstance creates a test InstanceConfig with common values
func newTestInstance(name string) config.InstanceConfig {
	return config.InstanceConfig{
		Name:         name,
		Shape:        "VM.Standard.E4.Flex",
		SubnetID:     testSubnetID,
		ImageOCID:    testImageOCID,
		SSHPublicKey: testSSHPublicKey,
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name          string
		computeCfg    ComputeCfg
		expectedError bool
	}{
		{
			name: "Valid configuration",
			computeCfg: ComputeCfg{
				ComputeConfig: config.ComputeConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-123",
					},
					Instances: []config.InstanceConfig{
						{
							Name:         "instance-1",
							Shape:        "VM.Standard.E4.Flex",
							SubnetID:     "subnet-123",
							ImageOCID:    "ocid1.image.oc1..example",
							SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name: "Empty compartment ID",
			computeCfg: ComputeCfg{
				ComputeConfig: config.ComputeConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "",
					},
					Instances: []config.InstanceConfig{
						{
							Name:         "instance-1",
							Shape:        "VM.Standard.E4.Flex",
							SubnetID:     "subnet-123",
							ImageOCID:    "ocid1.image.oc1..example",
							SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
						},
					},
				},
			},
			expectedError: true,
		},
		{
			name: "No instances defined",
			computeCfg: ComputeCfg{
				ComputeConfig: config.ComputeConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-123",
					},
					Instances: []config.InstanceConfig{},
				},
			},
			expectedError: true,
		},
		{
			name: "Instance missing name",
			computeCfg: ComputeCfg{
				ComputeConfig: config.ComputeConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-123",
					},
					Instances: []config.InstanceConfig{
						{
							Shape:        "VM.Standard.E4.Flex",
							SubnetID:     "subnet-123",
							ImageOCID:    "ocid1.image.oc1..example",
							SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
						},
					},
				},
			},
			expectedError: true,
		},
		{
			name: "Instance missing shape",
			computeCfg: ComputeCfg{
				ComputeConfig: config.ComputeConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-123",
					},
					Instances: []config.InstanceConfig{
						{
							Name:         "instance-1",
							SubnetID:     "subnet-123",
							ImageOCID:    "ocid1.image.oc1..example",
							SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
						},
					},
				},
			},
			expectedError: true,
		},
		{
			name: "Instance missing subnet ID",
			computeCfg: ComputeCfg{
				ComputeConfig: config.ComputeConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-123",
					},
					Instances: []config.InstanceConfig{
						{
							Name:         "instance-1",
							Shape:        "VM.Standard.E4.Flex",
							ImageOCID:    "ocid1.image.oc1..example",
							SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
						},
					},
				},
			},
			expectedError: true,
		},
		{
			name: "Instance missing image OCID",
			computeCfg: ComputeCfg{
				ComputeConfig: config.ComputeConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-123",
					},
					Instances: []config.InstanceConfig{
						{
							Name:         "instance-1",
							Shape:        "VM.Standard.E4.Flex",
							SubnetID:     "subnet-123",
							SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
						},
					},
				},
			},
			expectedError: true,
		},
		{
			name: "Instance missing SSH public key",
			computeCfg: ComputeCfg{
				ComputeConfig: config.ComputeConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-123",
					},
					Instances: []config.InstanceConfig{
						{
							Name:      "instance-1",
							Shape:     "VM.Standard.E4.Flex",
							SubnetID:  "subnet-123",
							ImageOCID: "ocid1.image.oc1..example",
						},
					},
				},
			},
			expectedError: true,
		},
		{
			name: "Multiple valid instances",
			computeCfg: ComputeCfg{
				ComputeConfig: config.ComputeConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-123",
					},
					Instances: []config.InstanceConfig{
						{
							Name:         "instance-1",
							Shape:        "VM.Standard.E4.Flex",
							SubnetID:     "subnet-123",
							ImageOCID:    "ocid1.image.oc1..example",
							SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
						},
						{
							Name:         "instance-2",
							Shape:        "VM.Standard.E3.Flex",
							SubnetID:     "subnet-456",
							ImageOCID:    "ocid1.image.oc1..example",
							SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
						},
					},
				},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.computeCfg.ValidateConfig()
			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCreateInstance(t *testing.T) {
	tests := []struct {
		name          string
		computeCfg    ComputeCfg
		instanceIndex int
		expectedError bool
	}{
		{
			name: "Create instance with index 0",
			computeCfg: ComputeCfg{
				ComputeConfig: config.ComputeConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-123",
					},
					Instances: []config.InstanceConfig{
						{
							Name:         "test-instance-1",
							DisplayName:  "Test Instance 1",
							Shape:        "VM.Standard.E4.Flex",
							SubnetID:     "subnet-123",
							ImageOCID:    "ocid1.image.oc1..example",
							SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
						},
						{
							Name:         "test-instance-2",
							DisplayName:  "Test Instance 2",
							Shape:        "VM.Standard.E3.Flex",
							SubnetID:     "subnet-456",
							ImageOCID:    "ocid1.image.oc1..example",
							SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
						},
					},
				},
			},
			instanceIndex: 0,
			expectedError: false,
		},
		{
			name: "Create instance with index 1",
			computeCfg: ComputeCfg{
				ComputeConfig: config.ComputeConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-123",
					},
					Instances: []config.InstanceConfig{
						{
							Name:         "test-instance-1",
							Shape:        "VM.Standard.E4.Flex",
							SubnetID:     "subnet-123",
							ImageOCID:    "ocid1.image.oc1..example",
							SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
						},
						{
							Name:         "test-instance-2",
							Shape:        "VM.Standard.E3.Flex",
							SubnetID:     "subnet-456",
							ImageOCID:    "ocid1.image.oc1..example",
							SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
						},
					},
				},
			},
			instanceIndex: 1,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pulumi.RunErr(func(ctx *pulumi.Context) error {
				instance, err := tt.computeCfg.CreateInstance(ctx, tt.instanceIndex)
				if err != nil {
					return err
				}

				if instance == nil {
					t.Error("Expected instance to be created, but got nil")
				}

				return nil
			}, pulumi.WithMocks("project", "stack", ComputeMocks(0)))

			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCreateInstanceWithInvalidIndex(t *testing.T) {
	computeCfg := ComputeCfg{
		ComputeConfig: config.ComputeConfig{
			BaseConfig: config.BaseConfig{
				CompartmentID: "compartment-123",
			},
			Instances: []config.InstanceConfig{
				{
					Name:         "test-instance",
					Shape:        "VM.Standard.E4.Flex",
					SubnetID:     "subnet-123",
					ImageOCID:    "ocid1.image.oc1..example",
					SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
				},
			},
		},
	}

	tests := []struct {
		name          string
		instanceIndex int
	}{
		{"Negative index", -1},
		{"Index out of bounds", 1},
		{"Index way out of bounds", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := computeCfg.CreateInstance(nil, tt.instanceIndex)
			if err == nil {
				t.Errorf("Expected error with invalid index %d but got none", tt.instanceIndex)
			}
		})
	}
}

func TestCreateInstanceWithShapeConfig(t *testing.T) {
	tests := []struct {
		name          string
		computeCfg    ComputeCfg
		instanceIndex int
		expectedError bool
	}{
		{
			name: "Instance with OCPUCount only",
			computeCfg: ComputeCfg{
				ComputeConfig: config.ComputeConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-123",
					},
					Instances: []config.InstanceConfig{
						{
							Name:         "instance-with-ocpu",
							Shape:        "VM.Standard.E4.Flex",
							SubnetID:     "subnet-123",
							ImageOCID:    "ocid1.image.oc1..example",
							SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
							OCPUCount:    func() *float64 { v := 2.0; return &v }(),
						},
					},
				},
			},
			instanceIndex: 0,
			expectedError: false,
		},
		{
			name: "Instance with MemoryGB only",
			computeCfg: ComputeCfg{
				ComputeConfig: config.ComputeConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-123",
					},
					Instances: []config.InstanceConfig{
						{
							Name:         "instance-with-memory",
							Shape:        "VM.Standard.E4.Flex",
							SubnetID:     "subnet-123",
							ImageOCID:    "ocid1.image.oc1..example",
							SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
							MemoryGB:     func() *float64 { v := 16.0; return &v }(),
						},
					},
				},
			},
			instanceIndex: 0,
			expectedError: false,
		},
		{
			name: "Instance with both OCPUCount and MemoryGB",
			computeCfg: ComputeCfg{
				ComputeConfig: config.ComputeConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-123",
					},
					Instances: []config.InstanceConfig{
						{
							Name:         "instance-with-both",
							Shape:        "VM.Standard.E4.Flex",
							SubnetID:     "subnet-123",
							ImageOCID:    "ocid1.image.oc1..example",
							SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
							OCPUCount:    func() *float64 { v := 4.0; return &v }(),
							MemoryGB:     func() *float64 { v := 32.0; return &v }(),
						},
					},
				},
			},
			instanceIndex: 0,
			expectedError: false,
		},
		{
			name: "Instance without OCPUCount or MemoryGB",
			computeCfg: ComputeCfg{
				ComputeConfig: config.ComputeConfig{
					BaseConfig: config.BaseConfig{
						CompartmentID: "compartment-123",
					},
					Instances: []config.InstanceConfig{
						{
							Name:         "instance-without-shape-config",
							Shape:        "VM.Standard.E4.Flex",
							SubnetID:     "subnet-123",
							ImageOCID:    "ocid1.image.oc1..example",
							SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
						},
					},
				},
			},
			instanceIndex: 0,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pulumi.RunErr(func(ctx *pulumi.Context) error {
				instance, err := tt.computeCfg.CreateInstance(ctx, tt.instanceIndex)
				if err != nil {
					return err
				}

				if instance == nil {
					t.Error("Expected instance to be created, but got nil")
				}

				return nil
			}, pulumi.WithMocks("project", "stack", ComputeMocks(0)))

			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCreateAllInstances(t *testing.T) {
	computeCfg := ComputeCfg{
		ComputeConfig: config.ComputeConfig{
			BaseConfig: config.BaseConfig{
				CompartmentID: "compartment-123",
			},
			Instances: []config.InstanceConfig{
				{
					Name:         "instance-1",
					Shape:        "VM.Standard.E4.Flex",
					SubnetID:     "subnet-123",
					ImageOCID:    "ocid1.image.oc1..example",
					SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
				},
				{
					Name:         "instance-2",
					Shape:        "VM.Standard.E3.Flex",
					SubnetID:     "subnet-456",
					ImageOCID:    "ocid1.image.oc1..example",
					SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
				},
				{
					Name:         "instance-3",
					Shape:        "VM.Standard.E2.Flex",
					SubnetID:     "subnet-789",
					ImageOCID:    "ocid1.image.oc1..example",
					SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
				},
			},
		},
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		instances, err := computeCfg.CreateAllInstances(ctx)
		if err != nil {
			return err
		}

		if len(instances) != 3 {
			t.Errorf("Expected 3 instances, but got %d", len(instances))
		}

		for i, instance := range instances {
			if instance == nil {
				t.Errorf("Expected instance %d to be created, but got nil", i)
			}
		}

		return nil
	}, pulumi.WithMocks("project", "stack", ComputeMocks(0)))

	if err != nil {
		t.Errorf("Unexpected error creating all instances: %v", err)
	}
}

func TestCreateAllInstancesWithEmptySlice(t *testing.T) {
	computeCfg := ComputeCfg{
		ComputeConfig: config.ComputeConfig{
			BaseConfig: config.BaseConfig{
				CompartmentID: "compartment-123",
			},
			Instances: []config.InstanceConfig{},
		},
	}

	_, err := computeCfg.CreateAllInstances(nil)
	if err == nil {
		t.Errorf("Expected error with empty instances slice but got none")
	}
}

func TestCreateInstanceWithNilContext(t *testing.T) {
	computeCfg := ComputeCfg{
		ComputeConfig: config.ComputeConfig{
			BaseConfig: config.BaseConfig{
				CompartmentID: "compartment-123",
			},
			Instances: []config.InstanceConfig{
				{
					Name:         "test-instance",
					Shape:        "VM.Standard.E4.Flex",
					SubnetID:     "subnet-123",
					ImageOCID:    "ocid1.image.oc1..example",
					SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
				},
			},
		},
	}

	_, err := computeCfg.CreateInstance(nil, 0)
	if err == nil {
		t.Error("Expected error with nil context but got none")
	}
}

func TestGetInstance(t *testing.T) {
	computeCfg := ComputeCfg{
		ComputeConfig: config.ComputeConfig{
			BaseConfig: config.BaseConfig{
				CompartmentID: "compartment-123",
			},
			Instances: []config.InstanceConfig{
				{
					Name:         "instance-1",
					Shape:        "VM.Standard.E4.Flex",
					SubnetID:     "subnet-123",
					ImageOCID:    "ocid1.image.oc1..example",
					SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
				},
				{
					Name:         "instance-2",
					Shape:        "VM.Standard.E3.Flex",
					SubnetID:     "subnet-456",
					ImageOCID:    "ocid1.image.oc1..example",
					SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
				},
			},
		},
	}

	tests := []struct {
		name          string
		instanceName  string
		expectedError bool
	}{
		{
			name:          "Get existing instance",
			instanceName:  "instance-1",
			expectedError: false,
		},
		{
			name:          "Get another existing instance",
			instanceName:  "instance-2",
			expectedError: false,
		},
		{
			name:          "Get non-existent instance",
			instanceName:  "instance-3",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance, err := computeCfg.GetInstance(tt.instanceName)
			if tt.expectedError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectedError && instance.Name != tt.instanceName {
				t.Errorf("Expected instance name %s, but got %s", tt.instanceName, instance.Name)
			}
		})
	}
}

func TestCreateInstancesInSubnet(t *testing.T) {
	computeCfg := ComputeCfg{
		ComputeConfig: config.ComputeConfig{
			BaseConfig: config.BaseConfig{
				CompartmentID: "compartment-123",
			},
			Instances: []config.InstanceConfig{
				{
					Name:         "instance-1",
					Shape:        "VM.Standard.E4.Flex",
					SubnetID:     "subnet-123",
					ImageOCID:    "ocid1.image.oc1..example",
					SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
				},
				{
					Name:         "instance-2",
					Shape:        "VM.Standard.E3.Flex",
					SubnetID:     "subnet-123",
					ImageOCID:    "ocid1.image.oc1..example",
					SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
				},
				{
					Name:         "instance-3",
					Shape:        "VM.Standard.E2.Flex",
					SubnetID:     "subnet-456",
					ImageOCID:    "ocid1.image.oc1..example",
					SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
				},
			},
		},
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		instances, err := computeCfg.CreateInstancesInSubnet(ctx, "subnet-123")
		if err != nil {
			return err
		}

		if len(instances) != 2 {
			t.Errorf("Expected 2 instances in subnet-123, but got %d", len(instances))
		}

		for i, instance := range instances {
			if instance == nil {
				t.Errorf("Expected instance %d to be created, but got nil", i)
			}
		}

		return nil
	}, pulumi.WithMocks("project", "stack", ComputeMocks(0)))

	if err != nil {
		t.Errorf("Unexpected error creating instances in subnet: %v", err)
	}
}

func TestCreateInstancesInSubnetNoMatches(t *testing.T) {
	computeCfg := ComputeCfg{
		ComputeConfig: config.ComputeConfig{
			BaseConfig: config.BaseConfig{
				CompartmentID: "compartment-123",
			},
			Instances: []config.InstanceConfig{
				{
					Name:         "instance-1",
					Shape:        "VM.Standard.E4.Flex",
					SubnetID:     "subnet-123",
					ImageOCID:    "ocid1.image.oc1..example",
					SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
				},
				{
					Name:         "instance-2",
					Shape:        "VM.Standard.E3.Flex",
					SubnetID:     "subnet-456",
					ImageOCID:    "ocid1.image.oc1..example",
					SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
				},
			},
		},
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		instances, err := computeCfg.CreateInstancesInSubnet(ctx, "subnet-789")
		if err != nil {
			return err
		}

		if len(instances) != 0 {
			t.Errorf("Expected 0 instances in subnet-789, but got %d", len(instances))
		}

		return nil
	}, pulumi.WithMocks("project", "stack", ComputeMocks(0)))

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// Benchmarks

func BenchmarkValidateConfig(b *testing.B) {
	benchmarks := []struct {
		name          string
		instanceCount int
	}{
		{"Single instance", 1},
		{"10 instances", 10},
		{"100 instances", 100},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			instances := make([]config.InstanceConfig, bm.instanceCount)
			for i := 0; i < bm.instanceCount; i++ {
				instances[i] = newTestInstance(
					"instance-" + string(rune(i)),
				)
			}
			cfg := newTestComputeCfg(testCompartmentID, instances)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = cfg.ValidateConfig()
			}
		})
	}
}

func BenchmarkCreateAllInstances(b *testing.B) {
	benchmarks := []struct {
		name          string
		instanceCount int
	}{
		{"Single instance", 1},
		{"10 instances", 10},
		{"100 instances", 100},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			instances := make([]config.InstanceConfig, bm.instanceCount)
			for i := 0; i < bm.instanceCount; i++ {
				instances[i] = newTestInstance(
					"instance-" + string(rune(i)),
				)
			}
			cfg := newTestComputeCfg(testCompartmentID, instances)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = pulumi.RunErr(func(ctx *pulumi.Context) error {
					_, err := cfg.CreateAllInstances(ctx)
					return err
				}, pulumi.WithMocks("project", "stack", ComputeMocks(0)))
			}
		})
	}
}

func BenchmarkCreateInstance(b *testing.B) {
	cfg := newTestComputeCfg(testCompartmentID, []config.InstanceConfig{
		newTestInstance("benchmark-instance"),
	})

	for b.Loop() {
		_ = pulumi.RunErr(func(ctx *pulumi.Context) error {
			_, err := cfg.CreateInstance(ctx, 0)
			return err
		}, pulumi.WithMocks("project", "stack", ComputeMocks(0)))
	}
}

func TestCreateInstancesInSubnetWithEmptySlice(t *testing.T) {
	computeCfg := ComputeCfg{
		ComputeConfig: config.ComputeConfig{
			BaseConfig: config.BaseConfig{
				CompartmentID: "compartment-123",
			},
			Instances: []config.InstanceConfig{},
		},
	}

	_, err := computeCfg.CreateInstancesInSubnet(nil, "subnet-123")
	if err == nil {
		t.Errorf("Expected error with empty instances slice but got none")
	}
}

func TestCreateAllInstancesWithDuplicateNames(t *testing.T) {
	computeCfg := ComputeCfg{
		ComputeConfig: config.ComputeConfig{
			BaseConfig: config.BaseConfig{
				CompartmentID: "compartment-123",
			},
			Instances: []config.InstanceConfig{
				{
					Name:         "duplicate-instance",
					Shape:        "VM.Standard.E4.Flex",
					SubnetID:     "subnet-123",
					ImageOCID:    "ocid1.image.oc1..example",
					SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
				},
				{
					Name:         "duplicate-instance",
					Shape:        "VM.Standard.E3.Flex",
					SubnetID:     "subnet-456",
					ImageOCID:    "ocid1.image.oc1..example",
					SSHPublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
				},
			},
		},
	}

	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		instances, err := computeCfg.CreateAllInstances(ctx)
		if err != nil {
			return err
		}

		if len(instances) != 2 {
			t.Errorf("Expected 2 instances to be created with duplicate names, but got %d", len(instances))
		}

		return nil
	}, pulumi.WithMocks("project", "stack", ComputeMocks(0)))

	if err != nil {
		t.Errorf("Unexpected error creating instances with duplicate names: %v", err)
	}
}
