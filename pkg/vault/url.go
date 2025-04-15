package vault

import "os"

func URL(path string) string {
	addr := os.Getenv("VAULT_ADDR")
	return addr + path
}
