package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type BaseConfig struct {
	CompartmentID string `yaml:"compartment_id"`
	Region        string `yaml:"region,omitempty"`
}

type SubnetConfig struct {
	Name      string `yaml:"name"`
	CidrBlock string `yaml:"cidr_block"`
}

type TCPOptionConfig struct {
	MinPort int `yaml:"min_port"`
	MaxPort int `yaml:"max_port"`
}

type SecurityListConfig struct {
	DisplayName string            `yaml:"display_name"`
	Protocol    string            `yaml:"protocol"`
	Description string            `yaml:"description"`
	Destination string            `yaml:"destination"`
	Source      string            `yaml:"source"`
	Stateless   bool              `yaml:"stateless"`
	TCPOptions  []TCPOptionConfig `yaml:"tcp_options"`
}

type NetworkConfig struct {
	BaseConfig
	CidrBlock     string               `yaml:"cidr_block"`
	DisplayName   string               `yaml:"display_name"`
	Subnets       []SubnetConfig       `yaml:"subnets"`
	SecurityLists []SecurityListConfig `yaml:"security_lists"`
}

type ComputeConfig struct {
	BaseConfig
	InstanceShape string `yaml:"instance_shape"`
}

type BastionConfig struct {
	BaseConfig
}

type HeatwaveConfig struct {
	BaseConfig
}

type Config struct {
	Network  NetworkConfig  `yaml:"network"`
	Compute  ComputeConfig  `yaml:"compute"`
	Bastion  BastionConfig  `yaml:"bastion"`
	Heatwave HeatwaveConfig `yaml:"heatwave"`
}

func (c *Config) LoadFromYaml(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, c)
}
