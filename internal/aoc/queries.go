package aoc

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// getDescription fetches the and parses the html content
func (c *Client) GetDescription(year, day int) (HTMLContent, error) {
	// Create the request
	req, err := http.NewRequest("GET", DayURL(year, day), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Get site content
	resp, err := c.Request(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to parse HTML: %v", err)
	}

	// Extract the <main> tag
	mainContent, _ := doc.Find("main").Html()
	if mainContent == "" {
		return "", fmt.Errorf("No <main> tag found")
	}

	return HTMLContent(mainContent), nil
}

func (c *Client) GetExamples(year, day int) ([]string, error) {
	// Create the request
	req, err := http.NewRequest("GET", DayURL(year, day), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Get site content
	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse HTML: %v", err)
	}

	return parseExamples(doc), nil
}

func (c *Client) GetInput(year, day int) (string, error) {
	// Create the request
	req, err := http.NewRequest("GET", InputURL(year, day), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Get site content
	resp, err := c.Request(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to parse HTML: %v", err)
	}

	data, _ := doc.Find("body").Html()
	if data == "" {
		return "", fmt.Errorf("No <body> tag found")
	}

	return data, nil
}

// parses each code block after a p element if it contains the word "example"
func parseExamples(doc *goquery.Document) []string {
	dupe := make(map[string]struct{})
	// Find the desired <code> tags
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		// Check if the paragraph contains "For example:"
		if strings.Contains(strings.ToLower(s.Text()), "example") {
			// Get the next <pre> sibling after this paragraph
			preTag := s.NextFiltered("pre")
			if preTag.Length() > 0 {
				// Extract the code content inside <pre><code>
				dupe[preTag.Find("code").Text()] = struct{}{}
			}
		}
	})

	examples := make([]string, 0, 2)
	for k := range dupe {
		examples = append(examples, k)
	}
	return examples
}

type HTMLContent string
type Markdown = string

// convert html to markdown
func (c HTMLContent) ToMarkdown(year int) (Markdown, error) {
	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),
			commonmark.NewCommonmarkPlugin(),
		),
	)

	// custom render function to render <em> as bold instead of italic
	renderEmToBold := func(ctx converter.Context, w converter.Writer, node *html.Node) converter.RenderStatus {
		w.WriteString(" **")
		ctx.RenderChildNodes(ctx, w, node)
		w.WriteString("** ")

		return converter.RenderSuccess
	}

	conv.Register.RendererFor("em", converter.TagTypeInline, renderEmToBold, converter.PriorityEarly)

	markdown, err := conv.ConvertString(string(c), converter.WithDomain(fmt.Sprintf("%s/%d/day/", BaseURL, year)))
	if err != nil {
		return "", fmt.Errorf("Failed to convert HTML to Markdown: %v", err)
	}

	return markdown, nil
}
