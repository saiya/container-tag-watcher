package containerregistryclient

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/saiya/container_tag_watcher/logger"
)

func (c *client) GetImageDigest(platformName string, imageReference string) (string, error) {
	logAttrs := []interface{}{"platform", platformName, "imageReference", imageReference}

	platform, err := v1.ParsePlatform(platformName)
	if err != nil {
		return "", fmt.Errorf("failed to parse container platform specifier: %w", err)
	}

	nameRef, err := name.ParseReference(imageReference)
	if err != nil {
		return "", fmt.Errorf("failed to parse container image reference: %w", err)
	}

	logger.Get().Debugw("Fetching container image information", logAttrs...)
	img, err := remote.Image(nameRef, append(c.options, remote.WithPlatform(*platform))...)
	if err != nil {
		return "", fmt.Errorf("failed to get container image from registry: %w", err)
	}

	hash, err := img.Digest()
	if err != nil {
		return "", fmt.Errorf("cannot parse digest of image: %w", err)
	}

	hashStr := fmt.Sprintf("%s:%s", hash.Algorithm, hash.Hex)
	logAttrs = append(logAttrs, "hash", hashStr)
	logger.Get().Debugw("Fetched container image information", logAttrs...)
	return hashStr, nil
}
