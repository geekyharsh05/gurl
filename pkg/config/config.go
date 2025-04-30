package config

import "time"

// Config holds all configuration for an HTTP request
type Config struct {
	URL            string            // Target URL for the request
	Method         string            // HTTP method (GET, POST, etc.)
	Headers        map[string]string // HTTP headers to send
	Body           string            // Request body content
	Insecure       bool              // Whether to skip TLS certificate verification
	Timeout        time.Duration     // Request timeout duration
	FollowRedirect bool              // Whether to follow HTTP redirects
	MaxRedirects   int               // Maximum number of redirects to follow
	WaitTime       time.Duration     // Duration to wait before making the request
}

// NewDefaultConfig creates a Config with reasonable defaults
func NewDefaultConfig() Config {
	return Config{
		Method:       "GET",
		Headers:      make(map[string]string),
		Timeout:      30 * time.Second,
		MaxRedirects: 10,
	}
} 