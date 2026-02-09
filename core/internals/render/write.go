package render

import (
	"os"
	"path/filepath"
)

func WriteFile(path string, content string) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0o644)
}
