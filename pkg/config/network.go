package config

import (
	"fmt"
	"os/exec"
	"strings"
)

// NetworkConfig représente la configuration réseau
type NetworkConfig struct {
	Version    int         `json:"version"`
	Config     []NetDevice `json:"config"`
	DNSServers []string    `json:"nameservers"`
	DNSDomain  string      `json:"domain"`
}

// NetDevice représente la configuration d'une interface réseau
type NetDevice struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	IPAddress string `json:"address"`
	Netmask   string `json:"netmask"`
	Gateway   string `json:"gateway"`
}

// ApplyNetworkConfig applique la configuration réseau sur Windows
func ApplyNetworkConfig(config *NetworkConfig) error {
	if len(config.Config) == 0 {
		return fmt.Errorf("aucune configuration réseau trouvée")
	}

	for _, device := range config.Config {
		// Configuration IP
		if err := setIPAddress(device); err != nil {
			return fmt.Errorf("erreur lors de la configuration IP: %v", err)
		}

		// Configuration Gateway
		if err := setGateway(device); err != nil {
			return fmt.Errorf("erreur lors de la configuration de la passerelle: %v", err)
		}
	}

	// Configuration DNS
	if err := setDNSServers(config.DNSServers, config.DNSDomain); err != nil {
		return fmt.Errorf("erreur lors de la configuration DNS: %v", err)
	}

	return nil
}

// setIPAddress configure l'adresse IP et le masque
func setIPAddress(device NetDevice) error {
	cmd := exec.Command("netsh", "interface", "ipv4", "set", "address",
		fmt.Sprintf("name=\"%s\"", device.Name),
		"source=static",
		fmt.Sprintf("addr=%s", device.IPAddress),
		fmt.Sprintf("mask=%s", device.Netmask))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erreur netsh: %v, output: %s", err, string(output))
	}

	fmt.Printf("✓ Adresse IP configurée: %s/%s\n", device.IPAddress, device.Netmask)
	return nil
}

// setGateway configure la passerelle par défaut
func setGateway(device NetDevice) error {
	if device.Gateway == "" {
		return nil // Pas de gateway à configurer
	}

	cmd := exec.Command("netsh", "interface", "ipv4", "set", "address",
		fmt.Sprintf("name=\"%s\"", device.Name),
		fmt.Sprintf("gateway=%s", device.Gateway),
		"gwmetric=1")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erreur netsh: %v, output: %s", err, string(output))
	}

	fmt.Printf("✓ Passerelle configurée: %s\n", device.Gateway)
	return nil
}

// setDNSServers configure les serveurs DNS
func setDNSServers(servers []string, domain string) error {
	if len(servers) == 0 {
		return nil // Pas de DNS à configurer
	}

	// Liste les interfaces réseau
	cmd := exec.Command("netsh", "interface", "show", "interface")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erreur lors de la liste des interfaces: %v", err)
	}

	// Trouve l'interface active
	lines := strings.Split(string(output), "\n")
	var interfaceName string
	for _, line := range lines {
		if strings.Contains(line, "Connected") {
			fields := strings.Fields(line)
			if len(fields) > 3 {
				interfaceName = strings.Join(fields[3:], " ")
				break
			}
		}
	}

	if interfaceName == "" {
		return fmt.Errorf("aucune interface réseau active trouvée")
	}

	// Configure les serveurs DNS
	for i, server := range servers {
		cmd = exec.Command("netsh", "interface", "ipv4", "set", "dns",
			fmt.Sprintf("name=\"%s\"", interfaceName),
			"source=static",
			fmt.Sprintf("addr=%s", server),
			fmt.Sprintf("register=primary"))

		if i > 0 {
			cmd = exec.Command("netsh", "interface", "ipv4", "add", "dns",
				fmt.Sprintf("name=\"%s\"", interfaceName),
				fmt.Sprintf("addr=%s", server),
				"index=2")
		}

		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("erreur lors de la configuration DNS: %v, output: %s", err, string(output))
		}
	}

	fmt.Printf("✓ Serveurs DNS configurés: %v\n", servers)
	if domain != "" {
		fmt.Printf("✓ Domaine DNS configuré: %s\n", domain)
	}

	return nil
}
