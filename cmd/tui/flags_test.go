package main

import (
	"errors"
	"flag"
	"io"
	"testing"
)

func TestParseFlags(t *testing.T) {
	f, err := parseFlags([]string{"--config", "c.json", "--db", "x.db", "--theme", "monochrome", "--no-mouse"}, io.Discard)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if f.configPath != "c.json" || f.dbPath != "x.db" || f.themeName != "monochrome" || !f.noMouse {
		t.Fatalf("unexpected flags %+v", f)
	}
}

func TestParseFlagsVersion(t *testing.T) {
	_, err := parseFlags([]string{"--version"}, io.Discard)
	if !errors.Is(err, errShowVersion) {
		t.Fatalf("expected version sentinel, got %v", err)
	}
}

func TestParseFlagsUnknown(t *testing.T) {
	_, err := parseFlags([]string{"--nope"}, io.Discard)
	if err == nil || errors.Is(err, flag.ErrHelp) {
		t.Fatalf("expected a parse error, got %v", err)
	}
}
