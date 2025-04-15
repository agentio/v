package write

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/agentio/v/pkg/pretty"
	"github.com/agentio/v/pkg/vault"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "write ENGINE FILE",
		Short: "Write a JSON file of secrets into a vault",
		RunE:  action,
		Args:  cobra.ExactArgs(2),
	}
	return cmd
}

func action(cmd *cobra.Command, args []string) error {
	engine := args[0]
	keys, err := vault.ReadKeys()
	if err != nil {
		return err
	}
	b, err := os.ReadFile(args[1])
	if err != nil {
		return err
	}
	secrets := make(map[string]map[string]string)
	err = json.Unmarshal(b, &secrets)
	if err != nil {
		return err
	}
	for k, v := range secrets {
		b, err := json.Marshal(PutKVStoreMetadataRequest{Data: v})
		if err != nil {
			return err
		}
		request, err := http.NewRequest("PUT", vault.URL("/v1/"+engine+"/data/"+k), bytes.NewBuffer(b))
		if err != nil {
			return err
		}
		request.Header.Set("Authorization", "Bearer "+keys.RootToken)
		request.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", string(pretty.JSON(body)))
	}
	return nil
}

type PutKVStoreMetadataRequest struct {
	Data map[string]string `json:"data"`
}
