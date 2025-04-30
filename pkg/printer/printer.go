package printer

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	
	"github.com/fatih/color"
	"github.com/geekyharsh05/gurl/pkg/response"
	"github.com/geekyharsh05/gurl/pkg/ui"
	"github.com/geekyharsh05/gurl/pkg/utils"
)

// Printer handles the formatting and output of HTTP responses
type Printer struct {
	Verbose       bool
	FormatJSON    bool
	OutputFile    string
}

// NewPrinter creates a new printer with the given options
func NewPrinter(verbose bool, formatJSON bool, outputFile string) *Printer {
	return &Printer{
		Verbose:    verbose,
		FormatJSON: formatJSON,
		OutputFile: outputFile,
	}
}

// Print formats and outputs the response
func (p *Printer) Print(resp *response.Response) error {
	if p.Verbose {
		p.printVerboseInfo(resp)
	}
	
	bodyContent := resp.Body
	
	// Check if response appears to be JSON
	isValidJSON := utils.IsJSON(bodyContent)
	
	// Check content type for JSON
	isJSONContentType := resp.IsJSON() || strings.Contains(strings.ToLower(resp.ContentType), "json")
	
	// Format JSON if it's valid JSON or has JSON content type
	if p.FormatJSON && (isValidJSON || isJSONContentType) {
		prettyJSON, err := utils.PrettyJSON(bodyContent)
		if err == nil {
			bodyContent = prettyJSON
		}
	}
	
	// Write to file if specified
	if p.OutputFile != "" {
		if err := os.WriteFile(p.OutputFile, bodyContent, 0644); err != nil {
			color.Red("Error writing to file: %v", err)
			return err
		}
		color.Green("Response written to %s", p.OutputFile)
		return nil
	}
	
	// Output the response body with color based on content type
	if isValidJSON || isJSONContentType {
		ui.DisplaySectionHeader("RESPONSE BODY")
		
		if p.FormatJSON {
			// Use direct coloring here for important elements
			lines := strings.Split(string(bodyContent), "\n")
			for _, line := range lines {
				// Color the keys blue
				keyPattern := regexp.MustCompile(`"([^"]+)"(\s*):`)
				line = keyPattern.ReplaceAllString(line, color.BlueString("\"$1\"") + "$2:")
				
				// Color string values green
				strPattern := regexp.MustCompile(`:\s*"([^"]*)"`)
				line = strPattern.ReplaceAllString(line, ": " + color.GreenString("\"$1\""))
				
				// Color number values magenta
				numPattern := regexp.MustCompile(`:\s*(-?\d+(\.\d+)?)`)
				line = numPattern.ReplaceAllString(line, ": " + color.MagentaString("$1"))
				
				// Color boolean and null values red
				line = regexp.MustCompile(`:\s*true`).ReplaceAllString(line, ": " + color.RedString("true"))
				line = regexp.MustCompile(`:\s*false`).ReplaceAllString(line, ": " + color.RedString("false"))
				line = regexp.MustCompile(`:\s*null`).ReplaceAllString(line, ": " + color.RedString("null"))
				
				fmt.Println(line)
			}
		} else {
			fmt.Println(string(bodyContent))
		}
	} else {
		fmt.Println(string(bodyContent))
	}
	
	// Show completion banner
	duration := fmt.Sprintf("%v", resp.TotalTime)
	ui.DisplayEndBanner(resp.Request, resp.StatusCode, duration)
	
	return nil
}

// printVerboseInfo prints detailed information about the response
func (p *Printer) printVerboseInfo(resp *response.Response) {
	// Display section headers and formatted content
	ui.DisplaySectionHeader("RESPONSE INFO")
	fmt.Println(ui.FormatStatusCode(resp.Status, resp.StatusCode))
	fmt.Println()
	
	// Headers section
	ui.DisplaySectionHeader("HEADERS")
	headers := ui.FormatHeaders(resp.Headers)
	for _, line := range headers {
		fmt.Println(line)
	}
	fmt.Println()
	
	// Timing section
	ui.DisplaySectionHeader("TIMING")
	timing := ui.FormatTiming(resp.WaitTime, resp.TotalTime, resp.RedirectsFollowed)
	for _, line := range timing {
		fmt.Println(line)
	}
	fmt.Println()
} 