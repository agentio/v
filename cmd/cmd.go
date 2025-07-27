package cmd

import (
	"github.com/agentio/v/cmd/env"
	"github.com/agentio/v/cmd/read"
	"github.com/agentio/v/cmd/token"
	"github.com/agentio/v/cmd/unseal"
	"github.com/agentio/v/cmd/write"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "v",
		Short: "Vault Tools",
	}
	cmd.AddCommand(env.Cmd())
	cmd.AddCommand(read.Cmd())
	cmd.AddCommand(token.Cmd())
	cmd.AddCommand(unseal.Cmd())
	cmd.AddCommand(write.Cmd())
	return cmd
}
