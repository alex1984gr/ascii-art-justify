// Package tests contains all unit tests for argument parsing
package tests

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"ascii-art/pipeline"
)

func TestArgs_SingleString_DefaultBanner(t *testing.T) {
	var out bytes.Buffer
	if code := pipeline.Run([]string{"hello"}, &out); code != 0 {
		t.Fatalf("expected success, got exit code %d", code)
	}
	if out.Len() == 0 {
		t.Fatal("expected non-empty output")
	}
}

func TestArgs_StringAndBanner(t *testing.T) {
	var out bytes.Buffer
	if code := pipeline.Run([]string{"hello", "shadow"}, &out); code != 0 {
		t.Fatalf("expected success, got exit code %d", code)
	}
}

func TestArgs_InvalidBanner(t *testing.T) {
	var out bytes.Buffer
	if code := pipeline.Run([]string{"hello", "badbanner"}, &out); code == 0 {
		t.Fatal("expected failure for invalid banner")
	}
}

func TestArgs_TooManyPositionals(t *testing.T) {
	var out bytes.Buffer
	if code := pipeline.Run([]string{"a", "shadow", "extra"}, &out); code == 0 {
		t.Fatal("expected failure for too many positional args")
	}
}

func TestArgs_InvalidFlagFormat(t *testing.T) {
	var out bytes.Buffer
	if code := pipeline.Run([]string{"--output", "hello"}, &out); code == 0 {
		t.Fatal("expected failure for invalid --output format")
	}
}

func TestArgs_UnknownFlag(t *testing.T) {
	var out bytes.Buffer
	if code := pipeline.Run([]string{"--unknown=1", "hello"}, &out); code == 0 {
		t.Fatal("expected failure for unknown flag")
	}
}

func TestArgs_ColorInputBanner(t *testing.T) {
	var out bytes.Buffer
	if code := pipeline.Run([]string{"--color=red", "hello", "thinkertoy"}, &out); code != 0 {
		t.Fatalf("expected success, got exit code %d", code)
	}
}

func TestArgs_ColorSubstringInputBanner(t *testing.T) {
	var out bytes.Buffer
	if code := pipeline.Run([]string{"--color=red", "ll", "hello", "standard"}, &out); code != 0 {
		t.Fatalf("expected success, got exit code %d", code)
	}
}

func TestArgs_FontFlagAndPositionalBannerConflict(t *testing.T) {
	var out bytes.Buffer
	if code := pipeline.Run([]string{"--font=standard", "hello", "shadow"}, &out); code == 0 {
		t.Fatal("expected failure when --font and positional banner are both provided")
	}
}

func TestArgs_OutputFlags(t *testing.T) {
	t.Run("--out", func(t *testing.T) {
		var out bytes.Buffer
		path := filepath.Join(t.TempDir(), "out-short.txt")
		if code := pipeline.Run([]string{"--out=" + path, "hello"}, &out); code != 0 {
			t.Fatalf("expected success, got exit code %d", code)
		}
	})

	t.Run("--output", func(t *testing.T) {
		var out bytes.Buffer
		path := filepath.Join(t.TempDir(), "out-long.txt")
		if code := pipeline.Run([]string{"--output=" + path, "hello"}, &out); code != 0 {
			t.Fatalf("expected success, got exit code %d", code)
		}
	})
}

func TestArgs_AlignValidTypes(t *testing.T) {
	t.Setenv("COLUMNS", "120")
	for _, align := range []string{"left", "center", "right", "justify"} {
		var out bytes.Buffer
		if code := pipeline.Run([]string{"--align=" + align, "hello", "standard"}, &out); code != 0 {
			t.Fatalf("expected success for align=%s, got %d", align, code)
		}
	}
}

func TestArgs_AlignInvalidType(t *testing.T) {
	var out bytes.Buffer
	if code := pipeline.Run([]string{"--align=middle", "hello", "standard"}, &out); code == 0 {
		t.Fatal("expected failure for invalid align type")
	}
}

func TestArgs_AlignInvalidFormat(t *testing.T) {
	var out bytes.Buffer
	if code := pipeline.Run([]string{"--align", "hello", "standard"}, &out); code == 0 {
		t.Fatal("expected failure for invalid --align format")
	}
}

func TestArgs_AlignChangesOutputPosition(t *testing.T) {
	t.Setenv("COLUMNS", "120")

	var leftOut bytes.Buffer
	if code := pipeline.Run([]string{"--align=left", "hello", "standard"}, &leftOut); code != 0 {
		t.Fatalf("expected left align success, got %d", code)
	}

	var rightOut bytes.Buffer
	if code := pipeline.Run([]string{"--align=right", "hello", "standard"}, &rightOut); code != 0 {
		t.Fatalf("expected right align success, got %d", code)
	}

	if leftOut.String() == rightOut.String() {
		t.Fatal("expected right aligned output to differ from left aligned output")
	}

	lines := strings.Split(strings.TrimRight(rightOut.String(), "\n"), "\n")
	foundPrefixed := false
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.HasPrefix(line, " ") {
			foundPrefixed = true
			break
		}
	}
	if !foundPrefixed {
		t.Fatal("expected right alignment to add left padding")
	}
}
