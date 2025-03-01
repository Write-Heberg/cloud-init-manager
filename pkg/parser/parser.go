package parser

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

// CloudConfig represents the structure of cloud-init configuration
type CloudConfig struct {
	Hostname    string   `yaml:"hostname" json:"hostname"`
	Users       []User   `yaml:"users" json:"users"`
	Network     Network  `yaml:"network" json:"network"`
	SSHKeys     []string `yaml:"ssh_authorized_keys" json:"ssh_authorized_keys"`
	RDPSettings RDP      `yaml:"rdp" json:"rdp"`
}

type User struct {
	Name              string   `yaml:"name" json:"name"`
	Password          string   `yaml:"passwd" json:"passwd"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys" json:"ssh_authorized_keys"`
	Sudo              string   `yaml:"sudo" json:"sudo"`
}

type Network struct {
	Version   int         `yaml:"version" json:"version"`
	Config    []NetConfig `yaml:"config" json:"config"`
	DNSServer []string    `yaml:"nameservers" json:"nameservers"`
	DNSDomain string      `yaml:"domain" json:"domain"`
}

type NetConfig struct {
	Type      string `yaml:"type" json:"type"`
	Name      string `yaml:"name" json:"name"`
	IPAddress string `yaml:"address" json:"address"`
	Netmask   string `yaml:"netmask" json:"netmask"`
	Gateway   string `yaml:"gateway" json:"gateway"`
}

type RDP struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	Port     int    `yaml:"port" json:"port"`
	Security string `yaml:"security" json:"security"`
}

// ParseYAML parses YAML content into CloudConfig structure
func ParseYAML(content []byte) (*CloudConfig, error) {
	var config CloudConfig
	err := yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %v", err)
	}
	return &config, nil
}

// ExportJSON converts CloudConfig to JSON format
func (c *CloudConfig) ExportJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}
