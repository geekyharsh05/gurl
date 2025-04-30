package cmd

import (
    "fmt"
    "os"
    "time"
    "github.com/spf13/cobra"
)

const (
    version = "1.0.0"
)

var rootCmd = &cobra.Command{
    Use:   "gurl",
    Short: "Modern HTTP client inspired by curl",
    Long:  `gurl is a fast, reliable HTTP client with JSON support and intuitive syntax`,
}

// Version command to display version information
var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Display version information",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("gurl version %s\n", version)
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

func init() {
    rootCmd.PersistentFlags().BoolP("insecure", "k", false, "Allow insecure server connections")
    rootCmd.PersistentFlags().Duration("timeout", 30*time.Second, "Request timeout")
    rootCmd.PersistentFlags().Int("max-redirects", 10, "Maximum number of redirects to follow")
    
    rootCmd.AddCommand(versionCmd)
}
