package containerregistryclient_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetImageDigest(t *testing.T) {
	for _, platform := range []string{"linux/amd64", "linux/arm64"} {
		t.Run(platform, func(t *testing.T) {
			hash, err := client.GetImageDigest(platform, "saiya/dsps:latest")
			assert.NoError(t, err)
			assert.NotEmpty(t, hash)
		})
	}
}
