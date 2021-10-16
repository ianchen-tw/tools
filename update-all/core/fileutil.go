package core

import (
	"log"
	"os"
	"path/filepath"
)

func ensureDirExists(fpath string) {
	os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
}

func getActualFileLoc(filename string) string {
	if UseWorkdirToFetch {
		mydir, _ := os.Getwd()
		return filepath.Join(mydir, filename)
	}
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}
	return filepath.Join(homedir, ".config", "update-all", filename)
}

// Flush byte data to file.
// Create one if not exists, overwrite otherwise
func flushToFile(filename string, data []byte) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	defer file.Sync()
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func ifFileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
