package cmd

import (
	"fmt"
	"os"

	"github.com/mitsimi/aocli/internal/aoc"
	"github.com/spf13/cobra"
)

var cfgFile string
var session string
var client *aoc.Client

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aocli",
	Short: "Convenient cli tool for Advent of Code to get going faster.",
	Long: `aocli is a convenient cli tool for Advent of Code so you never have to leave your editor.
It automatically can retreive the puzzle description and input and submit your answer.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if s := getSessionToken(cmd); s != "" {
			client = aoc.NewClient(s)
			return nil
		}

		return fmt.Errorf("no session token provided")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	rootCmd.PersistentFlags().StringVarP(&session, "session", "s", "", "session cookie from adventofcode.com")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// TODO load config from file
}

// getSessionToken returns the session token in order from the flag or config
func getSessionToken(cmd *cobra.Command) string {
	fmt.Printf("%+v\n", session)
	if session != "" {
		fmt.Println("Using session token from flag")
		return session
	}

	// TODO return session from config
	return ""
}
