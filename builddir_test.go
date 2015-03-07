package builddir_test

import (
	. "github.com/kusabashira/go-builddir"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var dir = &Dir{
	Name: "test",
	Content: []Item{
		&File{
			Name: "file_a",
			Content: []byte("content_a\n"),
		},
		&File{
			Name: "file_b",
			Content: []byte("content_b\n"),
		},
		&Dir{
			Name: "dir_a",
		},
		&Dir{
			Name: "dir_b",
			Content: []Item{
				&File{
					Name: "file_c",
					Content: []byte("content_c\n"),
				},
			},
		},
	},
}

func exist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func TestSetupDir(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed load working dir %v:", err)
	}
	defer os.RemoveAll("test")
	if err = dir.Build(wd); err != nil {
		t.Fatalf("failed create dir %v:", err)
	}

	var paths = []string{
		filepath.Join(wd, "test"),
		filepath.Join(wd, "test", "file_a"),
		filepath.Join(wd, "test", "file_b"),
		filepath.Join(wd, "test", "dir_a"),
		filepath.Join(wd, "test", "dir_b"),
		filepath.Join(wd, "test", "dir_b", "file_c"),
	}
	for _, path := range paths {
		if !exist(path) {
			t.Fatalf("%s doesn't exist:", path)
		}
	}
}

func TestSetupDirTwice(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed load working dir %v:", err)
	}
	defer os.RemoveAll("test")
	if err = dir.Build(wd); err != nil {
		t.Fatalf("failed 1st create dir %v:", err)
	}
	if err = dir.Build(wd); err != nil {
		t.Fatalf("failed 2nd create dir %v:", err)
	}
}

func TestNotOverWrite(t *testing.T) {
	wd, err := os.Getwd()
	if err = dir.Build(wd); err != nil {
		t.Fatalf("failed load working dir %v:", err)
	}
	defer os.RemoveAll("test")

	if err = dir.Build(wd); err != nil {
		t.Fatalf("failed 1st create dir %v:", err)
	}
	
	file_a := filepath.Join(wd, "test", "file_a")
	if err = ioutil.WriteFile(file_a, []byte("new\n"), 0644); err != nil {
		t.Fatal(err)
	}

	if err = dir.Build(wd); err != nil {
		t.Fatalf("failed 2nd create dir %v:", err)
	}

	a_content, err := ioutil.ReadFile(file_a)
	if err != nil {
		t.Fatal(err)
	}
	if string(a_content) != "new\n" {
		t.Fatalf("%v over writed by dir.Build:", file_a)
	}
}
