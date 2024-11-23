package aoc

import (
	"bytes"
	"errors"
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
	SubmissionIncorrect
	SubmissionWait
	SubmissionWrongLevel
	SubmissionError
)

func (so SubmissionOutcome) String() string {
	switch so {
	case SubmissionCorrect:
		return "Correct answer"
	case SubmissionIncorrect:
		return "Incorrect answer"
	case SubmissionWait:
		return "Wait a bit before submitting again"
	case SubmissionWrongLevel:
		return "You are solving the wrong level"
	case SubmissionError:
		return "Error submitting answer"
	default:
		return "Unknown outcome"
	}
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

	resp, err := fetchSiteBody(c, req)
	if err != nil {
		return SubmissionError, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return SubmissionError, fmt.Errorf("Failed to parse HTML: %v", err)
	}

	outcome := doc.Find("article > p").Text()
	if strings.Contains(outcome, "That's the right answer") {
		return SubmissionCorrect, nil
	} else if strings.Contains(outcome, "That's not the right answer") {
		return SubmissionIncorrect, nil
	} else if strings.Contains(outcome, "You gave an answer too recently") {
		return SubmissionWait, nil
	} else if strings.Contains(outcome, "You don't seem to be solving the right level") {
		return SubmissionWrongLevel, nil
	} else {
		return SubmissionError, errors.New("did not match any outcome")
	}
}
