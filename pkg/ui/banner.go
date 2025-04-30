package ui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	// Color functions for banner elements
	primaryColor   = color.New(color.FgCyan, color.Bold).SprintFunc()
	secondaryColor = color.New(color.FgHiMagenta).SprintFunc()
	accentColor    = color.New(color.FgHiYellow).SprintFunc()
	infoColor      = color.New(color.FgHiWhite).SprintFunc()
)

// DisplayWelcome shows a welcome banner when the CLI starts
func DisplayWelcome(version string) {
	// Clear output
	fmt.Print("\033[H\033[2J")
	
	// Print stylish banner
	banner := []string{
		"",
		"   ██████╗ ██╗   ██╗██████╗ ██╗     ",
		"  ██╔════╝ ██║   ██║██╔══██╗██║     ",
		"  ██║  ███╗██║   ██║██████╔╝██║     ",
		"  ██║   ██║██║   ██║██╔══██╗██║     ",
		"  ╚██████╔╝╚██████╔╝██║  ██║███████╗",
		"   ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚══════╝",
		"",
	}

	// Print the banner with color
	for _, line := range banner {
		fmt.Println(primaryColor(line))
	}

	// Print version and tagline
	fmt.Printf("  %s %s\n", secondaryColor("Version:"), accentColor(version))
	fmt.Printf("  %s %s\n", secondaryColor("Description:"), infoColor("Modern HTTP Client with a colorful interface"))
	fmt.Println()
}

// DisplayEndBanner shows a completion banner
func DisplayEndBanner(requestURL string, statusCode int, duration string) {
	width := 80
	headerLine := strings.Repeat("━", width)
	
	// Success or error color based on status code
	var statusColor func(a ...interface{}) string
	if statusCode >= 200 && statusCode < 300 {
		statusColor = color.New(color.FgGreen, color.Bold).SprintFunc()
	} else if statusCode >= 300 && statusCode < 400 {
		statusColor = color.New(color.FgYellow, color.Bold).SprintFunc()
	} else {
		statusColor = color.New(color.FgRed, color.Bold).SprintFunc()
	}
	
	// Print footer banner
	fmt.Println()
	fmt.Println(secondaryColor(headerLine))
	fmt.Printf(" %s %s\n", primaryColor("Request to:"), infoColor(requestURL))
	fmt.Printf(" %s %s\n", primaryColor("Status:"), statusColor(fmt.Sprintf("%d", statusCode)))
	fmt.Printf(" %s %s\n", primaryColor("Duration:"), accentColor(duration))
	fmt.Println(secondaryColor(headerLine))
	fmt.Println()
}

// DisplaySectionHeader shows a header for a section of output
func DisplaySectionHeader(title string) {
	width := 60
	padding := 2
	titleLength := len(title)
	
	// Calculate line lengths to ensure symmetry
	remainingWidth := width - titleLength - (padding * 2)
	leftSideLength := remainingWidth / 2
	rightSideLength := remainingWidth - leftSideLength
	
	// Ensure minimum lengths
	if leftSideLength < 1 {
		leftSideLength = 1
	}
	if rightSideLength < 1 {
		rightSideLength = 1
	}
	
	// Create the left and right sides
	leftSide := strings.Repeat("━", leftSideLength)
	rightSide := strings.Repeat("━", rightSideLength)
	
	// Build and print the header with consistent spacing
	header := fmt.Sprintf("%s%s%s%s%s", 
		secondaryColor(leftSide),
		strings.Repeat(" ", padding),
		primaryColor(title),
		strings.Repeat(" ", padding),
		secondaryColor(rightSide))
	
	fmt.Println(header)
} 