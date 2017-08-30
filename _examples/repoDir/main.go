package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/src-d/go-billy.v3/helper/chroot"
	"gopkg.in/src-d/go-git.v4"
	. "gopkg.in/src-d/go-git.v4/_examples"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

// Open an existing repository in a specific folder.
func main() {
	CheckArgs("<path>")
	path := os.Args[1]

	// We instance a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen(path)
	CheckIfError(err)

	p, err := root(r)
	CheckIfError(err)

	// Clean up the path
	if !filepath.IsAbs(p) {
		pwd, err := os.Getwd()
		CheckIfError(err)
		p = filepath.Join(pwd, p)
	}

	fmt.Println(filepath.Clean(p))
}

func root(r *git.Repository) (string, error) {
	// Try to grab the repository Storer
	s, ok := r.Storer.(*filesystem.Storage)
	if !ok {
		return "", errors.New("Repository storage is not filesystem.Storage")
	}

	// Try to get the underlying billy.Filesystem
	fs, ok := s.Filesystem().(*chroot.ChrootHelper)
	if !ok {
		return "", errors.New("Filesystem is not chroot.ChrootHelper")
	}

	return fs.Root(), nil
}
