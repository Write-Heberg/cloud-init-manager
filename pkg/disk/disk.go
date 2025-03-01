package disk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
			fmt.Printf("  → Vérification des fichiers cloud-init...\n")

			// Vérifie si c'est notre disque cloud-init
			if isCloudInitDisk(path) {
				fmt.Printf("\n✓ Disque Cloud-Init confirmé sur %s\n", drive)
				return &DiskInfo{
					Path:       path,
					MountPoint: path,
					Files:      []string{},
				}, nil
			} else {
				fmt.Printf("  → Pas de fichiers cloud-init sur ce lecteur\n")
			}
		}
	}

	if foundDrives == 0 {
		return nil, fmt.Errorf("aucun lecteur accessible trouvé sur le système")
	}

	return nil, fmt.Errorf("\nAucun disque Cloud-Init trouvé après vérification de %d lecteurs.\nAssurez-vous que:\n1. Le disque cloud-init est bien monté\n2. Il contient au moins un des fichiers suivants:\n   - meta-data\n   - user-data\n   - network-config\n3. Vous avez les droits administrateur", foundDrives)
}

// isCloudInitDisk vérifie si le lecteur contient les fichiers typiques de cloud-init
func isCloudInitDisk(path string) bool {
	// Liste des fichiers possibles de cloud-init
	commonFiles := []string{
		"meta-data",
		"user-data",
		"network-config",
		"vendor-data",
		"meta-data.json",
		"user-data.json",
		"network-config.json",
		"meta-data.yaml",
		"user-data.yaml",
		"network-config.yaml",
	}

	// Vérifie d'abord si le volume s'appelle "config-2"
	volumeInfo, err := getVolumeLabel(path)
	if err == nil && strings.Contains(strings.ToLower(volumeInfo), "config-2") {
		fmt.Printf("  → Volume nommé 'config-2' détecté!\n")
		return true
	}

	// Vérifie la présence des fichiers
	for _, file := range commonFiles {
		fullPath := filepath.Join(path, file)
		if _, err := os.Stat(fullPath); err == nil {
			fmt.Printf("  → Fichier trouvé: %s\n", file)
			return true
		}
	}

	return false
}

// getVolumeLabel tente de lire le nom du volume
func getVolumeLabel(path string) (string, error) {
	// Cette fonction est un placeholder - sous Windows, nous devrions utiliser l'API Windows
	// pour obtenir le vrai nom du volume
	return "", nil
}

// ListCloudInitFiles returns a list of files in the Cloud-Init disk
func (d *DiskInfo) ListCloudInitFiles() error {
	fmt.Printf("\nLecture du contenu du disque %s...\n", d.Path)

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
