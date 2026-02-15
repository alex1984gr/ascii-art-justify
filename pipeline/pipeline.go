// Package pipeline orchestrates the ASCII art generation process
package pipeline

import (
	"fmt" // Formatted I/O functions for printing error messages
	"io"  // Basic I/O interfaces like Writer
	"os"  // Operating system functions for stderr output
	"strings" // String manipulation functions like ReplaceAll
)

// Run executes the complete ASCII art pipeline from input to output.
// It takes command-line arguments and an output writer, returns 0 on success or 1 on error.
func Run(args []string, stdout io.Writer) int {
	// Parse command-line arguments into a configuration struct
	cfg, err := parseArgs(args)
	// If argument parsing failed, print usage message and exit with error code
	if err != nil {
		fmt.Fprintln(os.Stderr, usageMessage)
		return 1
	}

	// Replace literal \n sequences with actual newline characters
	cfg.input = strings.ReplaceAll(cfg.input, "\\n", "\n")

	// Validate the input text (check for empty, too long, or invalid characters)
	if err := ValidateInput(cfg.input); err != nil {
		// Print validation error to stderr and exit with error code
		fmt.Fprintln(os.Stderr, "invalid input:", err)
		return 1
	}

	// Load the specified banner font file (standard, shadow, or thinkertoy)
	banner, err := LoadBanner(cfg.font)
	if err != nil {
		// Print banner loading error to stderr and exit with error code
		fmt.Fprintln(os.Stderr, "failed loading banner:", err)
		return 1
	}

	// Get the terminal width for alignment calculations
	terminalWidth := getTerminalWidth()
	// If justify alignment is requested, modify input to add extra spaces between words
	if cfg.align == "justify" {
		// Justify must be done before rendering to adjust the input text itself
		cfg.input = justifyInputForTerminal(cfg.input, banner, terminalWidth)
	}

	// Split input into individual character tokens
	tokens := Tokenize(cfg.input)
	// Render tokens into ASCII art lines using the loaded banner
	lines := RenderLines(tokens, banner)
	// Apply horizontal alignment (left, center, right) to the rendered lines
	lines = applyAlignment(lines, cfg.align, terminalWidth)

	// If color option is specified, apply ANSI color codes to the lines
	if cfg.colorName != "" {
		// Color either the entire output or just the specified substring
		lines, err = ColorLinesWithBanner(lines, cfg.colorName, cfg.substring, banner)
		if err != nil {
			// Print color error to stderr and exit with error code
			fmt.Fprintln(os.Stderr, "color error:", err)
			return 1
		}
	}

	// Default output writer is stdout
	var w io.Writer = stdout
	// If output file is specified, create the file and use it as writer
	if cfg.outFile != "" {
		// Create or overwrite the output file
		f, err := os.Create(cfg.outFile)
		if err != nil {
			// Print file creation error to stderr and exit with error code
			fmt.Fprintln(os.Stderr, "failed creating output file:", err)
			return 1
		}
		// Ensure file is closed when function returns
		defer f.Close()
		// Use the file as output writer
		w = f
	}

	// Write the final ASCII art lines to the output writer
	if err := WriteOutput(lines, w); err != nil {
		// Print write error to stderr and exit with error code
		fmt.Fprintln(os.Stderr, "failed writing output:", err)
		return 1
	}

	// Return success code
	return 0
}
