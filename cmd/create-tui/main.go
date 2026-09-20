// Command create-tui scaffolds a new terminal application from this
// repository.
//
//	go run github.com/remcostoeten/reusable-tui/cmd/create-tui@latest myapp
//
// It downloads the repository at a ref, rewrites every import to the new
// module path, optionally strips the demo modules, and leaves behind a project
// that builds and runs. It depends on nothing outside the standard library, so
// the `go run` above resolves in a second rather than pulling the whole
// dependency graph of the shell it installs.
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// source is the repository this scaffolder copies from. It is also the import
// prefix every generated file has rewritten away.
const (
	sourceRepo   = "remcostoeten/reusable-tui"
	sourceModule = "github.com/" + sourceRepo
	sourceBinary = "tui"
)

// options is the resolved answer to every question the scaffolder asks.
type options struct {
	Dir    string
	Module string
	Name   string
	Ref    string
	Local  string
	Demo   bool
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "\n  error:", err)
		os.Exit(1)
	}
}

func run() error {
	module := flag.String("module", "", "go module path for the new project")
	name := flag.String("name", "", "application and binary name")
	ref := flag.String("ref", "master", "branch or tag to scaffold from")
	local := flag.String("local", "", "scaffold from a checkout on disk instead of downloading")
	demo := flag.Bool("demo", false, "keep the demo modules (home, manager, settings)")
	yes := flag.Bool("y", false, "accept every default without prompting")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		return errors.New("expected exactly one target directory")
	}

	opts, err := resolve(flag.Arg(0), *module, *name, *ref, *demo, *yes)
	if err != nil {
		return err
	}
	opts.Local = *local

	fmt.Printf("\n  module path   %s\n", opts.Module)
	fmt.Printf("  binary        %s\n", opts.Name)
	if opts.Demo {
		fmt.Printf("  demo modules  kept\n\n")
	} else {
		fmt.Printf("  demo modules  pruned\n\n")
	}

	return scaffold(opts)
}

// scaffold performs the whole build, reporting each step as it completes. A
// failure past the fetch leaves the target directory in place for inspection
// rather than deleting work the caller might want to look at.
func scaffold(o options) error {
	if o.Local != "" {
		if err := copyLocal(o.Local, o.Dir); err != nil {
			return fmt.Errorf("copy %s: %w", o.Local, err)
		}
		step("copied %s", o.Local)
	} else {
		if err := fetch(sourceRepo, o.Ref, o.Dir); err != nil {
			return fmt.Errorf("fetch %s@%s: %w", sourceRepo, o.Ref, err)
		}
		step("fetched %s@%s", sourceRepo, o.Ref)
	}

	if !o.Demo {
		if err := prune(o); err != nil {
			return fmt.Errorf("prune demo: %w", err)
		}
		step("pruned demo modules, wrote internal/modules/hello")
	}

	if err := renameBinary(o); err != nil {
		return fmt.Errorf("rename cmd/%s: %w", sourceBinary, err)
	}

	n, err := rewrite(o)
	if err != nil {
		return fmt.Errorf("rewrite imports: %w", err)
	}
	step("rewrote module path (%d files)", n)

	// After the rewrite, so that the link back to this repository is not
	// itself rewritten into a link to the new project.
	if err := write(o.Dir, "README.md", readme(o)); err != nil {
		return fmt.Errorf("write README: %w", err)
	}

	if err := tidy(o.Dir); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}
	step("go mod tidy")

	if err := initRepo(o.Dir); err != nil {
		fmt.Printf("  · skipped git init: %v\n", err)
	} else {
		step("git init")
	}

	fmt.Printf("\n  cd %s && make run\n\n", o.Dir)
	return nil
}

// resolve fills in whatever the flags left blank, prompting when there is a
// terminal to prompt on and falling back to the derived default when there is
// not.
func resolve(dir, module, name, ref string, demo, yes bool) (options, error) {
	dir = filepath.Clean(dir)
	if err := checkTarget(dir); err != nil {
		return options{}, err
	}

	base := sanitize(filepath.Base(dir))
	if base == "" {
		return options{}, fmt.Errorf("cannot derive an application name from %q", dir)
	}

	o := options{Dir: dir, Module: module, Name: name, Ref: ref, Demo: demo}
	interactive := !yes && isTerminal(os.Stdin)

	if o.Module == "" {
		suggested := guessModule(base)
		o.Module = suggested
		if interactive {
			o.Module = ask("module path", suggested)
		}
	}
	if o.Name == "" {
		o.Name = base
		if interactive {
			o.Name = sanitize(ask("application name", base))
		}
	}
	if !o.Demo && interactive && !demoFlagSet() {
		o.Demo = confirm("keep the demo modules", false)
	}

	if o.Name == "" {
		return options{}, errors.New("application name cannot be empty")
	}
	if strings.TrimSpace(o.Module) == "" {
		return options{}, errors.New("module path cannot be empty")
	}
	return o, nil
}

// checkTarget refuses to scaffold over anything that already has content, so a
// mistyped directory never clobbers a project.
func checkTarget(dir string) error {
	entries, err := os.ReadDir(dir)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return nil
	case err != nil:
		return err
	case len(entries) > 0:
		return fmt.Errorf("%s already exists and is not empty", dir)
	}
	return nil
}

// guessModule builds a plausible module path from the git host username, which
// is right often enough to be worth offering as the default.
func guessModule(base string) string {
	out, err := exec.Command("git", "config", "--get", "user.name").Output()
	user := strings.TrimSpace(string(out))
	if err != nil || user == "" || strings.ContainsAny(user, " \t") {
		return base
	}
	return "github.com/" + strings.ToLower(user) + "/" + base
}

// sanitize reduces a string to what is legal in a Go package name and a binary
// name at once: lowercase letters, digits and dashes collapsed from the rest.
func sanitize(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(s)) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_':
			b.WriteRune('-')
		}
	}
	return strings.Trim(b.String(), "-")
}

func step(format string, args ...any) {
	fmt.Printf("  ✓ "+format+"\n", args...)
}

var stdin = bufio.NewReader(os.Stdin)

func ask(label, fallback string) string {
	fmt.Printf("  %s [%s]: ", label, fallback)
	line, err := stdin.ReadString('\n')
	if err != nil {
		fmt.Println()
		return fallback
	}
	if answer := strings.TrimSpace(line); answer != "" {
		return answer
	}
	return fallback
}

func confirm(label string, fallback bool) bool {
	hint := "y/N"
	if fallback {
		hint = "Y/n"
	}
	switch strings.ToLower(ask(label, hint)) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		return fallback
	}
}

// demoFlagSet reports whether -demo was passed explicitly, so that an explicit
// `-demo=false` is not undone by the interactive prompt.
func demoFlagSet() bool {
	set := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "demo" {
			set = true
		}
	})
	return set
}

// isTerminal reports whether a file is attached to a terminal. Checking the
// mode bits avoids a dependency for a single syscall.
func isTerminal(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func tidy(dir string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// initRepo starts a fresh history. A failure here is reported but not fatal:
// the project is already usable without git.
func initRepo(dir string) error {
	for _, args := range [][]string{
		{"init", "-q"},
		{"add", "."},
		{"commit", "-q", "-m", "Initial commit"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git %s: %s", args[0], strings.TrimSpace(string(out)))
		}
	}
	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, `create-tui — scaffold a terminal application from %s

usage:
  create-tui [flags] <directory>

flags:
`, sourceModule)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
examples:
  create-tui myapp
  create-tui -module github.com/you/myapp -name myapp -y myapp
  create-tui -demo -ref v0.2.0 playground
`)
}
