package main

import (
	"fmt"
	"log"

	"cloud-init-manager/pkg/disk"
	"cloud-init-manager/pkg/parser"
)

func main() {
	log.Println("Starting Cloud-Init Manager...")

	// Detect Cloud-Init disk
	diskInfo, err := disk.DetectCloudInitDisk()
	if err != nil {
		log.Fatalf("Failed to detect Cloud-Init disk: %v", err)
	}

	// List Cloud-Init files
	err = diskInfo.ListCloudInitFiles()
	if err != nil {
		log.Fatalf("Failed to list Cloud-Init files: %v", err)
	}

	log.Printf("Found %d files on Cloud-Init disk", len(diskInfo.Files))

	// Read and parse each file
	for _, file := range diskInfo.Files {
		content, err := diskInfo.ReadCloudInitFile(file)
		if err != nil {
			log.Printf("Warning: Failed to read file %s: %v", file, err)
			continue
		}

		config, err := parser.ParseYAML(content)
		if err != nil {
			log.Printf("Warning: Failed to parse file %s: %v", file, err)
			continue
		}

		// Export configuration as JSON
		jsonData, err := config.ExportJSON()
		if err != nil {
			log.Printf("Warning: Failed to export JSON for file %s: %v", file, err)
			continue
		}

		fmt.Printf("Configuration from %s:\n%s\n", file, string(jsonData))
	}

	fmt.Println("Cloud-Init Manager completed successfully")
}
