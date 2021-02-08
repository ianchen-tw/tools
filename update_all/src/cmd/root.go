package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewCmdEdit Generate and edit command
func NewCmdEdit() *cobra.Command {
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

// NewCmdCInit crete defualt config file
func NewCmdCInit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create default config",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Run create command")
		},
	}
	return cmd
}

// NewCmdRoot Generate the root command
func NewCmdRoot() *cobra.Command {
	// TODO: add persistent option: debug
	// TODO: add persistent option: nohome
	// TODO: make option `dry` persistent
	cmd := &cobra.Command{
		Use:   "update-all",
		Short: "Update All",
		Long:  "Automatically run your routines",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	cmd.Flags().Bool("dry", false, "Dry run")
	cmd.Flags().Bool("force", false, "Force to run all routines")

	cmd.AddCommand(NewCmdEdit())
	cmd.AddCommand(NewCmdCInit())
	return cmd
}

var rootCmd = NewCmdRoot()

// Execute our rootCmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
