package download

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

// ProgressWriter tracks download progress
type ProgressWriter struct {
	Total            int64
	Downloaded       int64
	LastPercent      int
	LastUpdate       time.Time
	Start            time.Time
	Filename         string
	ShowProgressBar  bool
}

// Downloader handles file downloads with resume capability
type Downloader struct {
	URL              string
	OutputDir        string
	Filename         string
	NumConnections   int
	ShowProgress     bool
	Resume           bool
	Timeout          time.Duration
	Insecure         bool
	FollowRedirect   bool
	MaxRedirects     int
	Retries          int
	RetryDelay       time.Duration
}

// NewDownloader creates a new downloader instance
func NewDownloader(url, outputDir, filename string, showProgress, resume bool, insecure bool, timeout time.Duration, noRedirect bool, maxRedirects int) *Downloader {
	return &Downloader{
		URL:            url,
		OutputDir:      outputDir,
		Filename:       filename,
		NumConnections: 1, // Single connection for now
		ShowProgress:   showProgress,
		Resume:         resume,
		Timeout:        timeout,
		Insecure:       insecure,
		FollowRedirect: !noRedirect,
		MaxRedirects:   maxRedirects,
		Retries:        3, // Default number of retries
		RetryDelay:     2 * time.Second,
	}
}

// Download starts the download process
func (d *Downloader) Download() error {
	var lastErr error
	
	// Try the download with retries
	for attempt := 0; attempt <= d.Retries; attempt++ {
		if attempt > 0 {
			// Wait before retry
			sleepTime := d.RetryDelay * time.Duration(attempt)
			if d.ShowProgress {
				color.Yellow("Retry %d/%d after %v...", attempt, d.Retries, sleepTime)
			}
			time.Sleep(sleepTime)
		}
		
		// Try the download
		err := d.executeDownload()
		if err == nil {
			return nil // Success
		}
		
		lastErr = err
		if d.ShowProgress {
			color.Red("Download attempt %d failed: %v", attempt+1, err)
		}
	}
	
	return fmt.Errorf("download failed after %d attempts: %w", d.Retries+1, lastErr)
}

// executeDownload performs the actual download operation
func (d *Downloader) executeDownload() error {
	// Create output directory if it doesn't exist
	if d.OutputDir != "" {
		if err := os.MkdirAll(d.OutputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Determine filename if not provided
	if d.Filename == "" {
		d.Filename = filepath.Base(d.URL)
		// If URL ends with slash, use "index.html"
		if d.Filename == "" || d.Filename == "." {
			d.Filename = "index.html"
		}
	}

	fullPath := filepath.Join(d.OutputDir, d.Filename)
	
	// Check if file exists for resume
	fileSize := int64(0)
	if d.Resume {
		info, err := os.Stat(fullPath)
		if err == nil {
			fileSize = info.Size()
		}
	}

	// Create HTTP request
	req, err := http.NewRequest(http.MethodGet, d.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set Range header for resume
	if d.Resume && fileSize > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", fileSize))
	}

	// Set user agent
	req.Header.Set("User-Agent", "gurl/1.0")

	// Make HTTP request
	client := &http.Client{
		Timeout: d.Timeout,
	}
	
	// Configure TLS
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: d.Insecure,
		},
	}
	client.Transport = transport
	
	// Configure redirect policy
	if !d.FollowRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else if d.MaxRedirects > 0 {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) >= d.MaxRedirects {
				return http.ErrUseLastResponse
			}
			return nil
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("download request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return fmt.Errorf("bad response status: %s", resp.Status)
	}

	// Check if server supports range requests for resume
	if d.Resume && fileSize > 0 && resp.StatusCode != http.StatusPartialContent {
		// Server doesn't support range requests, start from beginning
		fileSize = 0
		color.Yellow("Resume not supported by server, restarting download")
	}

	// Get content length
	contentLength := resp.ContentLength
	if contentLength == -1 {
		contentLength = 0 // Unknown size
	}
	
	// Total expected size
	totalSize := fileSize + contentLength
	
	// Display download information
	if d.ShowProgress {
		fmt.Printf("Downloading: %s\n", color.CyanString(d.URL))
		fmt.Printf("Saving to: %s\n", color.GreenString(fullPath))
		if contentLength > 0 {
			fmt.Printf("Size: %s\n", formatSize(totalSize))
		}
		if d.Resume && fileSize > 0 {
			fmt.Printf("Resuming from: %s\n", formatSize(fileSize))
		}
		fmt.Println()
	}

	// Open file for writing
	var file *os.File
	if d.Resume && fileSize > 0 {
		file, err = os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY, 0644)
	} else {
		file, err = os.Create(fullPath)
	}
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Create progress writer
	progressWriter := &ProgressWriter{
		Total:           totalSize,
		Downloaded:      fileSize,
		LastUpdate:      time.Now(),
		Start:           time.Now(),
		Filename:        d.Filename,
		ShowProgressBar: d.ShowProgress,
	}

	// Start download
	_, err = io.Copy(io.MultiWriter(file, progressWriter), resp.Body)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	// Print final newline after progress bar
	if d.ShowProgress {
		fmt.Println()
	}

	// Show download summary
	elapsed := time.Since(progressWriter.Start)
	avgSpeed := float64(progressWriter.Downloaded) / elapsed.Seconds()
	
	color.Green("\nDownload complete!")
	fmt.Printf("Total size: %s\n", formatSize(progressWriter.Downloaded))
	fmt.Printf("Time: %s\n", formatDuration(elapsed))
	fmt.Printf("Average speed: %s/s\n", formatSize(int64(avgSpeed)))

	return nil
}

// Write implements io.Writer for tracking download progress
func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.Downloaded += int64(n)
	
	// Update progress not more than once every 100ms
	if pw.ShowProgressBar && time.Since(pw.LastUpdate) > 100*time.Millisecond {
		pw.updateProgress()
	}
	
	return n, nil
}

// updateProgress displays the progress bar
func (pw *ProgressWriter) updateProgress() {
	width := 50 // Progress bar width
	
	// Calculate percentage
	percent := 0
	if pw.Total > 0 {
		percent = int(float64(pw.Downloaded) / float64(pw.Total) * 100)
	}
	
	// Only update if percentage changed or every 0.5 second
	if percent != pw.LastPercent || time.Since(pw.LastUpdate) > 500*time.Millisecond {
		pw.LastPercent = percent
		pw.LastUpdate = time.Now()
		
		// Calculate speed
		elapsed := time.Since(pw.Start)
		speed := float64(pw.Downloaded) / elapsed.Seconds()
		
		// Generate progress bar
		done := int(float64(width) * float64(pw.Downloaded) / float64(pw.Total))
		if pw.Total == 0 {
			// If total is unknown, use a spinner instead
			spinner := []string{"|", "/", "-", "\\"}
			spinnerPos := int(time.Since(pw.Start).Milliseconds()/250) % len(spinner)
			done = 0
			fmt.Printf("\r%s ", spinner[spinnerPos])
		} else {
			fmt.Printf("\r[")
		}
		
		if pw.Total > 0 {
			// Print progress bar
			for i := 0; i < width; i++ {
				if i < done {
					fmt.Print("=")
				} else if i == done {
					fmt.Print(">")
				} else {
					fmt.Print(" ")
				}
			}
			fmt.Printf("] %3d%%", percent)
		}
		
		// Show downloaded / total size and speed
		if pw.Total > 0 {
			fmt.Printf(" %s/%s", formatSize(pw.Downloaded), formatSize(pw.Total))
		} else {
			fmt.Printf(" %s", formatSize(pw.Downloaded))
		}
		
		// Show speed
		fmt.Printf(" %s/s", formatSize(int64(speed)))
		
		// Estimate time remaining if total is known
		if pw.Total > 0 {
			remaining := float64(pw.Total-pw.Downloaded) / speed
			fmt.Printf(" ETA %s", formatDuration(time.Duration(remaining)*time.Second))
		}
		
		// Pad with spaces for clean updates
		fmt.Print("     ")
	}
}

// formatSize returns a human-readable size string
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// formatDuration returns a human-readable duration string
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
} 