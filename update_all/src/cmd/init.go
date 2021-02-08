package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCmdInit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create default config",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Run create command")
		},
	}
	return cmd
}
