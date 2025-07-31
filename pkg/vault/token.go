package vault

import (
	"os"
	"strings"
)

func TokenFromFile() (string, error) {
	b, err := os.ReadFile("secrets/vault_token")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}
