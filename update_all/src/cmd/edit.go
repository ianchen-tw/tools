package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCmdEdit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Open config file",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Run edit command")
		},
	}
	cmd.Flags().Bool("path", false, "Show config path")
	return cmd
}
