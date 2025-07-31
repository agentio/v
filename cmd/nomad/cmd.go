package nomad

import (
	"github.com/agentio/v/cmd/nomad/setup"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "nomad",
	}
	cmd.AddCommand(setup.Cmd())

	return cmd
}
