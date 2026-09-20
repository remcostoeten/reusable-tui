package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// maxFileSize caps a single extracted file. Nothing in the template is close
// to it; the limit exists so a malformed archive cannot fill the disk.
const maxFileSize = 8 << 20

// skipped are paths, relative to the repository root, that never belong in a
// scaffolded project: the scaffolder itself, this repository's README — which
// documents the scaffolder, and would be rewritten into nonsense — and a stale
// nested checkout that predates the current architecture.
var skipped = map[string]bool{
	".git":           true,
	"cmd/create-tui": true,
	"reusable-tui":   true,
	"README.md":      true,
}

// copyLocal populates dir from a checkout on disk instead of a download,
// applying the same exclusions. It is how the scaffolder is exercised against
// unreleased changes, before a ref exists to fetch them from.
func copyLocal(src, dir string) error {
	root := filepath.Clean(src)

	return filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		if rel == "." {
			return os.MkdirAll(dir, 0o755)
		}
		if skip(filepath.ToSlash(rel)) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		target := filepath.Join(dir, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if !d.Type().IsRegular() {
			return nil
		}

		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()

		info, err := d.Info()
		if err != nil {
			return err
		}
		return writeFile(f, target, info.Mode())
	})
}

// fetch downloads a repository ref as a tarball and extracts it into dir,
// dropping the archive's top-level directory and everything in skipped.
//
// Codeload serves the same archive GitHub's "Download ZIP" button does, with
// no authentication and no git binary required.
func fetch(repo, ref, dir string) error {
	url := fmt.Sprintf("https://codeload.github.com/%s/tar.gz/%s", repo, ref)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %s (is %q a real branch or tag?)", url, resp.Status, ref)
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer func() { _ = gz.Close() }()

	return extract(tar.NewReader(gz), dir)
}

func extract(tr *tar.Reader, dir string) error {
	for {
		header, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		rel, ok := strip(header.Name)
		if !ok || skip(rel) {
			continue
		}

		target, err := safeJoin(dir, rel)
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := writeFile(tr, target, header.FileInfo().Mode()); err != nil {
				return err
			}
		}
	}
}

func writeFile(r io.Reader, target string, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}

	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode.Perm())
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	if _, err := io.Copy(f, io.LimitReader(r, maxFileSize)); err != nil {
		return err
	}
	return f.Close()
}

// strip removes the archive's top-level directory, which codeload names after
// the repository and the ref. The root entry itself yields no path.
func strip(name string) (string, bool) {
	name = path.Clean(strings.TrimPrefix(name, "./"))
	_, rest, found := strings.Cut(name, "/")
	if !found || rest == "" || rest == "." {
		return "", false
	}
	return rest, true
}

// skip reports whether a path is excluded, matching directories by prefix so
// that one entry covers a whole subtree.
func skip(rel string) bool {
	for p := rel; p != "." && p != "/"; p = path.Dir(p) {
		if skipped[p] {
			return true
		}
	}
	return false
}

// safeJoin refuses a path that would escape the target directory, which is
// what stops a hostile archive from writing outside it.
func safeJoin(dir, rel string) (string, error) {
	target := filepath.Join(dir, filepath.FromSlash(rel))
	if !strings.HasPrefix(target, filepath.Clean(dir)+string(os.PathSeparator)) {
		return "", fmt.Errorf("archive entry escapes target directory: %q", rel)
	}
	return target, nil
}
