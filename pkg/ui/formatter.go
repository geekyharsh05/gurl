package ui

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/geekyharsh05/gurl/pkg/response"
)

var (
	// Header styling
	headerKey   = color.New(color.FgBlue, color.Bold).SprintFunc()
	headerValue = color.New(color.FgWhite).SprintFunc()
	
	// Timing styling
	timingLabel = color.New(color.FgYellow, color.Bold).SprintFunc()
	timingValue = color.New(color.FgHiYellow).SprintFunc()
	
	// JSON styling
	jsonKey     = color.New(color.FgHiBlue).SprintFunc()
	jsonString  = color.New(color.FgGreen).SprintFunc()
	jsonNumber  = color.New(color.FgHiMagenta).SprintFunc()
	jsonBoolean = color.New(color.FgRed).SprintFunc()
	jsonNull    = color.New(color.FgRed, color.Bold).SprintFunc()
	
	// Status code styling
	statusSuccess = color.New(color.FgGreen, color.Bold).SprintFunc()
	statusRedirect = color.New(color.FgYellow, color.Bold).SprintFunc()
	statusClientError = color.New(color.FgRed, color.Bold).SprintFunc()
	statusServerError = color.New(color.FgHiRed, color.Bold).SprintFunc()
)

// FormatHeaders returns headers with colorful formatting
func FormatHeaders(headers map[string][]string) []string {
	var lines []string
	
	for k, v := range headers {
		if len(v) > 0 {
			line := fmt.Sprintf("%s: %s", headerKey(k), headerValue(v[0]))
			lines = append(lines, line)
		}
	}
	
	return lines
}

// FormatTiming returns timing information with colorful formatting
func FormatTiming(waitTime, requestTime time.Duration, redirects int) []string {
	var lines []string
	
	if waitTime > 0 {
		line := fmt.Sprintf("%s: %s", timingLabel("Wait time"), timingValue(waitTime))
		lines = append(lines, line)
	}
	
	lines = append(lines, fmt.Sprintf("%s: %s", timingLabel("Request time"), timingValue(requestTime)))
	
	if redirects > 0 {
		line := fmt.Sprintf("%s: %s", timingLabel("Redirects"), timingValue(redirects))
		lines = append(lines, line)
	}
	
	return lines
}

// FormatStatusCode returns a color-coded status code string
func FormatStatusCode(status string, code int) string {
	var formatter func(a ...interface{}) string
	
	switch {
	case code >= 200 && code < 300:
		formatter = statusSuccess
	case code >= 300 && code < 400:
		formatter = statusRedirect
	case code >= 400 && code < 500:
		formatter = statusClientError
	case code >= 500:
		formatter = statusServerError
	default:
		formatter = color.New(color.FgWhite).SprintFunc()
	}
	
	return formatter(status)
}

// FormatJSONSample colorizes a small sample of JSON for display
func FormatJSONSample(body []byte, maxLen int) string {
	sample := string(body)
	if len(sample) > maxLen {
		sample = sample[:maxLen] + "..."
	}
	
	// Simple colorization for JSON preview
	sample = highlightJSONSyntax(sample)
	
	return sample
}

// FormatJSON returns a colorized version of the provided JSON
func FormatJSON(jsonData []byte) string {
	// Return the fully colorized JSON
	return highlightJSONSyntax(string(jsonData))
}

// highlightJSONSyntax adds color to JSON elements
func highlightJSONSyntax(json string) string {
	// Create a reader for each line of the JSON
	lines := strings.Split(json, "\n")
	coloredLines := make([]string, len(lines))
	
	// Process each line individually for more reliable highlighting
	for i, line := range lines {
		// Highlight keys - matches "key": pattern
		re := regexp.MustCompile(`"([^"]+)"\s*:`)
		line = re.ReplaceAllStringFunc(line, func(match string) string {
			parts := re.FindStringSubmatch(match)
			if len(parts) > 1 {
				return jsonKey("\""+parts[1]+"\"") + ":"
			}
			return match
		})
		
		// Highlight string values - matches : "value" pattern
		re = regexp.MustCompile(`:\s*"([^"]*)"`)
		line = re.ReplaceAllStringFunc(line, func(match string) string {
			parts := re.FindStringSubmatch(match)
			if len(parts) > 1 {
				return ": " + jsonString("\""+parts[1]+"\"")
			}
			return match
		})
		
		// Highlight numbers - matches : 123 or : 123.456 pattern
		re = regexp.MustCompile(`:\s*(-?\d+(\.\d+)?)`)
		line = re.ReplaceAllStringFunc(line, func(match string) string {
			parts := re.FindStringSubmatch(match)
			if len(parts) > 1 {
				return ": " + jsonNumber(parts[1])
			}
			return match
		})
		
		// Highlight booleans and null
		line = regexp.MustCompile(`:\s*true`).ReplaceAllString(line, ": "+jsonBoolean("true"))
		line = regexp.MustCompile(`:\s*false`).ReplaceAllString(line, ": "+jsonBoolean("false"))
		line = regexp.MustCompile(`:\s*null`).ReplaceAllString(line, ": "+jsonNull("null"))
		
		coloredLines[i] = line
	}
	
	return strings.Join(coloredLines, "\n")
}

// GetResponsePreview generates a colorful preview of the response
func GetResponsePreview(resp *response.Response) string {
	status := FormatStatusCode(resp.Status, resp.StatusCode)
	
	var preview string
	if resp.IsJSON() {
		preview = FormatJSONSample(resp.Body, 100) // Show first 100 chars of JSON
	} else if len(resp.Body) > 0 {
		if len(resp.Body) > 100 {
			preview = string(resp.Body[:100]) + "..."
		} else {
			preview = string(resp.Body)
		}
	} else {
		preview = "[No body]"
	}
	
	result := fmt.Sprintf("%s\nContent-Type: %s\nSize: %d bytes\n\nPreview: %s",
		status,
		color.HiWhiteString(resp.ContentType),
		len(resp.Body),
		preview)
	
	return result
} 