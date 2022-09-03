package parser

import (
	"fmt"
	"path/filepath"
	"strings"
)

// key = input string up to underscore, not including it. id = first integer AFTER underscore and BEFORE '.' character
func id(s string) int {
	total := 0
	countZero := false

	for _, ch := range s {
		if ch == '0' && !countZero {
			continue
		}

		countZero = true
		total *= 10
		addend := int(ch - '0')
		total += addend
	}

	return total
}

type fileMeta struct {
	extension  string
	baseName   string
	occurences int
}

type fileDescriptor struct {
	meta *fileMeta
	id   int
	dir  string
}

func processNames(names []string) *map[string]fileDescriptor {
	meta := make(map[string]*fileMeta)
	descriptors := make(map[string]fileDescriptor)

	for _, name := range names {
		dir, file := filepath.Split(name)

		fields := strings.FieldsFunc(file, func(r rune) bool {
			return r == '.' || r == '_'
		})

		if existing, ok := meta[fields[0]]; ok {
			existing.occurences++
			descriptors[name] = fileDescriptor{
				existing,
				id(fields[1]),
				dir,
			}
		} else {
			meta[fields[0]] = &fileMeta{
				fields[2],
				strings.Title(fields[0]),
				1,
			}

			v := meta[fields[0]]

			descriptors[name] = fileDescriptor{
				v,
				id(fields[1]),
				dir,
			}
		}
	}

	return &descriptors
}

func makeName(desc fileDescriptor) string {
	return fmt.Sprintf("%s%s (%d of %d).%s", desc.dir, desc.meta.baseName, desc.id, desc.meta.occurences, desc.meta.extension)
}

type FileManager interface {
	List(root string) []string
	Rename(old, new string) error
}

func Rename(root string, fmgr FileManager) {
	filenames := fmgr.List(root)

	fileMeta := *processNames(filenames)

	for _, name := range filenames {
		if m, ok := fileMeta[name]; ok {
			fmgr.Rename(name, makeName(m))
		} else {
			panic(fmt.Sprintf("File %s received from file manager, but not populated in file metadata", name))
		}
	}
}
