package fileManager

import (
	"io/fs"
	"os"
	"path/filepath"
)

type FileManager struct {
}

func (f *FileManager) List(root string) []string {
	res := []string{}

	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		res = append(res, path)

		return nil
	})

	return res
}

func (f *FileManager) Rename(old, new string) error {
	if err := os.Rename(old, new); err != nil {
		return err
	}

	return nil
}
