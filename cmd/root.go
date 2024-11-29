package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitsimi/aocli/internal/aoc"
	"github.com/mitsimi/aocli/internal/config"
	"github.com/spf13/cobra"
)

var cfgFlag string
var sessionFlag string

var conf *config.Config
var client *aoc.Client

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aocli",
	Short: "Convenient cli tool for Advent of Code to get going faster.",
	Long: `aocli is a convenient cli tool for Advent of Code so you never have to leave your editor.
It automatically can retreive the puzzle description and input and submit your answer.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if s := getSessionToken(); s != "" {
			client = aoc.NewClient(s)
			return nil
		}

		return fmt.Errorf("no session token provided")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Execute()
	fmt.Println("")
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFlag, "config", "c", "", "config file")
	rootCmd.PersistentFlags().StringVarP(&sessionFlag, "session", "s", "", "session cookie from adventofcode.com")
}

func initConfig() {
	if cfgFlag != "" {
		if _, err := os.Stat(cfgFlag); err == nil {
			c, err := config.Parse(cfgFlag)
			if err == nil {
				conf = c
				return
			}
		}
	}

	// because no config file was provided through the flag, we first look in the project directory
	// because the command can be executed from the root, the year folder depending on the config or the day folder
	// we look in the current directory and two parent directories
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 3; i++ {
		entries, err := os.ReadDir(wd)
		if err != nil {
			fmt.Println("Error reading directory:", err)
			return
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if strings.Contains(entry.Name(), ".aocli") {
				c, err := config.Parse(filepath.Join(wd, entry.Name()))
				if err == nil {
					conf = c
					return
				}
			}
		}

		wd = filepath.Dir(wd)
	}

	// if we still have no config, we look in the home and home/.config directories
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = filepath.WalkDir(home, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return filepath.SkipDir
		}

		base := filepath.Base(path)
		if strings.Contains(base, ".aocli") {
			c, err := config.Parse(path)
			if err == nil {
				conf = c
				return filepath.SkipAll
			}
		}

		return nil
	})
	if errors.Is(err, filepath.SkipAll) {
		return
	}

	err = filepath.WalkDir(filepath.Join(home, ".config"), func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return filepath.SkipDir
		}

		base := filepath.Base(path)
		if strings.Contains(base, ".aocli") {
			c, err := config.Parse(path)
			if err == nil {
				conf = c
				return filepath.SkipAll
			}
		}

		return nil
	})
	if errors.Is(err, filepath.SkipAll) {
		return
	}
}

// getSessionToken returns the session token in order from the flag or config
func getSessionToken() string {
	if sessionFlag != "" {
		//fmt.Println("Using session token from flag")
		return sessionFlag
	}

	// TODO return session from config
	return ""
}
