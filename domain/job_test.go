package domain_test

import (
	"encoder/domain"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestNewJob(t *testing.T) {
	video := domain.NewVideo()

	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	job, err := domain.NewJob("path", "Converted", video)

	require.NotNil(t, job)
	require.Nil(t, err)
}
