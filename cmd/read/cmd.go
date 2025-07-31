package read

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/agentio/v/pkg/vault"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read ENGINE",
		Short: "Read a file of secrets from a vault",
		RunE:  action,
		Args:  cobra.ExactArgs(1),
	}
	return cmd
}

func action(cmd *cobra.Command, args []string) error {
	engine := args[0]
	k, err := vault.ReadKeys()
	if err != nil {
		log.Printf("no configuration file")
		//return err
	}
	token, err := vault.TokenFromFile()
	if err != nil {
		token = k.RootToken
	}
	var responseData ListKVStoreMetadataResponse
	{
		url := vault.URL("/v1/" + engine + "/metadata?list=true")
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		request.Header.Set("Authorization", "Bearer "+token)
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
		err = json.Unmarshal(body, &responseData)
		if err != nil {
			return err
		}
	}
	secrets := make(map[string]map[string]string)
	for _, key := range responseData.Data.Keys {
		request, err := http.NewRequest("GET", vault.URL("/v1/"+engine+"/data/"+key), nil)
		if err != nil {
			return err
		}
		request.Header.Set("Authorization", "Bearer "+token)
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
		var responseData GetKVStoreEntryResponse
		err = json.Unmarshal(body, &responseData)
		if err != nil {
			return err
		}
		secrets[key] = responseData.Data.Data
	}
	b, err := json.MarshalIndent(secrets, "", "  ")
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(b)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write([]byte("\n"))
	return err
}

type ListKVStoreMetadataResponse struct {
	Data struct {
		Keys []string `json:"keys"`
	} `json:"data"`
}

type GetKVStoreEntryResponse struct {
	Data struct {
		Data map[string]string `json:"data"`
	} `json:"data"`
}
