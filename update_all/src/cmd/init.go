package cmd

import (
	"update_all/src/core"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func initConfigFile(cmd *cobra.Command, args []string) {
	writeNewFile := true
	if core.IfRoutineFileExists() {
		log.Info("File exists")
		survey.AskOne(&survey.Confirm{
			Message: "Override existing config file?",
		}, &writeNewFile)
	}
	if writeNewFile {
		routines := core.DefaultRoutines()
		core.FlushRoutines(routines)
		log.Info("Write default config successfully")
	} else {
		log.Warn("Keep existing config file")
	}
}
func newCmdInit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create default config",
		Run:   initConfigFile,
	}
	return cmd
}
