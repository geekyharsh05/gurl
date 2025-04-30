package cmd

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
    
    "github.com/fatih/color"
    "github.com/spf13/cobra"
    "github.com/geekyharsh05/gurl/pkg/client"
)

var requestCmd = &cobra.Command{
    Use:   "request [URL]",
    Short: "Make HTTP request",
    Args:  cobra.ExactArgs(1),
    Run:   makeRequest,
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        if len(args) == 0 {
            return fmt.Errorf("URL is required")
        }
        return nil
    },
}

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

func makeRequest(cmd *cobra.Command, args []string) {
    method := cmd.Flag("method").Value.String()
    headers, _ := cmd.Flags().GetStringArray("header")
    parsedHeaders := parseHeaders(headers, nil)
    
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
    
    cfg := client.Config{
        URL:            args[0],
        Method:         method,
        Body:           cmd.Flag("data").Value.String(),
        Headers:        parsedHeaders,
        Insecure:       cmd.Flag("insecure").Changed,
        Timeout:        timeout,
        FollowRedirect: followRedirects,
        MaxRedirects:   maxRedirects,
    }

    resp, err := client.NewClient(cfg).Execute()
    if err != nil {
        color.Red("Error: %v", err)
        os.Exit(1)
    }

    outputFile, _ := cmd.Flags().GetString("output")
    formatJSON, _ := cmd.Flags().GetBool("json")
    
    // Auto-detect and format JSON responses unless explicitly disabled
    autoDetectJSON := true
    if cmd.Flag("no-pretty").Changed {
        autoDetectJSON = false
    }
    
    printResponse(resp, cmd.Flag("verbose").Changed, outputFile, formatJSON || autoDetectJSON)
}

func printResponse(resp *client.Response, verbose bool, outputFile string, formatJSON bool) {
    if verbose {
        color.Cyan("%s %s", resp.Proto, resp.Status)
        for k, v := range resp.Headers {
            color.Cyan("%s: %s", k, v[0])
        }
        
        // Print timing and redirect info
        color.Yellow("Time: %v", resp.TotalTime)
        if resp.RedirectsFollowed > 0 {
            color.Yellow("Redirects: %d", resp.RedirectsFollowed)
        }
        
        fmt.Println()
    }
    
    bodyContent := resp.Body
    
    // Check if response appears to be JSON
    isValidJSON := isJSON(bodyContent)
    
    // Check content type for JSON
    isJSONContentType := false
    if resp.ContentType != "" {
        contentType := strings.ToLower(resp.ContentType)
        if strings.Contains(contentType, "application/json") || 
           strings.Contains(contentType, "application/ld+json") {
            isJSONContentType = true
        }
    }
    
    // Format JSON if it's valid JSON or has JSON content type
    if formatJSON && (isValidJSON || isJSONContentType) {
        var obj interface{}
        if err := json.Unmarshal(bodyContent, &obj); err == nil {
            // Convert back to JSON with pretty formatting
            formattedJSON, err := json.MarshalIndent(obj, "", "  ")
            if err == nil {
                bodyContent = formattedJSON
            }
        }
    }
    
    // Write to file if specified
    if outputFile != "" {
        if err := os.WriteFile(outputFile, bodyContent, 0644); err != nil {
            color.Red("Error writing to file: %v", err)
        } else {
            color.Green("Response written to %s", outputFile)
        }
        return
    }
    
    // Output the response body with color based on content type
    if isValidJSON || isJSONContentType {
        fmt.Println(color.CyanString(string(bodyContent)))
    } else {
        fmt.Println(string(bodyContent))
    }
}

// Helper to parse headers
func parseHeaders(headers []string, err error) map[string]string {
    if err != nil {
        return make(map[string]string)
    }
    
    h := make(map[string]string)
    for _, header := range headers {
        split := strings.SplitN(header, ":", 2)
        if len(split) == 2 {
            h[strings.TrimSpace(split[0])] = strings.TrimSpace(split[1])
        }
    }
    return h
}

// Check if byte slice is valid JSON
func isJSON(data []byte) bool {
    var js json.RawMessage
    return json.Unmarshal(data, &js) == nil
}

// Custom flag type for array values
type ArrayFlags []string

func (a *ArrayFlags) String() string { 
    return strings.Join(*a, ",") 
}

func (a *ArrayFlags) Set(value string) error {
    *a = append(*a, value)
    return nil
}

func (a *ArrayFlags) Type() string {
    return "stringArray"
}
