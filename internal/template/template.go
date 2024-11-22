package template

import (
	"os"
)

// templateExists checks if there is a template folder
func FolderExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func CopyContent(src, dst string) error {
	srcFS := os.DirFS(src)
	return os.CopyFS(dst, srcFS)
}
