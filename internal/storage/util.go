package storage

import (
	"os"
)

// pathExists checks if the provided path exists.
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) { // path does not exist
			return false, nil
		} else { // path exists, but some error occurred
			return false, err
		}
	}
	return true, nil // path exists
}
