package aoc

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Level int

const (
	LevelOne Level = 1
	LevelTwo
)

func ParseLevel(s string) Level {
	if s == "2" {
		return LevelTwo
	}
	return LevelOne
}

type SubmissionOutcome int

const (
	SubmissionCorrect SubmissionOutcome = iota
	SubmissionIncorrectTooHigh
	SubmissionIncorrectTooLow
	SubmissionWait
	SubmissionWrongLevel
	SubmissionOthersAnswer
	SubmissionError
)

func (so SubmissionOutcome) String() string {
	switch so {
	case SubmissionCorrect:
		return "Correct answer"
	case SubmissionIncorrectTooHigh:
		return "Incorrect answer. Answer is too high."
	case SubmissionIncorrectTooLow:
		return "Incorrect answer. Answer is too low."
	case SubmissionWait:
		return "Wait a bit before submitting again"
	case SubmissionWrongLevel:
		return "You are solving the wrong level"
	case SubmissionOthersAnswer:
		return "Incorrect answer, but the correct one for someone else."
	case SubmissionError:
		fallthrough
	default:
		return "Error submitting answer"
	}
}

type UnknownResponseError struct {
	StatusCode int
	Response   string
}

func (e UnknownResponseError) Error() string {
	return fmt.Sprintf("Unknown response: %d \n%s", e.StatusCode, e.Response)
}

func (c *Client) SubmitAnswer(level Level, year, day int, answer string) (SubmissionOutcome, error) {
	data := url.Values{}
	data.Set("level", fmt.Sprintf("%d", level))
	data.Set("answer", fmt.Sprintf("%s", answer))

	req, err := http.NewRequest("POST", SubmitURL(year, day), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return SubmissionError, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Request(req)
	if err != nil {
		return SubmissionError, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return SubmissionError, fmt.Errorf("Failed to parse HTML: %v", err)
	}

	outcome := doc.Find("main > article > p").Text()

	switch {
	case strings.Contains(outcome, "That's the right answer"):
		return SubmissionCorrect, nil
	case strings.Contains(outcome, "That's not the right answer"):
		return SubmissionIncorrect, nil
	case strings.Contains(outcome, "You gave an answer too recently"):
		return SubmissionWait, nil
	case strings.Contains(outcome, "You don't seem to be solving the right level"):
		return SubmissionWrongLevel, nil
	case strings.Contains(outcome, "for someone else"):
		return SubmissionOthersAnswer, nil
	default:
		return SubmissionError, UnknownResponseError{StatusCode: resp.StatusCode, Response: outcome}
	}
}
