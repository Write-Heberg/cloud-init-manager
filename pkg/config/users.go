package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// UserConfig représente la configuration d'un utilisateur
type UserConfig struct {
	Name              string   `json:"name"`
	Password          string   `json:"passwd"`
	SSHAuthorizedKeys []string `json:"ssh_authorized_keys"`
	Sudo              string   `json:"sudo"`
}

// ApplyUserConfig applique la configuration utilisateur
func ApplyUserConfig(config *UserConfig) error {
	// Vérifie si l'utilisateur existe
	userExists := checkUserExists(config.Name)

	if !userExists {
		// Crée l'utilisateur s'il n'existe pas
		if err := createUser(config.Name, config.Password); err != nil {
			return fmt.Errorf("erreur lors de la création de l'utilisateur: %v", err)
		}
		fmt.Printf("✓ Utilisateur créé: %s\n", config.Name)
	} else if config.Password != "" {
		// Met à jour le mot de passe si fourni
		if err := setPassword(config.Name, config.Password); err != nil {
			return fmt.Errorf("erreur lors de la modification du mot de passe: %v", err)
		}
		fmt.Printf("✓ Mot de passe mis à jour pour: %s\n", config.Name)
	}

	// Configure les clés SSH
	if len(config.SSHAuthorizedKeys) > 0 {
		if err := configureSSHKeys(config.Name, config.SSHAuthorizedKeys); err != nil {
			return fmt.Errorf("erreur lors de la configuration SSH: %v", err)
		}
	}

	return nil
}

// checkUserExists vérifie si l'utilisateur existe
func checkUserExists(username string) bool {
	cmd := exec.Command("net", "user", username)
	return cmd.Run() == nil
}

// createUser crée un nouvel utilisateur
func createUser(username, password string) error {
	cmd := exec.Command("net", "user", username, password, "/add")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erreur net user: %v, output: %s", err, string(output))
	}
	return nil
}

// setPassword modifie le mot de passe d'un utilisateur
func setPassword(username, password string) error {
	cmd := exec.Command("net", "user", username, password)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erreur net user: %v, output: %s", err, string(output))
	}
	return nil
}

// configureSSHKeys configure les clés SSH autorisées pour l'utilisateur
func configureSSHKeys(username string, keys []string) error {
	// Crée le dossier .ssh
	sshDir := filepath.Join("C:\\Users", username, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return fmt.Errorf("erreur création dossier SSH: %v", err)
	}

	// Écrit les clés dans authorized_keys
	authKeysPath := filepath.Join(sshDir, "authorized_keys")
	content := ""
	for _, key := range keys {
		content += key + "\n"
	}

	if err := os.WriteFile(authKeysPath, []byte(content), 0600); err != nil {
		return fmt.Errorf("erreur écriture authorized_keys: %v", err)
	}

	fmt.Printf("✓ Clés SSH configurées pour: %s\n", username)
	return nil
}
