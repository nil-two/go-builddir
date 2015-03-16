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

type Items []Item

func (i *Items) Build(root string) error {
	for _, i := range []Item(*i) {
		err := i.Build(root)
		if err != nil {
			return err
		}
	}
	return nil
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
	Content Items
}

func (d *Dir) Build(root string) error {
	path := filepath.Join(root, d.Name)
	if !exist(path) {
		if err := os.Mkdir(path, 0755); err != nil {
			return err
		}
	}
	return d.Content.Build(path)
}
