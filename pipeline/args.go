// Package pipeline contains argument parsing logic for the ASCII art application
package pipeline

import (
	"fmt"    // Formatted I/O functions for creating error messages
	"strings" // String manipulation functions like HasPrefix, TrimPrefix, ToLower
)

// usageMessage is the help text displayed when arguments are invalid
const usageMessage = "Usage: go run . [OPTION] [STRING] [BANNER]\n\nExample: go run . --align=right something standard"

// runConfig holds all configuration options parsed from command-line arguments
type runConfig struct {
	font      string // Banner font name (standard, shadow, thinkertoy)
	outFile   string // Output file path (empty means stdout)
	input     string // Text to convert to ASCII art
	colorName string // ANSI color name for coloring output
	substring string // Substring to color (empty means color all)
	align     string // Alignment type (left, center, right, justify)
	fontSet   bool   // Whether --font flag was explicitly set
}

// parseArgs parses command-line arguments and returns a runConfig or an error.
// It handles both optional flags (--font, --color, --align, --out) and positional arguments.
func parseArgs(args []string) (runConfig, error) {
	// Initialize config with default values: standard font and left alignment
	cfg := runConfig{font: "standard", align: "left"}
	// Slice to collect non-flag positional arguments
	positionals := make([]string, 0, len(args))

	// First pass: process all arguments to separate flags from positionals
	for _, arg := range args {
		switch {
		// Handle --font=<value> flag
		case strings.HasPrefix(arg, "--font="):
			// Extract font name after the equals sign
			cfg.font = strings.TrimPrefix(arg, "--font=")
			// Mark that font was explicitly set (affects positional arg parsing)
			cfg.fontSet = true
			// Reject empty font name
			if cfg.font == "" {
				return runConfig{}, fmt.Errorf("empty font")
			}
		// Handle invalid --font without value
		case arg == "--font":
			return runConfig{}, fmt.Errorf("invalid font format")
		// Handle --out=<value> flag
		case strings.HasPrefix(arg, "--out="):
			// Extract output file path after the equals sign
			cfg.outFile = strings.TrimPrefix(arg, "--out=")
			// Reject empty output file path
			if cfg.outFile == "" {
				return runConfig{}, fmt.Errorf("empty out file")
			}
		// Handle --output=<value> flag (alias for --out)
		case strings.HasPrefix(arg, "--output="):
			// Extract output file path after the equals sign
			cfg.outFile = strings.TrimPrefix(arg, "--output=")
			// Reject empty output file path
			if cfg.outFile == "" {
				return runConfig{}, fmt.Errorf("empty output file")
			}
		// Handle invalid --out or --output without value
		case arg == "--out" || arg == "--output":
			return runConfig{}, fmt.Errorf("invalid output format")
		// Handle --color=<value> flag
		case strings.HasPrefix(arg, "--color="):
			// Extract color name after the equals sign
			cfg.colorName = strings.TrimPrefix(arg, "--color=")
			// Reject empty color name
			if cfg.colorName == "" {
				return runConfig{}, fmt.Errorf("empty color")
			}
		// Handle invalid --color without value
		case arg == "--color":
			return runConfig{}, fmt.Errorf("invalid color format")
		// Handle --align=<value> flag
		case strings.HasPrefix(arg, "--align="):
			// Extract alignment type and convert to lowercase
			cfg.align = strings.ToLower(strings.TrimPrefix(arg, "--align="))
			// Validate alignment type is one of: left, center, right, justify
			if !isAlignType(cfg.align) {
				return runConfig{}, fmt.Errorf("invalid align type")
			}
		// Handle invalid --align without value
		case arg == "--align":
			return runConfig{}, fmt.Errorf("invalid align format")
		// Reject any unknown flags starting with --
		case strings.HasPrefix(arg, "--"):
			return runConfig{}, fmt.Errorf("unknown option")
		// Collect non-flag arguments as positionals
		default:
			positionals = append(positionals, arg)
		}
	}

	// If no color flag was specified, parse positionals as: [input] or [input banner]
	if cfg.colorName == "" {
		switch len(positionals) {
		// Single positional: it's the input text
		case 1:
			cfg.input = positionals[0]
		// Two positionals: input and banner name
		case 2:
			// If font was already set via --font flag, reject extra argument
			if cfg.fontSet {
				return runConfig{}, fmt.Errorf("too many args")
			}
			// Validate second argument is a valid banner name
			if !isBannerName(positionals[1]) {
				return runConfig{}, fmt.Errorf("invalid banner")
			}
			// First positional is input, second is banner
			cfg.input = positionals[0]
			cfg.font = positionals[1]
		// Any other number of positionals is invalid
		default:
			return runConfig{}, fmt.Errorf("invalid args")
		}
		// Return successfully parsed config
		return cfg, nil
	}

	// If color flag was specified, parse positionals as: [input], [substring input], or [substring input banner]
	switch len(positionals) {
	// Single positional: it's the input text (color entire output)
	case 1:
		cfg.input = positionals[0]
	// Two positionals: could be [input banner] or [substring input]
	case 2:
		// If font not set and second arg is a banner name, treat as [input banner]
		if !cfg.fontSet && isBannerName(positionals[1]) {
			cfg.input = positionals[0]
			cfg.font = positionals[1]
		} else {
			// Otherwise treat as [substring input] for coloring specific substring
			cfg.substring = positionals[0]
			cfg.input = positionals[1]
		}
	// Three positionals: [substring input banner]
	case 3:
		// If font was already set via --font flag, reject extra argument
		if cfg.fontSet {
			return runConfig{}, fmt.Errorf("too many args")
		}
		// Validate third argument is a valid banner name
		if !isBannerName(positionals[2]) {
			return runConfig{}, fmt.Errorf("invalid banner")
		}
		// First is substring to color, second is input, third is banner
		cfg.substring = positionals[0]
		cfg.input = positionals[1]
		cfg.font = positionals[2]
	// Any other number of positionals is invalid
	default:
		return runConfig{}, fmt.Errorf("invalid args")
	}

	// Return successfully parsed config
	return cfg, nil
}

// isBannerName checks if the given name is a valid banner font name.
// Valid names are: standard, shadow, thinkertoy
func isBannerName(name string) bool {
	return name == "standard" || name == "shadow" || name == "thinkertoy"
}

// isAlignType checks if the given alignment type is valid.
// Valid types are: left, center, right, justify
func isAlignType(align string) bool {
	return align == "left" || align == "center" || align == "right" || align == "justify"
}
