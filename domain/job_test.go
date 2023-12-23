package domain_test

import (
	"encoder/domain"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestNewJobWorks(t *testing.T) {
	video := domain.NewVideo(uuid.NewV4().String(), "resource-id", "path")

	job, err := domain.NewJob(video)

	require.NotNil(t, job)
	require.Nil(t, err)
}

func TestNewJobHasPendingStatus(t *testing.T) {
	video := domain.NewVideo(uuid.NewV4().String(), "resource-id", "path")

	job, err := domain.NewJob(video)

	require.Nil(t, err)
	require.Equal(t, job.Status, domain.JOB_STATUS_PENDING)
}
