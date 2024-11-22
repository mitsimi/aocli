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

// application/x-www-form-urlencoded
// POST https://adventofcode.com/{}/day/{}/answer
// Payload: {level: number, answer: number}
// url encoded: level=1&answer=21

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
		return "correct"
	case Incorrect:
		return "incorrect"
	case Wait:
		return "wait"
	case WrongLevel:
		return "wrong level"
	case Error:
		return "error"
	default:
		return "unknow outcome"
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
