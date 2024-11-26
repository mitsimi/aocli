package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mitsimi/aocli/internal/aoc"
	"github.com/spf13/cobra"
)

func getNewClient(cmd *cobra.Command) *aoc.Client {
	return aoc.NewClient(getSessionToken(cmd))
}

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

// getDefaultDay returns the latest day where an event is available
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

// getDay returns the day from the flag or uses the default.
// The default is the current or last unlocked day.
func getDay(cmd *cobra.Command) int {
	if day, _ := cmd.Flags().GetInt("day"); day != 0 {
		return day
	}

	return getDefaultDay()
}

// getYear returns the year from the flag, config or default.
// The default is the current or last event year
func getYear(cmd *cobra.Command) int {
	// function to use to convert the year shorthand ex. 19 to 2019
	// we do not care about the century before 2000 because advent of code started in 2015
	conv := func(year int) int {
		if year < 100 {
			return year + 2000
		}
		return year
	}

	if year, _ := cmd.Flags().GetInt("year"); year != 0 {
		return conv(year)
	}

	// if year := viper.GetInt("year"); year != 0 {
	// 	return conv(year)
	// }

	return getDefaultYear()
}

// getSessionToken returns the session token in order from the flag or config
func getSessionToken(cmd *cobra.Command) string {
	fmt.Printf("%+v\n", session)
	if session != "" {
		fmt.Println("Using session token from flag")
		return session
	}

	// TODO return session from config
	return "some session"
}
