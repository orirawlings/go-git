package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	. "gopkg.in/src-d/go-git.v4/_examples"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Print log for commits with changes to certain file path
func main() {
	CheckArgs("<repoDir> <path>")
	repoDir := os.Args[1]
	path := os.Args[2]

	// We open the repository at given directory
	r, err := git.PlainOpen(repoDir)
	CheckIfError(err)

	Info(fmt.Sprintf("git log -- %s", path))

	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)

	// ... retrieves the commit history
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	CheckIfError(err)

	// ... just iterates over the commits, printing it
	err = cIter.ForEach(filterByChangesToPath(r, path, func(c *object.Commit) error {
		fmt.Println(c)
		return nil
	}))
	CheckIfError(err)
}

type memo map[plumbing.Hash]plumbing.Hash

// filterByChangesToPath provides a CommitIter callback that only
// invokes a delegate callback for commits that include changes to
// the content of path.
func filterByChangesToPath(r *git.Repository, path string, callback func(*object.Commit) error) func(*object.Commit) error {
	m := make(memo)
	return func(c *object.Commit) error {
		if err := ensure(m, c, path); err != nil {
			return err
		}
		if c.NumParents() == 0 && !m[c.Hash].IsZero() {
			// c is a root commit containing the path
			return callback(c)
		}
		// Compare the path in c with the path in each of its parents
		for _, p := range c.ParentHashes {
			if _, ok := m[p]; !ok {
				pc, err := r.CommitObject(p)
				if err != nil {
					return err
				}
				if err := ensure(m, pc, path); err != nil {
					return err
				}
			}
			if m[p] != m[c.Hash] {
				// contents at path are different from parent
				return callback(c)
			}
		}
		return nil
	}
}

// ensure our memoization includes a mapping from commit hash to
// the hash of path contents.
func ensure(m memo, c *object.Commit, path string) error {
	if _, ok := m[c.Hash]; !ok {
		t, err := c.Tree()
		if err != nil {
			return err
		}
		te, err := t.FindEntry(path)
		if err == object.ErrDirectoryNotFound {
			m[c.Hash] = plumbing.ZeroHash
			return nil
		} else if err != nil {
			if !strings.ContainsRune(path, '/') {
				// path is in root directory of project, but not found in this commit
				m[c.Hash] = plumbing.ZeroHash
				return nil
			}
			return err
		}
		m[c.Hash] = te.Hash
	}
	return nil
}
