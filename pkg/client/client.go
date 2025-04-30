package client

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"time"
	
	"github.com/geekyharsh05/gurl/pkg/config"
	"github.com/geekyharsh05/gurl/pkg/response"
)

type Client struct {
	cfg config.Config
}

// NewClient creates a new Client with the given configuration
func NewClient(cfg config.Config) *Client {
	// Set default max redirects if following redirects
	if cfg.FollowRedirect && cfg.MaxRedirects == 0 {
		cfg.MaxRedirects = 10 // Default max redirects
	}
	
	return &Client{cfg: cfg}
}

// Execute performs the HTTP request according to the configuration
// and returns a Response object with the results
func (c *Client) Execute() (*response.Response, error) {
	// Wait for the specified duration before making the request
	if c.cfg.WaitTime > 0 {
		time.Sleep(c.cfg.WaitTime)
	}
	
	// Create the HTTP request with method, URL and body
	req, err := http.NewRequest(c.cfg.Method, c.cfg.URL, bytes.NewBufferString(c.cfg.Body))
	if err != nil {
		return nil, err
	}

	// Set headers from configuration
	for k, v := range c.cfg.Headers {
		req.Header.Set(k, v)
	}

	// Set default User-Agent if not specified
	if _, exists := c.cfg.Headers["User-Agent"]; !exists {
		req.Header.Set("User-Agent", "gurl/1.0")
	}

	// Configure HTTP client with timeout and TLS settings
	client := &http.Client{
		Timeout: c.cfg.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: c.cfg.Insecure,
			},
		},
	}
	
	// Configure redirect policy based on settings
	redirectsFollowed := 0
	if !c.cfg.FollowRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else if c.cfg.MaxRedirects > 0 {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) >= c.cfg.MaxRedirects {
				return http.ErrUseLastResponse
			}
			redirectsFollowed = len(via)
			return nil
		}
	}

	// Execute request and track time
	startTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	totalTime := time.Since(startTime)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	// Extract content type for easier access
	contentType := ""
	if ctHeader := resp.Header.Get("Content-Type"); ctHeader != "" {
		contentType = ctHeader
	}

	// Return a structured response with all relevant information
	return &response.Response{
		Status:           resp.Status,
		StatusCode:       resp.StatusCode,
		Proto:            resp.Proto,
		Headers:          resp.Header,
		Body:             body,
		TotalTime:        totalTime,
		WaitTime:         c.cfg.WaitTime,
		RedirectsFollowed: redirectsFollowed,
		ContentType:      contentType,
		Request:          c.cfg.URL,
	}, nil
}
