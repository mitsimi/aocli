package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mitsimi/aocli/internal/aoc"
)

// getDayFromCurrentDir returns the day number from the current working directory
func getDayFromCurrentDir() (string, error) {
	// Get working dir, which equals to the location from where the command is issued
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	base := filepath.Base(wd)
	base = strings.ToLower(base)
	if strings.HasPrefix(base, "day") && len(base) == 5 {
		return base[3:], nil // Extract day number
	}
	return "", fmt.Errorf("not in a day folder")
}

// getDefaultYear returns the latest year where a event is available
func getDefaultYear() int {
	yearVal := time.Now().Year()
	if time.Now().Month() != time.December {
		// default to lasts years because this years isn't open yet
		yearVal -= 1
	}
	return yearVal
}

// getDefaultDay returns the latest day where a event is available
func getDefaultDay() int {
	if d, err := getDayFromCurrentDir(); err == nil {
		n, err := strconv.Atoi(d)
		if err == nil {
			return n
		}
	}

	var dayVal int = time.Now().Day()
	if unlocked, _ := aoc.IsDayUnlocked(getDefaultYear(), dayVal); !unlocked {
		// default to lasts years because this years isn't open yet
		dayVal -= 1
	}
	return dayVal
}
