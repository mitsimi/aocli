/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download the puzzle description, examples and inputs",
	Long:  `Download the puzzle description, examples and inputs`,
	Args:  cobra.NoArgs,
	Run:   executeDownload,
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().IntP("year", "y", 0, "puzzle year (year of current or last event. Can be specified in the config file)")
	downloadCmd.Flags().IntP("day", "d", 0, "puzzle day (current/last unlocked day (during Advent of Code month) or is inferred from the current folder)")

	downloadCmd.Flags().BoolP("descriptionOnly", "D", false, "only download the description")
	downloadCmd.Flags().BoolP("examplesOnly", "E", false, "only download the examples")
	downloadCmd.Flags().BoolP("inputsOnly", "I", false, "only download the inputs")
}

func executeDownload(cmd *cobra.Command, args []string) {
	fmt.Println("download called")
}
