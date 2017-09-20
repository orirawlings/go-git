package main

import (
	"fmt"
	"os"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/client"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"gopkg.in/src-d/go-git.v4/utils/ioutil"

	. "gopkg.in/src-d/go-git.v4/_examples"
)

// List remote references
func main() {
	CheckArgs("<path>", "<remote>")
	path := os.Args[1]
	remoteName := os.Args[2]

	r, err := git.PlainOpen(path)
	CheckIfError(err)

	remote, err := r.Remote(remoteName)
	CheckIfError(err)

	Info("git ls-remote %s", remoteName)

	// You can pass in a custom auth for the transport if necessary
	ar, err := lsRemote(remote, nil)
	CheckIfError(err)

	refs, err := ar.IterReferences()
	CheckIfError(err)

	refs.ForEach(func(ref *plumbing.Reference) error {
		fmt.Println(ref)
		return nil
	})
	CheckIfError(err)
}

// lsRemote returns the references contained in the remote
func lsRemote(remote *git.Remote, auth transport.AuthMethod) (memory.ReferenceStorage, error) {
	url := remote.Config().URLs[0]
	s, err := newUploadPackSession(url, auth)
	if err != nil {
		return nil, err
	}
	defer ioutil.CheckClose(s, &err)

	ar, err := s.AdvertisedReferences()
	if err != nil {
		return nil, err
	}

	return ar.AllReferences()
}

func newUploadPackSession(url string, auth transport.AuthMethod) (transport.UploadPackSession, error) {
	c, ep, err := newClient(url)
	if err != nil {
		return nil, err
	}

	return c.NewUploadPackSession(ep, auth)
}

func newClient(url string) (transport.Transport, transport.Endpoint, error) {
	ep, err := transport.NewEndpoint(url)
	if err != nil {
		return nil, nil, err
	}

	c, err := client.NewClient(ep)
	if err != nil {
		return nil, nil, err
	}

	return c, ep, err
}
