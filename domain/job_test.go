package domain_test

import (
	"encoder/domain"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestNewJob(t *testing.T) {
	video := domain.NewVideo(uuid.NewV4().String(), "resource-id", "path")

	job, err := domain.NewJob("path", "Converted", video)

	require.NotNil(t, job)
	require.Nil(t, err)
}
