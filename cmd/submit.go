package cmd

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/mitsimi/aocli/internal/aoc"
	"github.com/spf13/cobra"
)

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit [flags] [answer]",
	Short: "Submit puzzle answer",
	Long: `Submit your puzzle answer without leaving your editor. 
The answer may be provided as an argument or through the file flag. You also can pipe your answer into the command.`,
	Args:      cobra.MaximumNArgs(1),
	ValidArgs: []string{"answer"},
	RunE:      executeSubmit,
}

func init() {
	rootCmd.AddCommand(submitCmd)

	var yearVal int = time.Now().Year()
	if time.Now().Month() != time.December {
		// default to lasts years because this years isn't open yet
		yearVal -= 1
	}
	submitCmd.Flags().IntP("year", "y", yearVal, "Puzzle year. Defaults to year of current or last Advent of Code event. Can be specified in the config file.")

	var dayVal int = time.Now().Day()
	if unlocked, _ := aoc.IsDayUnlocked(yearVal, dayVal); !unlocked {
		// default to lasts years because this years isn't open yet
		dayVal -= 1
	}
	submitCmd.Flags().IntP("day", "d", dayVal, "Puzzle day. Defaults to current/last unlocked day (during Advent of Code month) or is inferred from the current folder")

	submitCmd.Flags().IntP("level", "l", 1, "Puzzle level. Defaults to 1")

	submitCmd.Flags().StringP("file", "f", "", "File containing the answer")
}

func executeSubmit(cmd *cobra.Command, args []string) error {
	answer, err := getAnswer(cmd, args)
	if err != nil {
		return err
	}

	level, _ := cmd.Flags().GetInt("level")

	year, _ := cmd.Flags().GetInt("year")

	day, _ := cmd.Flags().GetInt("day")

	outcome, err := client.SubmitAnswer(aoc.Level(level), year, day, answer)
	if err != nil {
		return fmt.Errorf("Error submitting answer: %v", err)
	}

	fmt.Printf("%s\n", outcome)
	return nil
}

// getAnswer returns the answer from the stdin, file or argument
func getAnswer(cmd *cobra.Command, args []string) (string, error) {
	answer, err := readStdin()
	if err != nil {
		return "", err
	}

	if answer != "" {
		return answer, nil
	}

	path, _ := cmd.Flags().GetString("file")
	// check if there is a file
	if path != "" {
		// read the file
		file, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(file), nil
	}

	if len(args) == 0 {
		return "", fmt.Errorf("No answer provided")
	}

	return args[0], nil
}

// checkStdin checks if stdin has data available
func checkStdin() bool {
	// Check if stdin has data
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	// If data is available (piped input), return true
	if info.Mode()&os.ModeCharDevice == 0 {
		return true
	}
	// If no data is available (interactive), return false
	return false
}

// readStdin reads the data from stdin if available and prints it
func readStdin() (string, error) {
	if !checkStdin() {
		return "", nil
	}

	var data string
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		data = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("Error reading stdin: %v", err)
	}

	return data, nil
}