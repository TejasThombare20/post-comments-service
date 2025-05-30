package utils

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

// HTMLSanitizer handles HTML content sanitization and processing
type HTMLSanitizer struct {
	policy *bluemonday.Policy
}

// NewHTMLSanitizer creates a new HTML sanitizer instance
func NewHTMLSanitizer() *HTMLSanitizer {
	// Create a policy that allows common rich text formatting
	policy := bluemonday.NewPolicy()

	// Allow basic text formatting
	policy.AllowElements("p", "br", "span", "div")

	// Allow text styling
	policy.AllowElements("strong", "b", "em", "i", "u", "s", "del", "ins")

	// Allow headings
	policy.AllowElements("h1", "h2", "h3", "h4", "h5", "h6")

	// Allow lists
	policy.AllowElements("ul", "ol", "li")

	// Allow links with href attribute
	policy.AllowAttrs("href").OnElements("a")
	policy.AllowElements("a")

	// Allow basic styling attributes
	policy.AllowAttrs("style").OnElements("span", "div", "p")

	// Allow class attributes for styling
	policy.AllowAttrs("class").OnElements("span", "div", "p", "strong", "em", "b", "i")

	// Allow blockquotes
	policy.AllowElements("blockquote")

	// Allow code elements
	policy.AllowElements("code", "pre")

	return &HTMLSanitizer{
		policy: policy,
	}
}

// SanitizeHTML sanitizes HTML content using the bluemonday policy
func (h *HTMLSanitizer) SanitizeHTML(content string) string {
	return h.policy.Sanitize(content)
}

// IsHTMLContent checks if the content appears to be HTML
func (h *HTMLSanitizer) IsHTMLContent(content string) bool {
	// Check for common HTML tags
	htmlTagPattern := regexp.MustCompile(`<[^>]+>`)
	return htmlTagPattern.MatchString(content)
}

// ValidateHTMLContent validates that HTML content is safe and well-formed
func (h *HTMLSanitizer) ValidateHTMLContent(content string) error {
	if content == "" {
		return nil
	}

	// Check for potentially dangerous patterns
	dangerousPatterns := []string{
		`<script`,
		`javascript:`,
		`onload=`,
		`onerror=`,
		`onclick=`,
		`onmouseover=`,
		`<iframe`,
		`<object`,
		`<embed`,
		`<form`,
	}

	lowerContent := strings.ToLower(content)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerContent, pattern) {
			return fmt.Errorf("potentially dangerous HTML content detected: %s", pattern)
		}
	}

	return nil
}

// ConvertPlainTextToHTML converts plain text to HTML by wrapping it in a span tag
func (h *HTMLSanitizer) ConvertPlainTextToHTML(content string) string {
	// Escape HTML entities in plain text
	escaped := html.EscapeString(content)

	// Convert line breaks to <br> tags
	escaped = strings.ReplaceAll(escaped, "\n", "<br>")

	// Wrap in span tag
	return "<span>" + escaped + "</span>"
}

// ProcessCommentContent processes comment content based on whether it's HTML or plain text
func (h *HTMLSanitizer) ProcessCommentContent(content string) string {
	if content == "" {
		return ""
	}

	// Check if content is already HTML
	if h.IsHTMLContent(content) {
		// Sanitize the HTML content
		return h.SanitizeHTML(content)
	}

	// Convert plain text to HTML and then sanitize
	htmlContent := h.ConvertPlainTextToHTML(content)
	return h.SanitizeHTML(htmlContent)
}

// StripHTMLTags removes all HTML tags and returns plain text (useful for previews)
func (h *HTMLSanitizer) StripHTMLTags(content string) string {
	// Use bluemonday's StripTagsPolicy to remove all HTML tags
	stripPolicy := bluemonday.StripTagsPolicy()
	return stripPolicy.Sanitize(content)
}
