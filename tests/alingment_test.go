// Package tests contains unit tests for the ASCII art application
package tests

import (
	"bytes"   // Provides buffer for capturing output
	"strings" // String manipulation functions for parsing test output
	"testing" // Go's testing framework

	"ascii-art/pipeline" // The pipeline package being tested
)

// TestAlingment_RightAddsMoreLeftPaddingThanLeft verifies that right alignment
// adds more left padding than left alignment (right-aligned text should be pushed to the right).
func TestAlingment_RightAddsMoreLeftPaddingThanLeft(t *testing.T) {
	// Set terminal width to 120 columns for this test
	t.Setenv("COLUMNS", "120")

	// Create a buffer to capture left-aligned output
	var left bytes.Buffer
	// Run pipeline with left alignment and capture output
	if code := pipeline.Run([]string{"--align=left", "hello", "standard"}, &left); code != 0 {
		// If pipeline returns non-zero exit code, test fails
		t.Fatalf("expected left alignment success, got %d", code)
	}

	// Create a buffer to capture right-aligned output
	var right bytes.Buffer
	// Run pipeline with right alignment and capture output
	if code := pipeline.Run([]string{"--align=right", "hello", "standard"}, &right); code != 0 {
		// If pipeline returns non-zero exit code, test fails
		t.Fatalf("expected right alignment success, got %d", code)
	}

	// Count leading spaces in the first non-empty line of left-aligned output
	leftPad := leadingSpaces(firstNonEmptyLine(left.String()))
	// Count leading spaces in the first non-empty line of right-aligned output
	rightPad := leadingSpaces(firstNonEmptyLine(right.String()))

	// Verify that right alignment has more left padding than left alignment
	if rightPad <= leftPad {
		// If right padding is not greater, test fails
		t.Fatalf("expected right alignment to have more left padding (left=%d, right=%d)", leftPad, rightPad)
	}
}

// TestAlingment_CenterPaddingBetweenLeftAndRight verifies that center alignment
// produces padding that is between left alignment (minimal) and right alignment (maximal).
func TestAlingment_CenterPaddingBetweenLeftAndRight(t *testing.T) {
	// Set terminal width to 120 columns for this test
	t.Setenv("COLUMNS", "120")

	// Create a buffer to capture left-aligned output
	var left bytes.Buffer
	// Run pipeline with left alignment and capture output
	if code := pipeline.Run([]string{"--align=left", "hello", "standard"}, &left); code != 0 {
		// If pipeline returns non-zero exit code, test fails
		t.Fatalf("expected left alignment success, got %d", code)
	}

	// Create a buffer to capture center-aligned output
	var center bytes.Buffer
	// Run pipeline with center alignment and capture output
	if code := pipeline.Run([]string{"--align=center", "hello", "standard"}, &center); code != 0 {
		// If pipeline returns non-zero exit code, test fails
		t.Fatalf("expected center alignment success, got %d", code)
	}

	// Create a buffer to capture right-aligned output
	var right bytes.Buffer
	// Run pipeline with right alignment and capture output
	if code := pipeline.Run([]string{"--align=right", "hello", "standard"}, &right); code != 0 {
		// If pipeline returns non-zero exit code, test fails
		t.Fatalf("expected right alignment success, got %d", code)
	}

	// Count leading spaces in the first non-empty line of left-aligned output
	leftPad := leadingSpaces(firstNonEmptyLine(left.String()))
	// Count leading spaces in the first non-empty line of center-aligned output
	centerPad := leadingSpaces(firstNonEmptyLine(center.String()))
	// Count leading spaces in the first non-empty line of right-aligned output
	rightPad := leadingSpaces(firstNonEmptyLine(right.String()))

	// Verify that center padding is strictly between left and right padding
	if !(centerPad > leftPad && centerPad < rightPad) {
		// If center padding is not in the middle, test fails
		t.Fatalf("expected center padding to be between left and right (left=%d, center=%d, right=%d)", leftPad, centerPad, rightPad)
	}
}

// TestAlingment_JustifyExpandsLineWidth verifies that justify alignment
// expands the line width by adding extra spaces between words to fill the terminal width.
func TestAlingment_JustifyExpandsLineWidth(t *testing.T) {
	// Set terminal width to 120 columns for this test
	t.Setenv("COLUMNS", "120")

	// Create a buffer to capture left-aligned output
	var left bytes.Buffer
	// Run pipeline with left alignment (normal spacing) and capture output
	if code := pipeline.Run([]string{"--align=left", "how are you", "shadow"}, &left); code != 0 {
		// If pipeline returns non-zero exit code, test fails
		t.Fatalf("expected left alignment success, got %d", code)
	}

	// Create a buffer to capture justified output
	var justify bytes.Buffer
	// Run pipeline with justify alignment (expanded spacing) and capture output
	if code := pipeline.Run([]string{"--align=justify", "how are you", "shadow"}, &justify); code != 0 {
		// If pipeline returns non-zero exit code, test fails
		t.Fatalf("expected justify alignment success, got %d", code)
	}

	// Get the first non-empty line from left-aligned output
	leftLine := firstNonEmptyLine(left.String())
	// Get the first non-empty line from justified output
	justifyLine := firstNonEmptyLine(justify.String())

	// Verify that justified line is wider than left-aligned line (due to extra spaces)
	if len(justifyLine) <= len(leftLine) {
		// If justified line is not wider, test fails
		t.Fatalf("expected justify line to be wider than left line (left=%d, justify=%d)", len(leftLine), len(justifyLine))
	}
}

// firstNonEmptyLine extracts the first line from output that contains non-whitespace characters.
// This is useful for testing alignment since ASCII art may have empty lines at the start.
func firstNonEmptyLine(out string) string {
	// Split output into lines (remove trailing newlines first, then split by newline)
	for _, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
		// Check if line has any non-whitespace content
		if strings.TrimSpace(line) != "" {
			// Return the first line with actual content
			return line
		}
	}
	// If no non-empty line found, return empty string
	return ""
}

// leadingSpaces counts the number of space characters at the beginning of a string.
// This is used to measure the left padding added by alignment operations.
func leadingSpaces(s string) int {
	// Initialize counter for leading spaces
	count := 0
	// Iterate through each byte in the string
	for i := 0; i < len(s); i++ {
		// Stop counting when we encounter a non-space character
		if s[i] != ' ' {
			break
		}
		// Increment counter for each space found
		count++
	}
	// Return the total count of leading spaces
	return count
}
