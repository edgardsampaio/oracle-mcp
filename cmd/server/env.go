package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"oracle-mcp/internal/config"
)

func prepareRuntime() error {
	dir, err := executableDir()
	if err != nil {
		return err
	}
	loadDotEnv(filepath.Join(dir, ".env"))
	for _, key := range []string{"ORACLE_MCP_CONFIG", "ORACLE_MCP_SQLITE_PATH"} {
		rebase(key, dir)
	}
	return nil
}

func executableDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		name, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		name = strings.TrimSpace(name)
		if name == "" || os.Getenv(name) != "" {
			continue
		}
		os.Setenv(name, value)
	}
}

func rebase(key, dir string) {
	value := os.Getenv(key)
	if value == "" || filepath.IsAbs(value) {
		return
	}
	os.Setenv(key, filepath.Join(dir, value))
}

func besideExecutable(dir, path string) string {
	if path == "" || filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(dir, path)
}

func applyRuntimePaths(cfg *config.Config) error {
	dir, err := executableDir()
	if err != nil {
		return err
	}
	cfg.Auth.SQLitePath = besideExecutable(dir, cfg.Auth.SQLitePath)
	return nil
}
