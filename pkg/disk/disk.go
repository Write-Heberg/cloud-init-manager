package disk

import (
	"fmt"
	"os"
	"path/filepath"
)

// DiskInfo represents information about a Cloud-Init disk
type DiskInfo struct {
	Path       string
	MountPoint string
	Files      []string
}

// DetectCloudInitDisk attempts to find the Cloud-Init disk on the system
func DetectCloudInitDisk() (*DiskInfo, error) {
	fmt.Println("\nDémarrage de la recherche du disque Cloud-Init...")
	fmt.Println("================================================")

	// Liste des lettres de lecteur possibles
	driveLetters := []string{"A:", "B:", "C:", "D:", "E:", "F:", "G:", "H:", "I:", "J:", "K:", "L:", "M:",
		"N:", "O:", "P:", "Q:", "R:", "S:", "T:", "U:", "V:", "W:", "X:", "Y:", "Z:"}

	foundDrives := 0
	for _, drive := range driveLetters {
		path := drive + "\\"

		// Vérifie si le lecteur existe et est accessible
		_, err := os.Stat(path)
		if err == nil {
			foundDrives++
			fmt.Printf("\nLecteur trouvé: %s\n", drive)

			// Vérifie le chemin OPENSTACK/LATEST
			openstackPath := filepath.Join(path, "OPENSTACK", "LATEST")
			fmt.Printf("  → Vérification du chemin %s\n", openstackPath)

			if isCloudInitPath(openstackPath) {
				fmt.Printf("\n✓ Configuration Cloud-Init trouvée dans %s\n", openstackPath)
				return &DiskInfo{
					Path:       path,
					MountPoint: openstackPath,
					Files:      []string{},
				}, nil
			}

			// Vérifie aussi le chemin OPENSTACK/CONTENT
			contentPath := filepath.Join(path, "OPENSTACK", "CONTENT")
			if _, err := os.Stat(contentPath); err == nil {
				fmt.Printf("  → Dossier CONTENT trouvé: %s\n", contentPath)
			}
		}
	}

	if foundDrives == 0 {
		return nil, fmt.Errorf("aucun lecteur accessible trouvé sur le système")
	}

	return nil, fmt.Errorf("\nAucune configuration Cloud-Init trouvée après vérification de %d lecteurs.\nChemin attendu: LECTEUR:\\OPENSTACK\\LATEST\nFichiers attendus:\n   - META_DATA.JSON\n   - USER_DATA\n   - VENDOR_DATA.JSON", foundDrives)
}

// isCloudInitPath vérifie si le chemin contient les fichiers typiques de cloud-init
func isCloudInitPath(path string) bool {
	// Liste des fichiers à vérifier
	requiredFiles := []string{
		"META_DATA.JSON",
		"USER_DATA",
		"VENDOR_DATA.JSON",
	}

	// Vérifie si le dossier existe
	if _, err := os.Stat(path); err != nil {
		return false
	}

	fmt.Printf("  → Vérification des fichiers dans %s\n", path)

	// Vérifie la présence de chaque fichier
	for _, file := range requiredFiles {
		fullPath := filepath.Join(path, file)
		if _, err := os.Stat(fullPath); err == nil {
			fmt.Printf("  → Fichier trouvé: %s\n", file)
		} else {
			fmt.Printf("  → Fichier manquant: %s\n", file)
			// On continue la vérification même si un fichier est manquant
		}
	}

	// Vérifie qu'au moins un des fichiers existe
	for _, file := range requiredFiles {
		if _, err := os.Stat(filepath.Join(path, file)); err == nil {
			return true
		}
	}

	return false
}

// ListCloudInitFiles returns a list of files in the Cloud-Init disk
func (d *DiskInfo) ListCloudInitFiles() error {
	fmt.Printf("\nLecture du contenu du dossier %s...\n", d.MountPoint)

	files, err := filepath.Glob(filepath.Join(d.MountPoint, "*"))
	if err != nil {
		return fmt.Errorf("impossible de lister les fichiers: %v", err)
	}

	d.Files = files
	return nil
}

// ReadCloudInitFile reads the content of a specific file from the Cloud-Init disk
func (d *DiskInfo) ReadCloudInitFile(filename string) ([]byte, error) {
	fullPath := filepath.Join(d.MountPoint, filename)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("impossible de lire le fichier %s: %v", filename, err)
	}
	return content, nil
}
