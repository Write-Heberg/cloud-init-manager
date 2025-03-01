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
	// Liste des lettres de lecteur possibles
	driveLetters := []string{"A:", "B:", "C:", "D:", "E:", "F:", "G:", "H:", "I:", "J:", "K:", "L:", "M:",
		"N:", "O:", "P:", "Q:", "R:", "S:", "T:", "U:", "V:", "W:", "X:", "Y:", "Z:"}

	for _, drive := range driveLetters {
		path := drive + "\\"

		// Vérifie si le lecteur existe et est accessible
		_, err := os.Stat(path)
		if err == nil {
			// Vérifie si c'est notre disque cloud-init
			if isCloudInitDisk(path) {
				return &DiskInfo{
					Path:       path,
					MountPoint: path,
					Files:      []string{},
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("cloud-init disk not found")
}

// isCloudInitDisk vérifie si le lecteur contient les fichiers typiques de cloud-init
func isCloudInitDisk(path string) bool {
	// Vérifie la présence de fichiers typiques de cloud-init
	commonFiles := []string{
		"meta-data",
		"user-data",
		"network-config",
	}

	for _, file := range commonFiles {
		fullPath := filepath.Join(path, file)
		if _, err := os.Stat(fullPath); err == nil {
			return true
		}
	}

	return false
}

// ListCloudInitFiles returns a list of files in the Cloud-Init disk
func (d *DiskInfo) ListCloudInitFiles() error {
	files, err := filepath.Glob(filepath.Join(d.MountPoint, "*"))
	if err != nil {
		return fmt.Errorf("failed to list Cloud-Init files: %v", err)
	}

	d.Files = files
	return nil
}

// ReadCloudInitFile reads the content of a specific file from the Cloud-Init disk
func (d *DiskInfo) ReadCloudInitFile(filename string) ([]byte, error) {
	fullPath := filepath.Join(d.MountPoint, filename)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", filename, err)
	}
	return content, nil
}
