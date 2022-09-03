package parser

import (
	"fmt"
	"reflect"
	"testing"
)

func TestId(t *testing.T) {
	given := "00000104"
	want := 104
	got := id(given)

	if got != want {
		t.Fatalf("Wanted %d\nGot %d", want, got)
	}
}

type fileMgr struct {
	filenames []string
	results   map[string]string
}

func (f *fileMgr) List(root string) []string {
	return f.filenames
}

func (f *fileMgr) Rename(o, n string) error {
	f.results[o] = n
	return nil
}

func TestRename(t *testing.T) {
	given := []string{"some/dir/bday_001.jpg", "some/dir/bday_002.jpg", "other/dir/bday_003.jpg", "some/dir/school_01.txt"}

	mock := fileMgr{
		given,
		make(map[string]string),
	}

	want := map[string]string{
		"some/dir/bday_001.jpg":  "some/dir/Bday (1 of 3).jpg",
		"some/dir/bday_002.jpg":  "some/dir/Bday (2 of 3).jpg",
		"other/dir/bday_003.jpg": "other/dir/Bday (3 of 3).jpg",
		"some/dir/school_01.txt": "some/dir/School (1 of 1).txt",
	}

	Rename("someroot", &mock)

	if !reflect.DeepEqual(want, mock.results) {
		fmt.Printf("\n%+v", want)
		fmt.Printf("\n%+v", mock.results)
		t.Fatalf("failed")
	}
}
