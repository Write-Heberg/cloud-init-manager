package disk

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

// DiskInfo represents information about a Cloud-Init disk
type DiskInfo struct {
	Path       string
	MountPoint string
	Files      []string
}

// DetectCloudInitDisk attempts to find the Cloud-Init disk on the system
func DetectCloudInitDisk() (*DiskInfo, error) {
	// Get all available drives
	drives, err := windows.GetLogicalDrives()
	if err != nil {
		return nil, fmt.Errorf("failed to get logical drives: %v", err)
	}

	// Check each drive for config-2
	for i := 0; i < 26; i++ {
		if drives&(1<<uint(i)) != 0 {
			driveLetter := string(rune('A' + i))
			path := fmt.Sprintf("%s:\\", driveLetter)

			// Check if this drive is named "config-2"
			volumeName := make([]uint16, windows.MAX_PATH+1)
			err := windows.GetVolumeInformation(
				windows.StringToUTF16Ptr(path),
				&volumeName[0],
				uint32(len(volumeName)),
				nil,
				nil,
				nil,
				nil,
				0,
			)
			if err == nil {
				name := windows.UTF16ToString(volumeName[:])
				if name == "config-2" {
					return &DiskInfo{
						Path:       path,
						MountPoint: path,
						Files:      []string{},
					}, nil
				}
			}
		}
	}

	// Si nous n'avons pas trouvé le disque, retournons une erreur
	return nil, fmt.Errorf("cloud-init disk (config-2) not found")
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
