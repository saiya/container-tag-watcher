package watcher_test

import (
	"fmt"
	"sync"
)

type stubClient struct {
	lock   sync.Mutex
	hashes map[string]string
}

func newStubClient() *stubClient {
	return &stubClient{
		hashes: map[string]string{},
	}
}

func (c *stubClient) SetImageDigest(platform string, imageReference string, hash string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.hashes[fmt.Sprintf("%s:%s", platform, imageReference)] = hash
}

func (c *stubClient) GetImageDigest(platform string, imageReference string) (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.hashes[fmt.Sprintf("%s:%s", platform, imageReference)], nil
}
