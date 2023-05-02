package config

import (
	"path/filepath"
	"runtime"
)

var (
	// Get current file full path from runtime
	_, b, _, _ = runtime.Caller(0)

	// ProjectRootPath Root folder of this project
	ProjectRootPath = filepath.Join(filepath.Dir(b), "../")
	PathFromRoot    = func(path string) string {
		return ProjectRootPath + path
	}
)
