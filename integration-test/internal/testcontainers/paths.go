package testcontainers

import (
	"path/filepath"
	"runtime"
)

func RepoRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..", ".."))
}

func Dockerfile(name string) string {
	return filepath.Join(RepoRoot(), "integration-test", "docker", name)
}

func PostgresInitScript() string {
	return filepath.Join(RepoRoot(), "integration-test", "fixtures", "postgres", "init.sql")
}
