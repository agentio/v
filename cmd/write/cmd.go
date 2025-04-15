package write

import "github.com/spf13/cobra"

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "write",
		Short: "Write a file of secrets into the vault",
		RunE:  action,
	}
	return cmd
}

func action(cmd *cobra.Command, args []string) error {
	return nil
}
