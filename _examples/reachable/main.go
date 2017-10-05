package main

import (
	"fmt"
	"os"

	"gopkg.in/src-d/go-git.v4"
	. "gopkg.in/src-d/go-git.v4/_examples"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// Prints out the name of each branch from which given commit can be reached.
func main() {
	CheckArgs("<path>", "<commit>")
	path := os.Args[1]
	commit := plumbing.NewHash(os.Args[2])

	r, err := git.PlainOpen(path)
	CheckIfError(err)

	memo := make(map[plumbing.Hash]bool)

	rs, err := r.References()
	CheckIfError(err)

	CheckIfError(rs.ForEach(func(ref *plumbing.Reference) error {
		n := ref.Name()
		if n.IsBranch() {
			b, err := r.Reference(n, true)
			if err != nil {
				return err
			}
			v, err := reaches(r, b.Hash(), commit, memo)
			if err != nil {
				return err
			}
			if v {
				fmt.Println(n)
			}
		}
		return nil
	}))
}

// reaches returns true if commit, c, can be reached from commit, start. Results are memoized in memo.
func reaches(r *git.Repository, start, c plumbing.Hash, memo map[plumbing.Hash]bool) (bool, error) {
	if v, ok := memo[start]; ok {
		return v, nil
	}
	if start == c {
		memo[start] = true
		return true, nil
	}
	co, err := r.CommitObject(start)
	if err != nil {
		return false, err
	}
	for _, p := range co.ParentHashes {
		v, err := reaches(r, p, c, memo)
		if err != nil {
			return false, err
		}
		if v {
			memo[start] = true
			return true, nil
		}
	}
	memo[start] = false
	return false, nil
}
