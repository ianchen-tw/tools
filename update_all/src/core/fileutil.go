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
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	file.Sync()
	return nil
}
