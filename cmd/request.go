package cmd

import (
    "fmt"
    "os"
    "strings"
    "time"
    
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
    rootCmd.AddCommand(requestCmd)
}

func makeRequest(cmd *cobra.Command, args []string) {
    cfg := client.Config{
        URL:     args[0],
        Method:  cmd.Flag("method").Value.String(),
        Body:    cmd.Flag("data").Value.String(),
        Headers: parseHeaders(cmd.Flags().GetStringArray("header")),
        Insecure: cmd.Flag("insecure").Changed,
        Timeout: 30 * time.Second,
    }

    resp, err := client.NewClient(cfg).Execute()
    if err != nil {
        color.Red("Error: %v", err)
        os.Exit(1)
    }

    printResponse(resp, cmd.Flag("verbose").Changed)
}

func printResponse(resp *client.Response, verbose bool) {
    if verbose {
        color.Cyan("%s %s", resp.Proto, resp.Status)
        for k, v := range resp.Headers {
            color.Cyan("%s: %s", k, v[0])
        }
        fmt.Println()
    }
    
    color.White(string(resp.Body))
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
