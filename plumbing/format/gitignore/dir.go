package gitignore

import (
	"io/ioutil"
	"strings"

	"gopkg.in/src-d/go-billy.v2"
)

const (
	commentPrefix = "#"
	eol           = "\n"
	gitDir        = ".git"
	gitignoreFile = ".gitignore"
)

// ReadPatterns reads gitignore patterns recursively traversing through the directory
// structure. The result is in the ascending order of priority (last higher).
func ReadPatterns(fs billy.Filesystem, path []string) (ps []Pattern, err error) {
	if f, err := fs.Open(fs.Join(append(path, gitignoreFile)...)); err == nil {
		defer f.Close()
		if data, err := ioutil.ReadAll(f); err == nil {
			for _, s := range strings.Split(string(data), eol) {
				if !strings.HasPrefix(s, commentPrefix) && len(strings.TrimSpace(s)) > 0 {
					ps = append(ps, ParsePattern(s, path))
				}
			}
		}
	}

	var fis []billy.FileInfo
	fis, err = fs.ReadDir(fs.Join(path...))
	if err != nil {
		return
	}
	for _, fi := range fis {
		if fi.IsDir() && fi.Name() != gitDir {
			var subps []Pattern
			subps, err = ReadPatterns(fs, append(path, fi.Name()))
			if err != nil {
				return
			}
			if len(subps) > 0 {
				ps = append(ps, subps...)
			}
		}
	}

	return
}
