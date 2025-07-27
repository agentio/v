package env

import (
	"fmt"

	"github.com/agentio/v/pkg/vault"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "export VAULT_ADDR and VAULT_TOKEN environment variables",
		RunE:  action,
		Args:  cobra.NoArgs,
	}
	return cmd
}

func action(cmd *cobra.Command, args []string) error {
	k, err := vault.ReadKeys()
	if err != nil {
		return err
	}

	fmt.Printf("export VAULT_ADDR=%s\n", "http://localhost:8200")
	fmt.Printf("export VAULT_TOKEN=%s\n", k.RootToken)
	return err
}
