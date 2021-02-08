package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newCmdRoot() *cobra.Command {
	// TODO: add persistent option: debug
	// TODO: add persistent option: nohome
	// TODO: make option `dry` persistent
	cmd := &cobra.Command{
		Use:   "update-all",
		Short: "Update All",
		Long:  "Automatically run your routines",
		Run:   startUpdateAll,
	}
	cmd.Flags().Bool("dry", false, "Dry run")
	cmd.Flags().BoolP("force", "f", false, "Force to run all routines")
	cmd.PersistentFlags().Bool("debug", false, "Start in debug mode")

	cmd.AddCommand(newCmdEdit())
	cmd.AddCommand(newCmdInit())
	return cmd
}

func startUpdateAll(cmd *cobra.Command, args []string) {
	isDebug, _ := cmd.Flags().GetBool("debug")
	fmt.Println("root: debug:", isDebug)
}

var rootCmd = newCmdRoot()

// Execute our rootCmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
