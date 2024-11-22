package aoc

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type Client struct {
	Session string
	Client  *http.Client

	unlockLocation *time.Location
}

// NewClient initializes a new client with a base URL and session token in a jar
func NewClient(token string, options ...Option) *Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Failed to create cookie jar: %v", err)
	}
	cookie := &http.Cookie{
		Name:  "session",
		Value: token,
		Path:  "/",
	}

	u, err := url.Parse(BaseURL)
	if err != nil {
		log.Fatalf("Failed to parse BaseURL: %v", err)
	}

	jar.SetCookies(u, []*http.Cookie{cookie})

	client := &Client{
		Client: &http.Client{
			Jar: jar,
		},

		unlockLocation: time.FixedZone("UTC-5", -5*60*60),
	}

	for _, opt := range options {
		opt(client)
	}

	return client
}

// isDayUnlocked checks if a challenge is unlocked based on the given year and day.
func (c *Client) isDayUnlocked(year int, day int) (bool, error) {
	// Define the release time in UTC-5 for the given year and day
	releaseTime := time.Date(year, time.December, day, 0, 0, 0, 0, c.unlockLocation)

	// Get the current time in UTC
	currentTime := time.Now().In(c.unlockLocation)

	// Check if the current time is after the release time
	return currentTime.After(releaseTime), nil
}

func (c *Client) Request(req *http.Request) ([]byte, error) {
	// Perform the request
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}

type Option func(c *Client)

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.Client.Timeout = timeout
	}
}

func WithTransport(transport http.RoundTripper) Option {
	return func(c *Client) {
		c.Client.Transport = transport
	}
}

func WithRedirectPolicy(f func(req *http.Request, via []*http.Request) error) Option {
	return func(c *Client) {
		c.Client.CheckRedirect = f
	}
}
