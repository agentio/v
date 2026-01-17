package get

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/agentio/v/pkg/vault"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get COLLECTION SECRET",
		Short: "Get a secret from a vault",
		RunE:  action,
		Args:  cobra.ExactArgs(2),
	}
	return cmd
}

func action(cmd *cobra.Command, args []string) error {
	collection := args[0]
	name := args[1]
	k, err := vault.ReadKeys("")
	if err != nil {
		log.Printf("no configuration file")
		//return err
	}
	token, err := vault.TokenFromFile()
	if err != nil {
		token = k.RootToken
	}

	request, err := http.NewRequest("GET", vault.URL("/v1/"+collection+"/data/"+name), nil)
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
	var secret string
	for _, v := range responseData.Data.Data {
		secret = v
		break
	}

	cmd.OutOrStdout().Write([]byte(secret))

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
