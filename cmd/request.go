package cmd

import (
    "fmt"
    "os"
    
    "github.com/fatih/color"
    "github.com/spf13/cobra"
    "github.com/geekyharsh05/gurl/pkg/client"
    "github.com/geekyharsh05/gurl/pkg/config"
    "github.com/geekyharsh05/gurl/pkg/printer"
    "github.com/geekyharsh05/gurl/pkg/utils"
)

// requestCmd handles HTTP requests with various options like headers, method, body
var requestCmd = &cobra.Command{
    Use:   "request [URL]",
    Short: "Make HTTP request",
    Args:  cobra.ExactArgs(1),
    Run:   makeRequest,
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        if len(args) == 0 {
            return fmt.Errorf("URL is required")
        }
        
        // Disable colors if requested
        if noColor, _ := rootCmd.PersistentFlags().GetBool("no-color"); noColor {
            color.NoColor = true
        }
        
        return nil
    },
}

// init registers flags for the request command and adds it to the root command
func init() {
    requestCmd.Flags().StringP("method", "X", "GET", "HTTP method")
    requestCmd.Flags().StringP("data", "d", "", "Request body")
    requestCmd.Flags().StringArrayP("header", "H", []string{}, "Custom headers")
    requestCmd.Flags().BoolP("verbose", "v", false, "Show verbose output")
    requestCmd.Flags().StringP("output", "o", "", "Write response to file")
    requestCmd.Flags().BoolP("json", "j", false, "Format response as JSON")
    requestCmd.Flags().BoolP("follow", "L", false, "Follow redirects")
    requestCmd.Flags().Bool("json-request", false, "Set Content-Type to application/json")
    requestCmd.Flags().Bool("form", false, "Set Content-Type to application/x-www-form-urlencoded")
    requestCmd.Flags().Bool("no-pretty", false, "Disable automatic JSON formatting")
    rootCmd.AddCommand(requestCmd)
}

// makeRequest handles the request execution and response processing
func makeRequest(cmd *cobra.Command, args []string) {
    // Create a configuration object
    cfg := buildConfig(cmd, args)

    // Show waiting message if verbose mode is enabled
    if cfg.WaitTime > 0 && cmd.Flag("verbose").Changed {
        color.Yellow("Waiting for %v before making request...", cfg.WaitTime)
    }

    // Execute the request
    resp, err := client.NewClient(cfg).Execute()
    if err != nil {
        color.Red("Error: %v", err)
        os.Exit(1)
    }

    // Create a printer for formatted output
    outputFile, _ := cmd.Flags().GetString("output")
    formatJSON, _ := cmd.Flags().GetBool("json")
    
    // Auto-detect and format JSON responses unless explicitly disabled
    autoDetectJSON := true
    if cmd.Flag("no-pretty").Changed {
        autoDetectJSON = false
    }
    
    printer := printer.NewPrinter(
        cmd.Flag("verbose").Changed,
        formatJSON || autoDetectJSON,
        outputFile,
    )
    
    // Print the response
    printer.Print(resp)
}

// buildConfig creates a config object from command flags
func buildConfig(cmd *cobra.Command, args []string) config.Config {
    method := cmd.Flag("method").Value.String()
    headers, _ := cmd.Flags().GetStringArray("header")
    parsedHeaders := utils.ParseHeaders(headers)
    
    // Set content type shortcuts
    if cmd.Flag("json-request").Changed {
        parsedHeaders["Content-Type"] = "application/json"
    } else if cmd.Flag("form").Changed {
        parsedHeaders["Content-Type"] = "application/x-www-form-urlencoded"
    }
    
    followRedirects, _ := cmd.Flags().GetBool("follow")
    
    // Get timeout from persistent flag
    timeout, _ := rootCmd.PersistentFlags().GetDuration("timeout")
    
    // Get max redirects from persistent flag
    maxRedirects, _ := rootCmd.PersistentFlags().GetInt("max-redirects")
    
    // Get wait time from persistent flag
    waitTime, _ := rootCmd.PersistentFlags().GetDuration("wait-time")
    
    return config.Config{
        URL:            args[0],
        Method:         method,
        Body:           cmd.Flag("data").Value.String(),
        Headers:        parsedHeaders,
        Insecure:       cmd.Flag("insecure").Changed,
        Timeout:        timeout,
        FollowRedirect: followRedirects,
        MaxRedirects:   maxRedirects,
        WaitTime:       waitTime,
    }
}
