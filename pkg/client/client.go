package client

import (
    "bytes"
    "crypto/tls"
    "io"
    "net/http"
    "time"
)

type Config struct {
    URL      string
    Method   string
    Headers  map[string]string
    Body     string
    Insecure bool
    Timeout  time.Duration
}

type Response struct {
    Status  string
    Proto   string
    Headers map[string][]string
    Body    []byte
}

func NewClient(cfg Config) *Client {
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

    // Configure client
    client := &http.Client{
        Timeout: c.cfg.Timeout,
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: c.cfg.Insecure,
            },
        },
    }

    // Execute request
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Read response
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return &Response{
        Status:  resp.Status,
        Proto:   resp.Proto,
        Headers: resp.Header,
        Body:    body,
    }, nil
}
