package path

import (
	"fmt"
	"os"
	"strings"
)

func ExpandDirAndEnv(path string) string {
	path = ExpandDir(path)
	return os.ExpandEnv(path)
}

func ExpandDir(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Invalid path: %v\n", err)
		os.Exit(-1)
	}

	return home + path[1:]
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
