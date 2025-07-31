package setup

import (
	"github.com/agentio/v/cmd/nomad/setup/vault"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "setup",
	}

	cmd.AddCommand(vault.Cmd())
	return cmd
}
