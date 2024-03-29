package cmd

import (
	"fmt"
	"os"
	"update-all/core"

	"github.com/kyokomi/emoji/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// FlagDebug Debug mode
	flagDebug bool

	//FlagNoHome Do not fetch user home folder
	flagNoHome bool

	//FlagDryRun Do not execute routine
	flagSkipExecute bool

	//FlagForceUpdate Ignore minimum run period
	flagForceUpdate bool
)

func newCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-all",
		Short: "Update All",
		Long:  "Automatically run your routines",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.Info(emoji.Sprint(":sparkles:Started!!"))
			if flagDebug {
				log.SetLevel(log.DebugLevel)
				log.Debug("Start cobra PreRun hook")
				log.Info("Start in Debug mode")
				core.IsDebug = true
			}
			if flagNoHome {
				core.UseWorkdirToFetch = true
			}
		},
		Run: startUpdateAll,
	}
	cmd.PersistentFlags().BoolVarP(&flagSkipExecute, "dry", "", false, "Dry run, do not execute routines")
	cmd.PersistentFlags().BoolVarP(&flagForceUpdate, "force", "f", false, "Force to run all routines")
	cmd.PersistentFlags().BoolVarP(&flagDebug, "debug", "", false, "Start in debug mode")
	cmd.PersistentFlags().BoolVarP(&flagNoHome, "nohome", "", false, "use working dir to store/read config")

	cmd.AddCommand(newCmdEdit())
	cmd.AddCommand(newCmdInit())
	return cmd
}

func startUpdateAll(cmd *cobra.Command, args []string) {
	cache := core.CreateRecordMap()
	err := cache.TryLoad()
	if err != nil {
		log.Warn(err)
	}
	log.Debug("Cache state: ", cache)
	defer cache.Flush()

	routines, err := core.LoadRoutines()
	if err != nil {
		// Can't find routine file
		log.Error("Unable to find file: ", core.GetRoutineFile())
		log.Error("Use `update-all init` to create a config file first")
		os.Exit(1)
	}
	for _, routine := range routines {
		cache.RunRoutineIfOutdated(routine, flagForceUpdate, flagSkipExecute)
	}
	log.Info(emoji.Sprint(":tada: Finished!!"))

}

var rootCmd = newCmdRoot()

// Execute our rootCmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
