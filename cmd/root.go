package cmd

import (
    "fmt"
    "os"
    "time"
    "github.com/spf13/cobra"
    "github.com/geekyharsh05/gurl/pkg/ui"
)

const (
    version = "1.0.0"
)

// rootCmd is the base command for the gurl CLI application
var rootCmd = &cobra.Command{
    Use:   "gurl",
    Short: "Modern HTTP client inspired by curl",
    Long:  `gurl is a fast, reliable HTTP client with JSON support and intuitive syntax`,
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        // Only show welcome banner for the root command and when not disabled
        if cmd.Name() == "gurl" && os.Getenv("GURL_NO_BANNER") != "1" {
            ui.DisplayWelcome(version)
        }
    },
    Run: func(cmd *cobra.Command, args []string) {
        // Check if version flag is set
        if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
            fmt.Printf("gurl version %s\n", version)
            return
        }

        // If no subcommand is provided, display the welcome banner and help
        if os.Getenv("GURL_NO_BANNER") != "1" {
            ui.DisplayWelcome(version)
        }
        cmd.Help()
    },
}

// Version command to display version information
var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Display version information",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("gurl version %s\n", version)
    },
}

// Execute runs the root command and handles any errors
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

// init sets up global flags and adds subcommands
func init() {
    rootCmd.Flags().BoolP("version", "V", false, "Display version information")
    rootCmd.PersistentFlags().BoolP("insecure", "k", false, "Allow insecure server connections")
    rootCmd.PersistentFlags().Duration("timeout", 30*time.Second, "Request timeout")
    rootCmd.PersistentFlags().Int("max-redirects", 10, "Maximum number of redirects to follow")
    rootCmd.PersistentFlags().Duration("wait-time", 0, "Wait for the specified duration before making the request")
    rootCmd.PersistentFlags().Bool("no-color", false, "Disable colorized output")
    
    rootCmd.AddCommand(versionCmd)
}
