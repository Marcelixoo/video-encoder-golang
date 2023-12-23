package domain_test

import (
	"encoder/domain"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestVideoIDIsNotUuid(t *testing.T) {
	invalidID := "abc"

	video := domain.NewVideo(invalidID, "any-example-id", "path-to-video-file")

	require.Error(t, video.Validate())
}

func TestVideoValidation(t *testing.T) {
	video := domain.NewVideo(uuid.NewV4().String(), "any-example-id", "path-to-video-file")

	err := video.Validate()

	require.Nil(t, err)
}

func TestVideoValidationFailureWhenNoResourceIDProvided(t *testing.T) {
	video := domain.Video{
		FilePath: "/tmp/123123123",
	}

	err := video.Validate()

	require.NotNil(t, err)
}

func TestVideoValidationGeneratesUUIDIfNotProvided(t *testing.T) {
	video := domain.Video{
		ResourceID: "123123123",
		FilePath:   "/tmp/123123123",
	}

	err := video.Validate()

	require.Nil(t, err)
}
