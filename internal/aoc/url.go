package aoc

import "fmt"

const BaseURL = "https://adventofcode.com"

func CalendarURL(year int) string {
	return fmt.Sprintf("%s/%d", BaseURL, year)
}

func DayURL(year, day int) string {
	return fmt.Sprintf("%s/%d/day/%d", BaseURL, year, day)
}

func InputURL(year, day int) string {
	return fmt.Sprintf("%s/%d/day/%d/input", BaseURL, year, day)
}

func SubmitURL(year, day int) string {
	return fmt.Sprintf("%s/%d/day/%d/answer", BaseURL, year, day)
}
