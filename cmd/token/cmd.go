package token

import (
	"time"

	"github.com/agentio/v/pkg/vault"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Copy the vault token to the pasteboard",
		RunE:  action,
		Args:  cobra.NoArgs,
	}
	return cmd
}

func action(cmd *cobra.Command, args []string) error {
	k, err := vault.ReadKeys("")
	if err != nil {
		return err
	}

	err = clipboard.Init()
	if err != nil {
		return err
	}

	clipboard.Write(clipboard.FmtText, []byte(k.RootToken))
	time.Sleep(100 * time.Millisecond)
	return err
}
