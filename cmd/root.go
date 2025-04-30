package cmd

import (
    "os"
    "time"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "gurl",
    Short: "Modern HTTP client inspired by curl",
    Long:  `gurl is a fast, reliable HTTP client with JSON support and intuitive syntax`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

func init() {
    rootCmd.PersistentFlags().BoolP("insecure", "k", false, "Allow insecure server connections")
    rootCmd.PersistentFlags().Duration("timeout", 30*time.Second, "Request timeout")
}
