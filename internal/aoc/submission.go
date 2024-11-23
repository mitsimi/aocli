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

type SubmissionOutcome int

const (
	Correct SubmissionOutcome = iota
	Incorrect
	Wait
	WrongLevel
	Error
)

func (so SubmissionOutcome) String() string {
	switch so {
	case Correct:
		return "Correct answer"
	case Incorrect:
		return "Incorrect answer"
	case Wait:
		return "Wait a bit before submitting again"
	case WrongLevel:
		return "You are solving the wrong level"
	case Error:
		return "Error submitting answer"
	default:
		return "Unknown outcome"
	}
}

func (c *Client) Submit(level Level, year, day int, answer int) (SubmissionOutcome, error) {
	data := url.Values{}
	data.Set("level", fmt.Sprintf("%d", level))
	data.Set("answer", fmt.Sprintf("%d", answer))

	req, err := http.NewRequest("POST", SubmitURL(year, day), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return Error, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := fetchSiteBody(c, req)
	if err != nil {
		return Error, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return Error, fmt.Errorf("Failed to parse HTML: %v", err)
	}

	outcome := doc.Find("article > p").Text()
	if strings.Contains(outcome, "That's the right answer") {
		return Correct, nil
	} else if strings.Contains(outcome, "That's not the right answer") {
		return Incorrect, nil
	} else if strings.Contains(outcome, "You gave an answer too recently") {
		return Wait, nil
	} else if strings.Contains(outcome, "You don't seem to be solving the right level") {
		return WrongLevel, nil
	} else {
		return Error, errors.New("did not match any outcome")
	}
}
