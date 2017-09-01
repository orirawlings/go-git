package main

import (
	"os"

	"gopkg.in/src-d/go-git.v4"
	. "gopkg.in/src-d/go-git.v4/_examples"
	"gopkg.in/src-d/go-git.v4/config"
)

// Example of a push with non-default options and progress reporting
func main() {
	CheckArgs("<repository-path>", "<remote>", "<refspec> [<refspec> ...]")
	path := os.Args[1]
	remote := os.Args[2]

	var refSpecs []config.RefSpec
	for _, s := range os.Args[3:] {
		refSpecs = append(refSpecs, config.RefSpec(s))
	}

	r, err := git.PlainOpen(path)
	CheckIfError(err)

	err = r.Push(&git.PushOptions{
		RemoteName: remote,
		RefSpecs:   refSpecs,
		Progress:   os.Stdout,
	})
	CheckIfError(err)
}
