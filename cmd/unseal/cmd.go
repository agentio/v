package unseal

import (
	"fmt"

	"github.com/spf13/cobra"
)

var verbose bool

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unseal",
		Short: "Unseal a local vault",
		RunE:  action,
	}
	return cmd
}

func action(cmd *cobra.Command, args []string) error {
	fmt.Printf("Everything looks good!\n")
	return nil
}
