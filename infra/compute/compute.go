package compute

import (
	"fmt"
	"infra/config"

	"github.com/pulumi/pulumi-oci/sdk/v3/go/oci/core"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ComputeCfg wraps the ComputeConfig and provides methods for managing compute resources
type ComputeCfg struct {
	config.ComputeConfig
}

// ValidateConfig validates the compute configuration
func (c *ComputeCfg) ValidateConfig() error {
	if c.CompartmentID == "" {
		return fmt.Errorf("compartment_id is required")
	}

	if len(c.Instances) == 0 {
		return fmt.Errorf("at least one instance must be defined")
	}

	for i, instance := range c.Instances {
		if instance.Name == "" {
			return fmt.Errorf("instance[%d]: name is required", i)
		}
		if instance.Shape == "" {
			return fmt.Errorf("instance[%d]: shape is required", i)
		}
		if instance.SubnetID == "" {
			return fmt.Errorf("instance[%d]: subnet_id is required", i)
		}
		if instance.ImageOCID == "" {
			return fmt.Errorf("instance[%d]: image_ocid is required", i)
		}
		if instance.SSHPublicKey == "" {
			return fmt.Errorf("instance[%d]: ssh_public_key is required", i)
		}
	}

	return nil
}

// CreateInstance creates a single compute instance
func (c *ComputeCfg) CreateInstance(ctx *pulumi.Context, instanceIndex int) (*core.Instance, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	if instanceIndex < 0 || instanceIndex >= len(c.Instances) {
		return nil, fmt.Errorf("instance index %d out of range", instanceIndex)
	}

	instance := c.Instances[instanceIndex]

	instanceArgs := &core.InstanceArgs{
		CompartmentId:      pulumi.String(c.CompartmentID),
		Shape:              pulumi.String(instance.Shape),
		AvailabilityDomain: pulumi.String("ad-1"),
		SourceDetails: &core.InstanceSourceDetailsArgs{
			SourceType: pulumi.String("image"),
			SourceId:   pulumi.String(instance.ImageOCID),
		},
		CreateVnicDetails: &core.InstanceCreateVnicDetailsArgs{
			SubnetId: pulumi.String(instance.SubnetID),
		},
		Metadata: pulumi.StringMap{
			"ssh_authorized_keys": pulumi.String(instance.SSHPublicKey),
		},
	}

	// Build ShapeConfig if OCPUCount or MemoryGB is specified
	if (instance.OCPUCount != nil && *instance.OCPUCount > 0) || (instance.MemoryGB != nil && *instance.MemoryGB > 0) {
		shapeConfig := &core.InstanceShapeConfigArgs{}
		if instance.OCPUCount != nil && *instance.OCPUCount > 0 {
			shapeConfig.Ocpus = pulumi.Float64(*instance.OCPUCount)
		}
		if instance.MemoryGB != nil && *instance.MemoryGB > 0 {
			shapeConfig.MemoryInGbs = pulumi.Float64(*instance.MemoryGB)
		}
		instanceArgs.ShapeConfig = shapeConfig
	}

	displayName := instance.DisplayName
	if displayName == "" {
		displayName = instance.Name
	}
	instanceArgs.DisplayName = pulumi.String(displayName)

	return core.NewInstance(ctx, instance.Name, instanceArgs)
}

// CreateAllInstances creates all instances defined in the configuration
func (c *ComputeCfg) CreateAllInstances(ctx *pulumi.Context) ([]*core.Instance, error) {
	if len(c.Instances) == 0 {
		return nil, fmt.Errorf("at least one instance must be defined")
	}

	var instances []*core.Instance

	for i := range c.Instances {
		instance, err := c.CreateInstance(ctx, i)
		if err != nil {
			return nil, fmt.Errorf("failed to create instance %s: %w", c.Instances[i].Name, err)
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

// GetInstance retrieves an instance config by name
func (c *ComputeCfg) GetInstance(name string) (*config.InstanceConfig, error) {
	for i := range c.Instances {
		if c.Instances[i].Name == name {
			return &c.Instances[i], nil
		}
	}
	return nil, fmt.Errorf("instance %s not found", name)
}

// CreateInstancesInSubnet creates all instances that belong to a specific subnet
func (c *ComputeCfg) CreateInstancesInSubnet(ctx *pulumi.Context, subnetID string) ([]*core.Instance, error) {
	if len(c.Instances) == 0 {
		return nil, fmt.Errorf("at least one instance must be defined")
	}

	var instances []*core.Instance

	for i, instance := range c.Instances {
		if instance.SubnetID == subnetID {
			createdInstance, err := c.CreateInstance(ctx, i)
			if err != nil {
				return nil, fmt.Errorf("failed to create instance %s in subnet %s: %w", instance.Name, subnetID, err)
			}
			instances = append(instances, createdInstance)
		}
	}

	return instances, nil
}
