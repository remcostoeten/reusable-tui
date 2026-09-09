package main

import (
	"errors"
	"flag"
	"io"
)

var errShowVersion = errors.New("version requested")

type flags struct {
	configPath string
	dbPath     string
	themeName  string
	noMouse    bool
}

func parseFlags(args []string, stderr io.Writer) (flags, error) {
	var f flags
	set := flag.NewFlagSet("tui", flag.ContinueOnError)
	set.SetOutput(stderr)
	set.StringVar(&f.configPath, "config", "", "path to config.json (default: XDG config dir)")
	set.StringVar(&f.dbPath, "db", "", "path to the SQLite database (default: XDG data dir)")
	set.StringVar(&f.themeName, "theme", "", "theme to start with, overriding the config for this run")
	set.BoolVar(&f.noMouse, "no-mouse", false, "disable mouse support")
	version := set.Bool("version", false, "print the version and exit")
	if err := set.Parse(args); err != nil {
		return f, err
	}
	if *version {
		return f, errShowVersion
	}
	return f, nil
}
