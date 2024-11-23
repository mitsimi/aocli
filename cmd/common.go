package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// getDayFromCurrentDir returns the day number from the current working directory
func getDayFromCurrentDir() (string, error) {
	// Get working dir, which equals to the location from where the command is issued
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	base := filepath.Base(wd)
	if strings.HasPrefix(base, "day") && len(base) == 5 {
		return base[3:], nil // Extract day number
	}
	return "", fmt.Errorf("not in a day folder")
}
