package builddir

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func exist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

type Item interface {
	Build(root string) error
}

type File struct {
	Name string
	Content []byte
}

func (f *File) Build(root string) error {
	path := filepath.Join(root, f.Name)
	if exist(path) {
		return nil
	}
	return ioutil.WriteFile(path, f.Content, 0644)
}

type Dir struct {
	Name string
	Content []Item
}

func (d *Dir) Build(root string) error {
	path := filepath.Join(root, d.Name)
	if !exist(path) {
		if err := os.Mkdir(path, 0755); err != nil {
			return err
		}
	}
	for _, i := range d.Content {
		err := i.Build(path)
		if err != nil {
			return err
		}
	}
	return nil
}
