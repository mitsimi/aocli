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

		cmd.SilenceUsage = true
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
	conf = &config.Config{}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return
	}

	path, err := findConfigInDir(home)
	if err == nil {
		c, err := config.Parse(path)
		if err == nil {
			config.Merge(conf, c)
		}
	}

	path, err = findConfigInProject()
	if err == nil {
		c, err := config.Parse(path)
		if err == nil {
			config.Merge(conf, c)
		}
	}

	if cfgFlag != "" {
		if _, err := os.Stat(cfgFlag); err == nil {
			c, err := config.Parse(cfgFlag)
			if err == nil {
				config.Merge(conf, c)
			}
		}
	}
}

func findConfigInProject() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for i := 0; i < 3; i++ {
		entries, err := os.ReadDir(wd)
		if err != nil {
			return "", fmt.Errorf("Error reading directory: %v", err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			if strings.Contains(entry.Name(), ".aocli") {
				return filepath.Join(wd, entry.Name()), nil
			}
		}

		wd = filepath.Dir(wd)
	}

	return "", fmt.Errorf("No config file found in project")
}

// findConfigInDir searches for a config file in the given directory
func findConfigInDir(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.Contains(entry.Name(), ".aocli") {
			return filepath.Join(dir, entry.Name()), nil
		}
	}

	return "", errors.New("No config file found in home directory")
}

// getSessionToken returns the session token in order from the flag or config
func getSessionToken() string {
	if sessionFlag != "" {
		return sessionFlag
	}

	if conf.Session != "" {
		return conf.Session
	}

	return ""
}
