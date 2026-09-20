package main

import (
	"os"
	"path/filepath"
)

// demoModules are the example modules. They exist to exercise every pattern
// end to end; a project that is not a demo starts without them.
var demoModules = []string{"home", "manager", "settings"}

// demoArtifacts are the files outside internal/modules that only make sense
// while the demo modules are present: the golden frame and the tests that
// navigate between demo routes.
var demoArtifacts = []string{
	"internal/app/app_test.go",
	"internal/app/testdata",
}

// prune replaces the demo with a single module named after the project, so
// that the scaffolded tree still shows how a screen is registered without
// carrying three screens nobody asked for.
//
// It runs before rewrite, which is why the files it generates still carry the
// source module path.
func prune(o options) error {
	for _, name := range demoModules {
		if err := os.RemoveAll(filepath.Join(o.Dir, "internal", "modules", name)); err != nil {
			return err
		}
	}
	for _, rel := range demoArtifacts {
		if err := os.RemoveAll(filepath.Join(o.Dir, filepath.FromSlash(rel))); err != nil {
			return err
		}
	}

	files := map[string]string{
		"internal/modules/hello/module.go": helloModuleGo,
		"internal/modules/hello/view.go":   helloViewGo,
		"internal/app/register.go":         registerGo,
		"internal/app/app_test.go":         appTestGo,
	}
	for rel, content := range files {
		if err := write(o.Dir, rel, content); err != nil {
			return err
		}
	}

	// The end-to-end test waits for a string the demo used to paint.
	return replaceInFile(filepath.Join(o.Dir, "internal", "app", "program_test.go"),
		`[]byte("Overview")`, `[]byte("Hello")`,
	)
}

func write(dir, rel, content string) error {
	target := filepath.Join(dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	return os.WriteFile(target, []byte(content), 0o644)
}
