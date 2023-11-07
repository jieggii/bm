// Package env serves the purpose of reading environmental variables.
package env

import (
	"os"
	"strings"
)

const VarStorageHome = "BM_STORAGE_HOME"
const VarStorageFileName = "BM_STORAGE_FILE_NAME"

const DefaultStorageHome = "~/.local/share/bm"
const DefaultStorageFileName = "books.json"

// Env represents necessary environmental variables.
type Env struct {
	// Path to the storage home.
	StorageHome string

	// Path to the storage file name.
	StorageFileName string
}

// Read reads necessary environmental variables.
func Read() *Env {
	storageHome, found := os.LookupEnv(VarStorageHome)
	if !found {
		storageHome = DefaultStorageHome

		// replace "~" symbol in the default storage home path to user's
		// home directory if possible:
		userHome, err := os.UserHomeDir()
		if err == nil {
			storageHome = strings.Replace(storageHome, "~", userHome, 1)
		}
	}

	storageFileName, found := os.LookupEnv(VarStorageFileName)
	if !found {
		storageFileName = DefaultStorageFileName
	}

	return &Env{
		StorageHome:     storageHome,
		StorageFileName: storageFileName,
	}
}
