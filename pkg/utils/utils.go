package utils

import "os"

// DestinationExists verifies that a destination exists, if supplied
func DestinationExists(path string) bool {
	if path == "." {
		return true
	}
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
