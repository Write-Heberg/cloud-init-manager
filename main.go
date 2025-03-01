package main

import (
	"bufio"
	"fmt"
	"os"

	"cloud-init-manager/pkg/disk"
	"cloud-init-manager/pkg/parser"
)

func main() {
	fmt.Println("=== Cloud-Init Manager v0.1 ===")
	fmt.Println("Recherche du disque Cloud-Init...")

	// Detect Cloud-Init disk
	diskInfo, err := disk.DetectCloudInitDisk()
	if err != nil {
		fmt.Printf("Erreur: %v\n", err)
		waitForEnter()
		return
	}

	fmt.Printf("Disque Cloud-Init trouvé: %s\n", diskInfo.Path)

	// List Cloud-Init files
	err = diskInfo.ListCloudInitFiles()
	if err != nil {
		fmt.Printf("Erreur lors de la lecture des fichiers: %v\n", err)
		waitForEnter()
		return
	}

	fmt.Printf("\nFichiers trouvés (%d):\n", len(diskInfo.Files))
	for _, file := range diskInfo.Files {
		fmt.Printf("- %s\n", file)
	}

	// Read and parse each file
	fmt.Println("\nLecture des fichiers de configuration...")
	for _, file := range diskInfo.Files {
		fmt.Printf("\nAnalyse du fichier: %s\n", file)

		content, err := diskInfo.ReadCloudInitFile(file)
		if err != nil {
			fmt.Printf("  Erreur de lecture: %v\n", err)
			continue
		}

		config, err := parser.ParseYAML(content)
		if err != nil {
			fmt.Printf("  Erreur de parsing YAML: %v\n", err)
			continue
		}

		// Export configuration as JSON
		jsonData, err := config.ExportJSON()
		if err != nil {
			fmt.Printf("  Erreur d'export JSON: %v\n", err)
			continue
		}

		fmt.Printf("Configuration trouvée:\n%s\n", string(jsonData))
	}

	fmt.Println("\nAnalyse terminée.")
	waitForEnter()
}

func waitForEnter() {
	fmt.Println("\nAppuyez sur Entrée pour quitter...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
