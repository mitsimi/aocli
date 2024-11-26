package cmd

import (
	"fmt"
	"os"

	"github.com/mitsimi/aocli/internal/aoc"
	"github.com/spf13/cobra"
)

var cfgFile string

var client *aoc.Client

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aocli",
	Short: "Convenient cli tool for Advent of Code to get going faster.",
	Long: `aocli is a convenient cli tool for Advent of Code so you never have to leave your editor.
It automatically can retreive the puzzle description and input and submit your answer.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	client = aoc.NewClient(getSessionToken(rootCmd), aoc.WithTransport(aoc.NewDebugTransport()))
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	rootCmd.PersistentFlags().StringP("session", "s", "", "session cookie value from adventofcode.com")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// TODO load config from file
}

func getSessionToken(cmd *cobra.Command) string {
	if cmd.Flags().Changed("session") {
		fmt.Println("Using session token from flag")
		return cmd.Flag("session").Value.String()
	}

	// TODO return session from config
	return "session"
}
