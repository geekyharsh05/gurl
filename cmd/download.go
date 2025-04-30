package cmd

import (
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/geekyharsh05/gurl/pkg/download"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download [URL]",
	Short: "Download files with wget-like functionality",
	Long: `Download files with progress bar, resume capability and more.
Examples:
  gurl download https://example.com/file.zip 
  gurl download -o myfile.zip https://example.com/file.zip
  gurl download -c -P ./downloads/ https://example.com/file.zip`,
	Args: cobra.ExactArgs(1),
	Run:  runDownload,
}

func init() {
	downloadCmd.Flags().StringP("output", "o", "", "Save file with the specified name")
	downloadCmd.Flags().StringP("directory", "P", "", "Save files to the specified directory")
	downloadCmd.Flags().BoolP("continue", "c", false, "Resume getting a partially-downloaded file")
	downloadCmd.Flags().BoolP("quiet", "q", false, "Quiet mode - don't show progress bar")
	downloadCmd.Flags().BoolP("insecure", "k", false, "Allow insecure server connections")
	downloadCmd.Flags().DurationP("timeout", "t", 30*time.Second, "Set timeout for download")
	downloadCmd.Flags().BoolP("no-redirect", "n", false, "Don't follow redirects")
	downloadCmd.Flags().Int("max-redirects", 10, "Maximum number of redirects to follow")
	downloadCmd.Flags().IntP("tries", "r", 3, "Number of retry attempts (0 for no retries)")
	downloadCmd.Flags().Duration("retry-delay", 2*time.Second, "Delay between retries")

	rootCmd.AddCommand(downloadCmd)
}

func runDownload(cmd *cobra.Command, args []string) {
	url := args[0]
	output, _ := cmd.Flags().GetString("output")
	directory, _ := cmd.Flags().GetString("directory")
	continueDownload, _ := cmd.Flags().GetBool("continue")
	quiet, _ := cmd.Flags().GetBool("quiet")
	insecure, _ := cmd.Flags().GetBool("insecure")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	noRedirect, _ := cmd.Flags().GetBool("no-redirect")
	maxRedirects, _ := cmd.Flags().GetInt("max-redirects")
	tries, _ := cmd.Flags().GetInt("tries")
	retryDelay, _ := cmd.Flags().GetDuration("retry-delay")

	// Convert directory to absolute path if not empty
	if directory != "" {
		absDir, err := filepath.Abs(directory)
		if err != nil {
			color.Red("Error resolving directory path: %v", err)
			os.Exit(1)
		}
		directory = absDir
	}

	// Create downloader with retry settings
	downloader := &download.Downloader{
		URL:            url,
		OutputDir:      directory,
		Filename:       output,
		NumConnections: 1,
		ShowProgress:   !quiet,
		Resume:         continueDownload,
		Timeout:        timeout,
		Insecure:       insecure,
		FollowRedirect: !noRedirect,
		MaxRedirects:   maxRedirects,
		Retries:        tries,
		RetryDelay:     retryDelay,
	}

	// Execute download
	if err := downloader.Download(); err != nil {
		color.Red("Download failed: %v", err)
		os.Exit(1)
	}
} 