package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

const fileName = "config.json"

type Config struct {
	Theme string `json:"theme"`
}

func Defaults() Config {
	return Config{Theme: "violet-dark"}
}

func Path(appName string) (string, error) {
	return xdg.ConfigFile(filepath.Join(appName, fileName))
}

func Load(path string) (Config, error) {
	raw, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Defaults(), nil
	}
	if err != nil {
		return Defaults(), err
	}
	cfg := Defaults()
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return Defaults(), err
	}
	return cfg, nil
}

func Save(path string, cfg Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(raw, '\n'), 0o644)
}
