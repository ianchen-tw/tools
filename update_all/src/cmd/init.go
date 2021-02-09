package cmd

import (
	"fmt"
	"update_all/src/core"

	"github.com/spf13/cobra"
)

func newCmdInit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create default config",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: if file exits, ask user to confirm override it
			routines := core.DefaultRoutines()
			core.FlushRoutines(routines)
			fmt.Println("Successfully init config")
		},
	}
	return cmd
}
