package main

import (
	"infra/compute"
	"infra/config"
	"infra/network"
	"log"

	"github.com/pulumi/pulumi-oci/sdk/go/oci/identity"
	"github.com/pulumi/pulumi-oci/sdk/go/oci/objectstorage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// TODO properly orchestrate resource creation
func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		var cfg config.Config
		err := cfg.LoadFromYaml("config/dev/config.yaml")
		if err != nil {
			return err
		}

		ncfg := network.NetCfg{NetworkConfig: cfg.Network}
		vcn, err := ncfg.CreateVCN(ctx, ncfg.DisplayName)

		if err != nil {
			log.Printf("Failed to create VCN with error: %v", err)
			return err
		}

		// Create security lists and get a map for easy reference
		securityLists, err := ncfg.CreateACLMap(ctx, vcn.ID().ElementType().String())
		if err != nil {
			log.Printf("Failed to create security lists with error: %v", err)
			return err
		}

		// Create subnets with their respective security lists attached based on subnet_name configuration
		subnets, err := ncfg.CreateAllSubnetsWithSecurityLists(ctx, vcn.ID().ElementType().String(), securityLists)
		if err != nil {
			log.Printf("Failed to create subnets with security lists with error: %v", err)
			return err
		}

		// Export subnet IDs for reference
		for i, subnet := range subnets {
			ctx.Export("subnet-"+string(rune(i)), subnet.ID())
		}

		ccfg := compute.ComputeCfg{ComputeConfig: cfg.Compute}
		instances, err := ccfg.CreateAllInstances(ctx)
		if err != nil {
			log.Printf("Failed to create compute instances with error: %v", err)
			return err
		}
		for i, instance := range instances {
			ctx.Export("instance-"+string(rune(i)), instance.ID())
		}

		myCompartment, err := identity.NewCompartment(ctx, "my-compartment", &identity.CompartmentArgs{
			Name:         pulumi.String("my-compartment"),
			Description:  pulumi.String("My description text"),
			EnableDelete: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		myNamespace := pulumi.All(myCompartment.CompartmentId).ApplyT(
			func(args []interface{}) (string, error) {
				namespace, err := objectstorage.GetNamespace(ctx, &objectstorage.GetNamespaceArgs{
					CompartmentId: pulumi.StringRef(args[0].(string)),
				})
				if err != nil {
					return "", err
				}
				return namespace.Namespace, nil
			},
		).(pulumi.StringOutput)

		myBucket, err := objectstorage.NewBucket(ctx, "my-bucket", &objectstorage.BucketArgs{
			Name:          pulumi.String("my-bucket"),
			Namespace:     myNamespace,
			CompartmentId: myCompartment.ID(),
		})
		if err != nil {
			return err
		}

		ctx.Export("name", myBucket.Name)

		return nil
	})
}
