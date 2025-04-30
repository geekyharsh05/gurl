package response

import "time"

// Response represents an HTTP response
type Response struct {
	Status           string
	StatusCode       int
	Proto            string
	Headers          map[string][]string
	Body             []byte
	TotalTime        time.Duration
	WaitTime         time.Duration // The duration waited before the request
	RedirectsFollowed int
	ContentType      string // Explicitly store content type for easy access
	Request          string // The URL that was requested
}

// IsJSON returns true if the response appears to be JSON based on the content type
func (r *Response) IsJSON() bool {
	if r.ContentType == "" {
		return false
	}
	
	contentType := r.ContentType
	return contentType == "application/json" || 
	       contentType == "application/ld+json" ||
	       contentType == "application/json; charset=utf-8"
} 