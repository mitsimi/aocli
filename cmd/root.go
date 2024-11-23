package cmd

import (
	"fmt"
	"os"

	"github.com/mitsimi/aocli/internal/aoc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	client = aoc.NewClient(viper.GetString("session"))


	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		if _, err := getDayFromCurrentDir(); err == nil {
			viper.AddConfigPath("..")
		} else {
			viper.AddConfigPath(".")
		}

		// Search config in home directory with name ".aocli" (without extension).
		viper.SetConfigType("toml")
		viper.SetConfigName(".aocli")
	}
	viper.SetDefault("year", getDefaultYear())

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
