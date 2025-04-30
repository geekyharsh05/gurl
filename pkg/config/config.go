package config

import "time"

// Config holds all configuration for an HTTP request
type Config struct {
	URL            string
	Method         string
	Headers        map[string]string
	Body           string
	Insecure       bool
	Timeout        time.Duration
	FollowRedirect bool
	MaxRedirects   int
	WaitTime       time.Duration // Duration to wait before making the request
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