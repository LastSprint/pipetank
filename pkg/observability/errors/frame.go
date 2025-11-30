package errors

import (
	"fmt"
	"strings"
)

var _frameFileRootPath = ""

func SetFrameFileRootPath(val string) {
	_frameFileRootPath = val
}

type frame struct {
	file, fn string
	ln       int
}

func (f frame) String() string {
	return fmt.Sprintf("[%s:%d]::%s", f.fn, f.ln, filePath(f.file, _frameFileRootPath))
}

//nolint:unused
func fnName(str string) string {
	v := strings.Split(str, ".")
	if len(v) != 2 {
		return str
	}

	return v[1]
}

func filePath(path, root string) string {
	items := strings.Split(path, "/")
	newItems := make([]string, 0, len(items)/2)
	foundRoot := false
	for _, it := range items {
		if it == root {
			foundRoot = true
		}

		if !foundRoot {
			continue
		}

		newItems = append(newItems, it)
	}

	if len(newItems) == 0 {
		return path
	}

	return strings.Join(newItems, "/")
}
