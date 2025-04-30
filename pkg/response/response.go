package response

import "time"

// Response represents an HTTP response with all relevant information
// including status, headers, body, timing, and redirect information
type Response struct {
	Status           string              // Status text (e.g., "200 OK")
	StatusCode       int                 // HTTP status code
	Proto            string              // Protocol version (e.g., "HTTP/1.1")
	Headers          map[string][]string // Response headers
	Body             []byte              // Response body
	TotalTime        time.Duration       // Total time taken for the request
	WaitTime         time.Duration       // The duration waited before the request
	RedirectsFollowed int                // Number of redirects followed
	ContentType      string              // Explicitly store content type for easy access
	Request          string              // The URL that was requested
}

// IsJSON returns true if the response appears to be JSON based on the content type
// Used to determine if responses should be formatted as JSON
func (r *Response) IsJSON() bool {
	if r.ContentType == "" {
		return false
	}
	
	contentType := r.ContentType
	return contentType == "application/json" || 
	       contentType == "application/ld+json" ||
	       contentType == "application/json; charset=utf-8"
} 