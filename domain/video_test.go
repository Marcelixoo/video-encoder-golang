package domain_test

import (
	"encoder/domain"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {
	video := domain.NewVideo() // convention to create objects from a pkg
	err := video.Validate()

	require.Error(t, err)
}

func TestVideoIDIsNotUuid(t *testing.T) {
	video := domain.NewVideo() // convention to create objects from a pkg

	video.ID = "abc" // not valid UUID
	video.ResourceID = "any-example-id"
	video.FilePath = "path-to-video-file"
	video.CreatedAt = time.Now()

	err := video.Validate()

	require.Error(t, err)
}

func TestVideoValidation(t *testing.T) {
	video := domain.NewVideo() // convention to create objects from a pkg

	video.ID = uuid.NewV4().String()
	video.ResourceID = "any-example-id"
	video.FilePath = "path-to-video-file"
	video.CreatedAt = time.Now()

	err := video.Validate()

	require.Nil(t, err)
}
