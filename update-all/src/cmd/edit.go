package cmd

import (
	"fmt"
	"update_all/src/core"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newCmdEdit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Open config file",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Run edit command")
			if !core.IfRoutineFileExists() {
				routines := core.DefaultRoutines()
				core.FlushRoutines(routines)
				log.Info("No routine config file founded, create one first...")
			}
			core.Run("code", core.GetRoutineFile())
		},
	}
	cmd.Flags().Bool("path", false, "Show config path")
	return cmd
}
