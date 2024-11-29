package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/mitsimi/aocli/internal/template"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new folder for the day",
	Long: `Creates a new folder with the contents of the template folder if present.
It will save the description, examples and inputs automatically into seperate files.`,
	Args: cobra.NoArgs,
	Run:  executeNew,
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().IntP("year", "y", 0, "puzzle year (year of current or last event. Can be specified in the config file)")
	newCmd.Flags().IntP("day", "d", 0, "puzzle day (current/last unlocked day (during Advent of Code month) or is inferred from the current folder)")
}

func executeNew(cmd *cobra.Command, args []string) {
	year := getYear(cmd)
	day := getDay(cmd)

	currentDir, err := os.Getwd()
	if err != nil {
		cmd.PrintErrln("Failed to get current directory:", err)
		return
	}

	yearFolder := fmt.Sprintf("%d", year)
	dayFolder := fmt.Sprintf("day%02d", day)
	targetDir := currentDir

	if regexp.MustCompile(`day\d{2}`).MatchString(filepath.Base(currentDir)) {
		targetDir = filepath.Dir(currentDir)
	}

	// if the current directory is not the year folder, then we must be inside the root folder
	if filepath.Base(currentDir) != yearFolder {
		targetDir = filepath.Join(currentDir, yearFolder)
	}

	targetDir = filepath.Join(targetDir, dayFolder)

	if err := createFolders(targetDir); err != nil {
		cmd.PrintErrln("Failed to create folders:", err)
		return
	}

	templateDir := filepath.Join(currentDir, "template")
	if _, err := os.Stat(templateDir); !os.IsNotExist(err) {
		if err := template.CopyContent(templateDir, targetDir); err != nil {
			cmd.PrintErrln("Failed to copy template files:", err)
			return
		}
	}

	if err := downloadPuzzleData(year, day, targetDir); err != nil {
		cmd.PrintErrln("Failed to download puzzle data:", err)
		return
	}

	cmd.Println("New folder created successfully")
}

func downloadPuzzleData(year, day int, destDir string) (err error) {
	err = downloadDescription(year, day, destDir)
	if err != nil {
		return err
	}

	err = downloadInput(year, day, destDir)
	if err != nil {
		return err
	}

	err = downloadExamples(year, day, destDir)
	if err != nil {
		return err
	}

	return nil
}

func createFolders(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
