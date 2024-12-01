package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

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

	dayFolder := fmt.Sprintf("day%02d", day)
	targetDir := currentDir

	// if the current directory is a day folder, then we must set the target directory to the parent folder so we don't create a nested day folder
	if regexp.MustCompile(`day\d{2}`).MatchString(filepath.Base(currentDir)) {
		targetDir = filepath.Dir(currentDir)
	}

	if conf.Structure == "multi-year" {
		yearReg := regexp.MustCompile(`\d{4}`)
		if yearReg.MatchString(filepath.Base(currentDir)) { // check if the current directory is a year folder
			// if yes we set the year to the year of the current directory
			year, _ = strconv.Atoi(yearReg.FindString(filepath.Base(currentDir)))
		} else {
			// if not we set the target directory to the year folder
			targetDir = filepath.Join(targetDir, fmt.Sprintf("%d", year))
		}
	}

	targetDir = filepath.Join(targetDir, dayFolder)

	if err := createFolders(targetDir); err != nil {
		cmd.PrintErrln("Failed to create folders:", err)
		return
	}

	templateDir := filepath.Join(currentDir, "template")
	if _, err := os.Stat(templateDir); !os.IsNotExist(err) {
		cmd.Println("Copying template files...")
		if err := template.CopyContent(templateDir, targetDir); err != nil {
			cmd.PrintErrln("Failed to copy template files:", err)
			return
		}
	}

	cmd.Println("Downloading puzzle data...")
	if err := downloadPuzzleData(year, day, targetDir); err != nil {
		cmd.PrintErrln("Failed to download puzzle data:", err)
		return
	}

	cmd.Println("Finished successfully!")
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

	err = downloadExample(year, day, destDir)
	if err != nil {
		return err
	}

	return nil
}

func createFolders(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
