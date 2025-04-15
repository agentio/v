package unseal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var verbose bool

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unseal",
		Short: "Unseal a local vault",
		RunE:  action,
	}
	return cmd
}

type VaultKeys struct {
	Keys       []string `json:"keys"`
	KeysBase64 []string `json:"keys_base64"`
	RootToken  string   `json:"root_token"`
}

type UnsealRequest struct {
	Key string `json:"key"`
}

func action(cmd *cobra.Command, args []string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configfilename := filepath.Join(home, ".config", "vault", "keys.json")
	b, err := os.ReadFile(configfilename)
	if err != nil {
		return err
	}
	var keys VaultKeys
	err = json.Unmarshal(b, &keys)
	if err != nil {
		return err
	}
	unseal := UnsealRequest{Key: keys.Keys[0]}
	unsealBytes, err := json.Marshal(unseal)
	if err != nil {
		return err
	}

	addr := os.Getenv("VAULT_ADDR")
	url := addr + "/v1/sys/unseal"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(unsealBytes))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	log.Printf("%+v", response)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	log.Printf("%s", string(body))
	fmt.Printf("Everything looks good!\n")
	return nil
}
