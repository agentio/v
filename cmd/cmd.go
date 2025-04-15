package cmd

import (
	"github.com/agentio/v/cmd/unseal"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "v",
		Short: "Vault Tools",
	}
	cmd.AddCommand(unseal.Cmd())
	return cmd
}
