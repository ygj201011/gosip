// Package runtime provides helpers to get and format call stack frame.
package runtime

import (
	"fmt"
	"go/build"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

const FmtFull = "%s (%s:%s)"
const FmtUndef = "???"

type Frame struct {
	File string
	Line int
	Func *runtime.Func
}

func (frm *Frame) String() string {
	if frm == nil {
		return FmtUndef
	}

	var file, line, fn string

	if frm.Func == nil {
		fn = FmtUndef
	} else {
		fn = filepath.Base(frm.Func.Name())
	}

	if frm.File == "" {
		file = FmtUndef
	} else {
		file = cleanPath(frm.File)
	}

	if frm.Line == 0 {
		line = FmtUndef
	} else {
		line = strconv.Itoa(frm.Line)
	}

	return fmt.Sprintf(FmtFull, fn, file, line)
}

// GetFrame returns call frame where GetFrame was called.
func GetFrame() (*Frame, bool) {
	return getFrame(0)
}

// GetFrameOffset works like GetFrame but with custom frames offset.
func GetFrameOffset(offset int) (*Frame, bool) {
	return getFrame(offset)
}

func getFrame(offset int) (*Frame, bool) {
	pc, file, line, ok := runtime.Caller(offset + 2)
	if !ok {
		return nil, false
	}

	return &Frame{file, line, runtime.FuncForPC(pc)}, true
}

func cleanPath(path string) string {
	dirs := filepath.SplitList(build.Default.GOPATH)
	// Sort in decreasing order by length so the longest matching prefix is removed
	sort.Stable(longestFirst(dirs))
	for _, dir := range dirs {
		srcdir := filepath.Join(dir, "src")
		rel, err := filepath.Rel(srcdir, path)
		// filepath.Rel can traverse parent directories, don't want those
		if err == nil && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return rel
		}
	}
	return path
}

type longestFirst []string

func (strs longestFirst) Len() int           { return len(strs) }
func (strs longestFirst) Less(i, j int) bool { return len(strs[i]) > len(strs[j]) }
func (strs longestFirst) Swap(i, j int)      { strs[i], strs[j] = strs[j], strs[i] }
