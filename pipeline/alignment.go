// Package pipeline contains text alignment utilities for ASCII art output.
package pipeline

import (
	"os"      // Operating system functions for environment variables
	"strconv" // String conversion functions (string to integer)
	"strings" // String manipulation functions like Repeat, Split, Join, Fields
)

// getTerminalWidth retrieves the terminal width from the COLUMNS environment variable.
// If COLUMNS is not set or invalid, it defaults to 80 characters.
func getTerminalWidth() int {
	// Get the COLUMNS environment variable value
	if cols := os.Getenv("COLUMNS"); cols != "" {
		// Try to convert the string value to an integer
		if w, err := strconv.Atoi(cols); err == nil && w > 0 {
			// Return the width if conversion succeeded and value is positive
			return w
		}
	}
	// Default to 80 columns if COLUMNS is not set or invalid
	return 80
}

// applyAlignment applies horizontal alignment (left, center, right) to ASCII art lines.
// Justify alignment is handled separately by justifyInputForTerminal before rendering.
// Returns a new slice with aligned lines based on the specified width.
func applyAlignment(lines []string, align string, width int) []string {
	// If no alignment specified, left alignment, or invalid width, return lines unchanged
	if align == "" || align == "left" || width <= 0 {
		return lines
	}
	// Justify is handled before rendering, so skip it here
	if align == "justify" {
		return lines
	}

	// Create a new slice to store the aligned lines
	aligned := make([]string, len(lines))
	// Process each line individually
	for i, line := range lines {
		// Calculate the current line width in characters
		lineWidth := len(line)
		// If line is already wider than terminal or empty, keep it as-is
		if lineWidth >= width || line == "" {
			aligned[i] = line
			continue
		}

		// Calculate how many spaces to add on the left side
		var leftPad int
		switch align {
		case "right":
			// Right align: all padding goes on the left
			leftPad = width - lineWidth
		case "center":
			// Center align: split padding evenly (integer division)
			leftPad = (width - lineWidth) / 2
		default:
			// Unknown alignment: no padding
			leftPad = 0
		}

		// If padding calculation resulted in zero or negative, keep line unchanged
		if leftPad <= 0 {
			aligned[i] = line
			continue
		}
		// Add the calculated padding spaces before the line
		aligned[i] = strings.Repeat(" ", leftPad) + line
	}
	// Return the slice of aligned lines
	return aligned
}

// justifyInputForTerminal modifies the input text to add extra spaces between words
// so that when rendered as ASCII art, the output fills the terminal width.
// This is done before rendering because we need to adjust the input text itself.
func justifyInputForTerminal(input string, banner map[string][]string, termWidth int) string {
	// If terminal width is invalid, return input unchanged
	if termWidth <= 0 {
		return input
	}
	// Get the width of a space character in the ASCII art font
	spaceWidth := glyphWidth(" ", banner)
	// If space width is invalid, return input unchanged
	if spaceWidth <= 0 {
		return input
	}

	// Split input into separate lines
	inputLines := strings.Split(input, "\n")
	// Process each line individually
	for i, rawLine := range inputLines {
		// Split line into words (whitespace-separated tokens)
		words := strings.Fields(rawLine)
		// Skip lines with less than 2 words (can't justify single word or empty line)
		if len(words) < 2 {
			continue
		}

		// Calculate total width of all words when rendered as ASCII art
		contentWidth := 0
		for _, w := range words {
			// Add the rendered width of each word
			contentWidth += renderedTextWidth(w, banner)
		}

		// Calculate number of gaps between words (n words = n-1 gaps)
		gaps := len(words) - 1
		// Calculate total width with single spaces between words
		baseWidth := contentWidth + gaps*spaceWidth
		// If content already fills or exceeds terminal width, just join with single spaces
		if baseWidth >= termWidth {
			inputLines[i] = strings.Join(words, " ")
			continue
		}

		// Calculate how many extra columns we need to fill
		extraColumns := termWidth - baseWidth
		// Convert extra columns to number of extra spaces needed
		extraSpaces := extraColumns / spaceWidth
		// Distribute extra spaces evenly across all gaps
		baseAdd := extraSpaces / gaps
		// Calculate remainder to distribute to first few gaps
		remainder := extraSpaces % gaps

		// Build the justified line using a string builder for efficiency
		var b strings.Builder
		// Iterate through words and add appropriate spacing
		for idx, w := range words {
			// Write the current word
			b.WriteString(w)
			// Skip adding spaces after the last word
			if idx == gaps {
				continue
			}
			// Start with 1 space plus the base additional spaces
			spaces := 1 + baseAdd
			// Add one more space to the first 'remainder' gaps for even distribution
			if idx < remainder {
				spaces++
			}
			// Write the calculated number of spaces
			b.WriteString(strings.Repeat(" ", spaces))
		}
		// Replace the original line with the justified version
		inputLines[i] = b.String()
	}

	// Join all lines back together with newlines and return
	return strings.Join(inputLines, "\n")
}

// renderedTextWidth calculates the total width of text when rendered as ASCII art.
// It sums up the width of each character's glyph in the banner font.
func renderedTextWidth(text string, banner map[string][]string) int {
	// Initialize total width counter
	width := 0
	// Iterate through each rune (character) in the text
	for _, r := range text {
		// Add the width of this character's ASCII art glyph
		width += glyphWidth(string(r), banner)
	}
	// Return the total calculated width
	return width
}

// glyphWidth returns the width in characters of a single character's ASCII art glyph.
// The width is determined by the length of the first line of the glyph.
// Returns 4 as default width if the character is not found in the banner.
func glyphWidth(ch string, banner map[string][]string) int {
	// Look up the character in the banner map
	if glyph, ok := banner[ch]; ok && len(glyph) > 0 {
		// Return the length of the first line (all lines should have same width)
		return len(glyph[0])
	}
	// Default width for unknown characters (matches the 4-space padding in renderLines)
	return 4
}
