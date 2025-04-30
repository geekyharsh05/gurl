package client

import (
    "bytes"
    "crypto/tls"
    "io"
    "net/http"
    "time"
)

type Config struct {
    URL            string
    Method         string
    Headers        map[string]string
    Body           string
    Insecure       bool
    Timeout        time.Duration
    FollowRedirect bool
    MaxRedirects   int
}

type Response struct {
    Status      string
    StatusCode  int
    Proto       string
    Headers     map[string][]string
    Body        []byte
    TotalTime   time.Duration
    RedirectsFollowed int
    ContentType string // Explicitly store content type for easy access
}

func NewClient(cfg Config) *Client {
    // Set default max redirects if following redirects
    if cfg.FollowRedirect && cfg.MaxRedirects == 0 {
        cfg.MaxRedirects = 10 // Default max redirects
    }
    
    return &Client{cfg: cfg}
}

type Client struct {
    cfg Config
}

func (c *Client) Execute() (*Response, error) {
    req, err := http.NewRequest(c.cfg.Method, c.cfg.URL, bytes.NewBufferString(c.cfg.Body))
    if err != nil {
        return nil, err
    }

    // Set headers
    for k, v := range c.cfg.Headers {
        req.Header.Set(k, v)
    }

    // Set default User-Agent if not specified
    if _, exists := c.cfg.Headers["User-Agent"]; !exists {
        req.Header.Set("User-Agent", "gurl/1.0")
    }

    // Configure client
    client := &http.Client{
        Timeout: c.cfg.Timeout,
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: c.cfg.Insecure,
            },
        },
    }
    
    // Configure redirect policy
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

    // Read response
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    // Extract content type for easier access
    contentType := ""
    if ctHeader := resp.Header.Get("Content-Type"); ctHeader != "" {
        contentType = ctHeader
    }

    return &Response{
        Status:           resp.Status,
        StatusCode:       resp.StatusCode,
        Proto:            resp.Proto,
        Headers:          resp.Header,
        Body:             body,
        TotalTime:        totalTime,
        RedirectsFollowed: redirectsFollowed,
        ContentType:      contentType,
    }, nil
}
