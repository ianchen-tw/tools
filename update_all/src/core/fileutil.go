package core

import "os"

// Flush byte data to file.
// Create one if not exists, overwrite otherwise
func flushToFile(filename string, data []byte) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
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
