package containerregistryclient

import "github.com/google/go-containerregistry/pkg/v1/remote"

type Client interface {
	GetImageDigest(platform string, imageReference string) (string, error)
}

type client struct {
	options []remote.Option
}

func Init(s *Settings) Client {
	return &client{
		options: []remote.Option{keyChain(s)},
	}
}
