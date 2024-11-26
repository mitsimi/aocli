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
	http.Client
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

	client := &Client{}
	client.Jar = jar

	for _, opt := range options {
		opt(client)
	}

	return client
}

// isDayUnlocked checks if a challenge is unlocked based on the given year and day.
func IsDayUnlocked(year int, day int) (bool, error) {
	var unlockLocation = time.FixedZone("UTC-5", -5*60*60)

	// Define the release time in UTC-5 for the given year and day
	releaseTime := time.Date(year, time.December, day, 0, 0, 0, 0, unlockLocation)

	// Get the current time in UTC
	currentTime := time.Now().In(unlockLocation)

	// Check if the current time is after the release time
	return currentTime.After(releaseTime), nil
}

type RequestError struct {
	StatusCode int
	Err        error
}

func (e RequestError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("unexpected status code: %d - %v", e.StatusCode, e.Err)
	}
	return fmt.Sprintf("unexpected status code: %d", e.StatusCode)
}

func (c *Client) Request(req *http.Request) ([]byte, error) {
	// Perform the request
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, RequestError{resp.StatusCode, nil}
	}

	// Read the response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return data, nil
}

type Option func(c *Client)

// WithTimeout sets the timeout duration for the client
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.Client.Timeout = timeout
	}
}

// WithTransport sets a custom transport for the client
func WithTransport(transport http.RoundTripper) Option {
	return func(c *Client) {
		c.Client.Transport = transport
	}
}

// WithRedirectPolicy sets a custom redirecct policy for the client
func WithRedirectPolicy(f func(req *http.Request, via []*http.Request) error) Option {
	return func(c *Client) {
		c.Client.CheckRedirect = f
	}
}

type DebugTransport struct {
	Transport http.RoundTripper
}

func NewDebugTransport() *DebugTransport {
	return &DebugTransport{
		Transport: http.DefaultTransport,
	}
}

func (d *DebugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// log url and body of the request
	log.Printf("Request: %s %s", req.Method, req.URL)
	if req.Body != nil {
		body, _ := io.ReadAll(req.Body)
		log.Printf("Body: %s", body)
	}
	return d.Transport.RoundTrip(req)
}
