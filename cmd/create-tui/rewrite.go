package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// rewritable are the file extensions that can name the module path or the
// binary. Anything else — golden files, binaries — is copied untouched.
var rewritable = map[string]bool{
	".go":   true,
	".mod":  true,
	".md":   true,
	".yml":  true,
	".yaml": true,
	".sh":   true,
	".toml": true,
	"":      true, // Makefile, Dockerfile
}

// rewrite replaces every mention of the source module and binary with the new
// project's, returning how many files changed.
//
// The module path is a plain string substitution rather than an AST rewrite
// because it appears in go.mod, the Makefile and the CI workflow as well as in
// imports, and it is specific enough that a false positive is not possible.
func rewrite(o options) (int, error) {
	replacer := strings.NewReplacer(
		sourceModule, o.Module,
		"./cmd/"+sourceBinary, "./cmd/"+o.Name,
		"cmd/"+sourceBinary+"/", "cmd/"+o.Name+"/",
		"bin/"+sourceBinary, "bin/"+o.Name,
	)

	changed := 0
	err := filepath.WalkDir(o.Dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == ".claude" {
			// Skills point back at the template repository on purpose.
			return filepath.SkipDir
		}
		if d.IsDir() || !rewritable[filepath.Ext(path)] {
			return nil
		}

		before, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		after := replacer.Replace(string(before))
		if after == string(before) {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(after), info.Mode().Perm()); err != nil {
			return err
		}
		changed++
		return nil
	})
	return changed, err
}

// renameBinary moves cmd/tui to cmd/<name> and renames the const the
// entrypoint reports itself under. The import paths inside it are fixed by
// rewrite, which runs after.
func renameBinary(o options) error {
	if o.Name == sourceBinary {
		return nil
	}

	from := filepath.Join(o.Dir, "cmd", sourceBinary)
	to := filepath.Join(o.Dir, "cmd", o.Name)
	if _, err := os.Stat(from); os.IsNotExist(err) {
		return nil
	}
	if err := os.Rename(from, to); err != nil {
		return err
	}

	return replaceInFile(filepath.Join(to, "main.go"),
		`name    = "demo"`, `name    = "`+o.Name+`"`,
		"// Command tui runs the demo application.",
		"// Command "+o.Name+" runs the application.",
	)
}

func replaceInFile(path string, pairs ...string) error {
	before, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	after := strings.NewReplacer(pairs...).Replace(string(before))
	if after == string(before) {
		return nil
	}
	return os.WriteFile(path, []byte(after), 0o644)
}
