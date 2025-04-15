package read

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Read a file of secrets from the vault",
		RunE:  action,
	}
	return cmd
}

func action(cmd *cobra.Command, args []string) error {
	return nil
}
