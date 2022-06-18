package generate

import (
	"os"
	"path"
	"strings"
)

func getRelativePath(root, name string) string {
	relativePath := strings.TrimPrefix(name, root)
	relativePath = strings.TrimPrefix(relativePath, "/")

	return relativePath
}

func makeDirAll(name string) error {
	return os.MkdirAll(path.Dir(name), 0o755)
}
