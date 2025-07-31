package vault

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type VaultKeys struct {
	Keys       []string `json:"keys"`
	KeysBase64 []string `json:"keys_base64"`
	RootToken  string   `json:"root_token"`
	Cluster    []string `json:"cluster"`
}

func ReadKeys(cluster string) (*VaultKeys, error) {
	if cluster == "" {
		cluster = "keys"
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	keyfile := filepath.Join(home, ".config", "vault", cluster+".json")
	b, err := os.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}
	keys := &VaultKeys{}
	err = json.Unmarshal(b, keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}
