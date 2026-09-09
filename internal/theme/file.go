package theme

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/adrg/xdg"
)

const (
	fileExtension = ".json"
	themesDirName = "themes"
)

var ErrNamelessTheme = errors.New("theme: file has no name")

type File struct {
	Name    string            `json:"name"`
	Extends string            `json:"extends,omitempty"`
	Tokens  map[string]string `json:"tokens"`
}

func UserDir(appName string) string {
	return filepath.Join(xdg.ConfigHome, appName, themesDirName)
}

func Export(t Theme) File {
	tokens := make(map[string]string, len(tokenOrder))
	for _, path := range tokenOrder {
		value, err := ReadToken(t, path)
		if err != nil {
			continue
		}
		tokens[path] = value
	}
	return File{Name: t.Name, Tokens: tokens}
}

func Apply(t Theme, tokens map[string]string) (Theme, error) {
	paths := make([]string, 0, len(tokens))
	for path := range tokens {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	failures := make([]error, 0, len(paths))
	for _, path := range paths {
		next, err := WriteToken(t, path, tokens[path])
		if err != nil {
			failures = append(failures, err)
			continue
		}
		t = next
	}
	return t, errors.Join(failures...)
}

func (f File) Build(base *Registry) (Theme, error) {
	if strings.TrimSpace(f.Name) == "" {
		return Theme{}, ErrNamelessTheme
	}
	start := base.Default()
	if f.Extends != "" {
		parent, ok := base.Get(f.Extends)
		if !ok {
			return Theme{}, fmt.Errorf("%w: %s extends %s", ErrUnknownTheme, f.Name, f.Extends)
		}
		start = parent
	}
	start.Name = f.Name
	built, err := Apply(start, f.Tokens)
	if err != nil {
		return Theme{}, fmt.Errorf("%s: %w", f.Name, err)
	}
	return built, nil
}

func WriteFile(path string, f File) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(raw, '\n'), 0o644)
}

func ReadFile(path string) (File, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return File{}, err
	}
	var f File
	if err := json.Unmarshal(raw, &f); err != nil {
		return File{}, fmt.Errorf("%s: %w", filepath.Base(path), err)
	}
	return f, nil
}

func LoadDir(dir string, base *Registry) ([]Theme, []error) {
	entries, err := os.ReadDir(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, []error{err}
	}
	themes := make([]Theme, 0, len(entries))
	failures := make([]error, 0)
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != fileExtension {
			continue
		}
		built, err := loadOne(filepath.Join(dir, entry.Name()), base)
		if err != nil {
			failures = append(failures, err)
			continue
		}
		themes = append(themes, built)
	}
	return themes, failures
}

func loadOne(path string, base *Registry) (Theme, error) {
	f, err := ReadFile(path)
	if err != nil {
		return Theme{}, err
	}
	return f.Build(base)
}

func ApplyOverrides(r *Registry, overrides map[string]map[string]string) []error {
	names := make([]string, 0, len(overrides))
	for name := range overrides {
		names = append(names, name)
	}
	sort.Strings(names)
	failures := make([]error, 0)
	for _, name := range names {
		current, ok := r.Get(name)
		if !ok {
			failures = append(failures, fmt.Errorf("%w: override for %s", ErrUnknownTheme, name))
			continue
		}
		patched, err := Apply(current, overrides[name])
		if err != nil {
			failures = append(failures, fmt.Errorf("override %s: %w", name, err))
			continue
		}
		r.Add(patched)
	}
	return failures
}

func Compose(dir string, overrides map[string]map[string]string) (*Registry, []error) {
	registry := Builtin()
	userThemes, failures := LoadDir(dir, registry)
	for _, t := range userThemes {
		registry.Add(t)
	}
	return registry, append(failures, ApplyOverrides(registry, overrides)...)
}
