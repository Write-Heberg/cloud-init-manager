package main

import (
	"encoding/json"
	"fmt"
	"os"

	"cloud-init-manager/pkg/config"
	"cloud-init-manager/pkg/disk"
)

func main() {
	fmt.Println("=== Cloud-Init Manager v0.1 ===")
	fmt.Println("Recherche et configuration du système...")

	// Detect Cloud-Init disk
	diskInfo, err := disk.DetectCloudInitDisk()
	if err != nil {
		fmt.Printf("Erreur: %v\n", err)
		waitForEnter()
		return
	}

	fmt.Printf("\nConfiguration trouvée dans: %s\n", diskInfo.MountPoint)
	fmt.Println("\nAnalyse des fichiers de configuration...")

	// Lecture des configurations
	var networkConfig config.NetworkConfig
	var userConfig config.UserConfig

	// Parcours les fichiers trouvés
	for _, file := range diskInfo.Files {
		content, err := diskInfo.ReadCloudInitFile(file)
		if err != nil {
			fmt.Printf("Erreur de lecture %s: %v\n", file, err)
			continue
		}

		// Parse selon le type de fichier
		switch {
		case contains(file, "META_DATA.JSON"):
			fmt.Println("\n=== Configuration Metadata ===")
			var metadata map[string]interface{}
			if err := json.Unmarshal(content, &metadata); err != nil {
				fmt.Printf("Erreur parsing metadata: %v\n", err)
				continue
			}
			fmt.Printf("✓ Metadata lu avec succès\n")

		case contains(file, "USER_DATA"):
			fmt.Println("\n=== Configuration Utilisateur ===")
			if err := json.Unmarshal(content, &userConfig); err != nil {
				fmt.Printf("Erreur parsing user-data: %v\n", err)
				continue
			}
			fmt.Printf("✓ Configuration utilisateur trouvée pour: %s\n", userConfig.Name)

		case contains(file, "VENDOR_DATA.JSON"):
			fmt.Println("\n=== Configuration Réseau ===")
			if err := json.Unmarshal(content, &networkConfig); err != nil {
				fmt.Printf("Erreur parsing network config: %v\n", err)
				continue
			}
			fmt.Printf("✓ Configuration réseau trouvée\n")
		}
	}

	// Application des configurations
	fmt.Println("\nApplication des configurations...")

	// Configure le réseau si des paramètres sont trouvés
	if networkConfig.Version > 0 {
		fmt.Println("\n=== Application de la configuration réseau ===")
		if err := config.ApplyNetworkConfig(&networkConfig); err != nil {
			fmt.Printf("Erreur configuration réseau: %v\n", err)
		}
	}

	// Configure l'utilisateur si des paramètres sont trouvés
	if userConfig.Name != "" {
		fmt.Println("\n=== Application de la configuration utilisateur ===")
		if err := config.ApplyUserConfig(&userConfig); err != nil {
			fmt.Printf("Erreur configuration utilisateur: %v\n", err)
		}
	}

	fmt.Println("\n=== Configuration terminée ===")
	waitForEnter()
}

func waitForEnter() {
	fmt.Println("\nAppuyez sur Entrée pour quitter...")
	os.Stdin.Read(make([]byte, 1))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}
