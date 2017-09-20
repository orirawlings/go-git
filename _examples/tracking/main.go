package main

import (
	"fmt"
	"os"

	"gopkg.in/src-d/go-git.v4"
	. "gopkg.in/src-d/go-git.v4/_examples"
	"gopkg.in/src-d/go-git.v4/plumbing/format/config"
)

// Print the tracking branch configuration for the repository
func main() {
	CheckArgs("<path>")
	path := os.Args[1]

	r, err := git.PlainOpen(path)
	CheckIfError(err)

	c, err := r.Config()
	CheckIfError(err)

	// There is no explicit library support for branch config sections so we will work with the raw parsed config file
	rc := c.Raw

	// Relevant config is in the "branch" section
	bs, err := section(rc, "branch")
	CheckIfError(err)

	// Each branch has its configuration in a separate subsection
	for _, b := range bs.Subsections {
		// Configuration is split into two separate options:
		//   the remote containing the tracked branch, "remote"
		//   the reference name of the tracked branch in the remote, "merge"
		// If the branch is tracking a local branch, "remote" will be "."
		fmt.Printf("%s\tremote=%q\tmerge=%q\n", b.Name, b.Option("remote"), b.Option("merge"))
	}
}

// section returns the section with name in the given config
func section(rc *config.Config, name string) (*config.Section, error) {
	for _, s := range rc.Sections {
		if s.Name == name {
			return s, nil
		}
	}
	return nil, fmt.Errorf("section %q not found", name)
}
