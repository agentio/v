package unseal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/agentio/v/pkg/pretty"
	"github.com/agentio/v/pkg/vault"
	"github.com/spf13/cobra"
)

var keyfile string

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unseal",
		Short: "Unseal a local vault",
		RunE:  action,
	}
	cmd.Flags().StringVar(&keyfile, "cluster", "", "cluster name")
	return cmd
}

type UnsealRequest struct {
	Key string `json:"key"`
}

func action(cmd *cobra.Command, args []string) error {
	k, err := vault.ReadKeys(keyfile)
	if err != nil {
		return err
	}
	unseal := UnsealRequest{}
	if len(k.Keys) > 0 {
		unseal.Key = k.Keys[0]
	} else if len(k.KeysBase64) > 0 {
		unseal.Key = k.KeysBase64[0]
	}
	unsealBytes, err := json.Marshal(unseal)
	if err != nil {
		return err
	}

	if len(k.Cluster) > 0 {
		for _, addr := range k.Cluster {
			request, err := http.NewRequest("POST", addr+"/v1/sys/unseal", bytes.NewBuffer(unsealBytes))
			if err != nil {
				return err
			}
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
			fmt.Printf("%s\n", pretty.JSON(body))
		}
		return nil
	}

	request, err := http.NewRequest("POST", vault.URL("/v1/sys/unseal"), bytes.NewBuffer(unsealBytes))
	if err != nil {
		return err
	}
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
	fmt.Printf("%s\n", pretty.JSON(body))
	return nil
}
