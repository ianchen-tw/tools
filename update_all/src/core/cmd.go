package core

import (
	"fmt"
	"os"
	"os/exec"
)

// Run run shell command
func Run(args ...string) error {
	exe, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}
	fmt.Println("Run:", args)
	cmd := exec.Command(exe, args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
