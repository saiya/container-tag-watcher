package containerregistryclient

import (
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/saiya/container_tag_watcher/logger"
)

type Settings struct {
	EnableAwsEcrSupport bool
}

func keyChain(s *Settings) remote.Option {
	keyChains := []authn.Keychain{authn.DefaultKeychain}

	if s.EnableAwsEcrSupport {
		logger.Get().Debugw("Enabling AWS ECR support...")
		authn.NewKeychainFromHelper(ecr.NewECRHelper())
	}

	return remote.WithAuthFromKeychain(authn.NewMultiKeychain(keyChains...))
}
