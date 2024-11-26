package cmd

import (
	"fmt"
	"os"

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
	fmt.Println("new called")
	// check config for wanted files structure
	// check current folder for creation process

	// create folders
	// copy template files
	// download description, examples and inputs
}

func createFolders(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
