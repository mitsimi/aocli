/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/mitsimi/aocli/internal/aoc"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download the puzzle description, examples and inputs",
	Long: `Download the puzzle description, examples and inputs.
The files will be saved in the current folder or in the folder specified by the output flag.`,
	Args: cobra.NoArgs,
	RunE: executeDownload,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// we download the input by default and if specified, because the input is generated per account we need a session token for it
		if (!contentFlagsChanged(cmd) || cmd.Flag("input").Changed) && getSessionToken() == "" {
			cmd.SilenceUsage = true
			return errors.New("a session token is required to download the input")
		}

		client = aoc.NewClient(getSessionToken())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().IntP("year", "y", 0, "puzzle year (year of current or last event. Can be specified in the config file)")
	downloadCmd.Flags().IntP("day", "d", 0, "puzzle day (current/last unlocked day (during Advent of Code month) or is inferred from the current folder)")

	downloadCmd.Flags().BoolP("description", "D", false, "download the description")
	downloadCmd.Flags().BoolP("examples", "E", false, "download the examples")
	downloadCmd.Flags().BoolP("input", "I", false, "download the input")

	downloadCmd.Flags().StringP("output", "o", "", "output folder (default is the current folder)")
}

func executeDownload(cmd *cobra.Command, args []string) error {
	dir, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}
	if dir == "" {
		dir = "."
	}

	if !contentFlagsChanged(cmd) {
		cmd.Flag("description").Value.Set("true")
		cmd.Flag("examples").Value.Set("true")
		cmd.Flag("input").Value.Set("true")
	}

	year := getYear(cmd)
	day := getDay(cmd)

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Failed to get current directory: %v", err)
	}

	yearReg := regexp.MustCompile(`\d{4}`)
	if yearReg.MatchString(filepath.Base(currentDir)) {
		year, _ = strconv.Atoi(yearReg.FindString(filepath.Base(currentDir)))
	} else if yearReg.MatchString(filepath.Dir(currentDir)) {
		year, _ = strconv.Atoi(yearReg.FindString(filepath.Dir(currentDir)))
	}

	if ok, _ := cmd.Flags().GetBool("description"); ok {
		cmd.Println("Downloading description...")
		err = downloadDescription(year, day, dir)
		if err != nil {
			return err
		}
	}

	if ok, _ := cmd.Flags().GetBool("examples"); ok {
		cmd.Println("Downloading examples...")
		err = downloadExample(year, day, dir)
		if err != nil {
			return err
		}
	}

	if ok, _ := cmd.Flags().GetBool("input"); ok {
		cmd.Println("Downloading input...")
		err = downloadInput(year, day, dir)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadDescription(year, day int, dir string) error {
	content, err := client.GetDescription(year, day)
	if err != nil {
		return err
	}

	md, err := content.ToMarkdown(year)
	if err != nil {
		return err
	}

	err = writeStringToFile(filepath.Join(dir, "description.md"), md)
	if err != nil {
		return err
	}

	return nil
}

func downloadExample(year, day int, dir string) error {
	content, err := client.GetExample(year, day)
	if err != nil {
		return err
	}

	err = writeStringToFile(filepath.Join(dir, "example"), content)
	if err != nil {
		return err
	}

	return nil
}

func downloadInput(year, day int, dir string) error {
	content, err := client.GetInput(year, day)
	if err != nil {
		return err
	}

	err = writeStringToFile(filepath.Join(dir, "input"), string(content))
	if err != nil {
		return err
	}

	return nil
}

func contentFlagsChanged(cmd *cobra.Command) bool {
	return cmd.Flag("description").Changed || cmd.Flag("examples").Changed || cmd.Flag("input").Changed
}

func writeStringToFile(filePath, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(string(content))
	if err != nil {
		return err
	}
	return nil
}
