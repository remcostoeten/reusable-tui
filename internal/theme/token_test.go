package theme

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestEveryTokenRoundTrips(t *testing.T) {
	base := VioletDark()
	exported := Export(base)
	if len(exported.Tokens) != len(Tokens()) {
		t.Fatalf("exported %d tokens, registry lists %d", len(exported.Tokens), len(Tokens()))
	}
	for _, path := range Tokens() {
		value, ok := exported.Tokens[path]
		if !ok {
			t.Fatalf("token %q missing from export", path)
		}
		written, err := WriteToken(base, path, value)
		if err != nil {
			t.Fatalf("token %q is listed but not writable: %v", path, err)
		}
		readBack, err := ReadToken(written, path)
		if err != nil {
			t.Fatalf("token %q is listed but not readable: %v", path, err)
		}
		if readBack != value {
			t.Fatalf("token %q round-tripped %q as %q", path, value, readBack)
		}
	}
}

func TestApplyRejectsBadValues(t *testing.T) {
	cases := map[string]map[string]string{
		"unknown token": {"accent.glow": "#FFFFFF"},
		"bad hex":       {TokenAccentActive: "not-a-color"},
		"bad border":    {TokenBorderFocusedSet: "squiggly"},
		"empty label":   {TokenStatusDangerLabel: ""},
	}
	for name, tokens := range cases {
		t.Run(name, applyFailureCase(tokens))
	}
}

func applyFailureCase(tokens map[string]string) func(*testing.T) {
	return func(t *testing.T) {
		if _, err := Apply(VioletDark(), tokens); err == nil {
			t.Fatal("expected an error")
		}
	}
}

func TestApplyAcceptsAnsiIndex(t *testing.T) {
	patched, err := Apply(VioletDark(), map[string]string{TokenAccentActive: "213"})
	if err != nil {
		t.Fatalf("ansi index rejected: %v", err)
	}
	if got := FormatColor(patched.Accent.Active); got != "213" {
		t.Fatalf("accent is %q", got)
	}
}

func TestFileExtendsBuiltin(t *testing.T) {
	built, err := File{
		Name:    "cyan",
		Extends: NameVioletDark,
		Tokens:  map[string]string{TokenAccentActive: "#22d3ee"},
	}.Build(Builtin())
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if built.Name != "cyan" {
		t.Fatalf("name is %q", built.Name)
	}
	if got := FormatColor(built.Accent.Active); got != "#22D3EE" {
		t.Fatalf("accent is %q", got)
	}
	if got := FormatColor(built.Base.Background); got != FormatColor(VioletDark().Base.Background) {
		t.Fatalf("inherited background is %q", got)
	}
}

func TestFileRejectsUnknownParent(t *testing.T) {
	_, err := File{Name: "orphan", Extends: "nope"}.Build(Builtin())
	if !errors.Is(err, ErrUnknownTheme) {
		t.Fatalf("expected ErrUnknownTheme, got %v", err)
	}
}

func TestComposeLoadsUserThemes(t *testing.T) {
	dir := t.TempDir()
	writeThemeFile(t, filepath.Join(dir, "cyan.json"), `{"name":"cyan","extends":"violet-dark","tokens":{"accent.active":"#22d3ee"}}`)
	writeThemeFile(t, filepath.Join(dir, "broken.json"), `{"name":"broken","tokens":{"accent.active":"nope"}}`)
	writeThemeFile(t, filepath.Join(dir, "notes.txt"), `ignored`)

	registry, failures := Compose(dir, map[string]map[string]string{
		NameMonochrome: {TokenAccentActive: "#FF00FF"},
	})
	if len(failures) != 1 {
		t.Fatalf("expected one failure, got %v", failures)
	}
	cyan, ok := registry.Get("cyan")
	if !ok {
		t.Fatal("user theme was not registered")
	}
	if got := FormatColor(cyan.Accent.Active); got != "#22D3EE" {
		t.Fatalf("user theme accent is %q", got)
	}
	mono, _ := registry.Get(NameMonochrome)
	if got := FormatColor(mono.Accent.Active); got != "#FF00FF" {
		t.Fatalf("override did not apply, accent is %q", got)
	}
}

func TestComposeMissingDirIsNotAnError(t *testing.T) {
	registry, failures := Compose(filepath.Join(t.TempDir(), "absent"), nil)
	if len(failures) != 0 {
		t.Fatalf("expected no failures, got %v", failures)
	}
	if len(registry.Names()) != len(Builtin().Names()) {
		t.Fatal("builtin themes are missing")
	}
}

func TestWriteFileRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "mono.json")
	if err := WriteFile(path, Export(Monochrome())); err != nil {
		t.Fatalf("write: %v", err)
	}
	back, err := ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	built, err := back.Build(Builtin())
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	for _, path := range Tokens() {
		want, _ := ReadToken(Monochrome(), path)
		got, _ := ReadToken(built, path)
		if want != got {
			t.Fatalf("token %q survived export as %q, expected %q", path, got, want)
		}
	}
}

func writeThemeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
