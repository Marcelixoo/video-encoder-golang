package repositories_test

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/database"
	"testing"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestRepositoryDbInsert(t *testing.T) {
	db, err := database.NewDbTest().Connect()
	if err != nil {
		t.Fatalf("could not establish connection to db %v", err)
	}
	defer db.Close()

	job, video, err := newJob(db, "output_path", "Pending")
	require.Nil(t, err)

	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Insert(job)
	j, err := repoJob.Find(job.ID)

	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, j.ID, job.ID)
	require.Equal(t, j.VideoID, video.ID)
}

func TestRepositoryDbUpdate(t *testing.T) {
	db, err := database.NewDbTest().Connect()
	if err != nil {
		t.Fatalf("could not establish connection to db %v", err)
	}
	defer db.Close()

	job, _, err := newJob(db, "output_path", "Pending")
	require.Nil(t, err)

	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Insert(job)
	job.Status = "Complete"
	repoJob.Update(job)
	j, err := repoJob.Find(job.ID)

	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, j.Status, job.Status)
}

func newVideo() *domain.Video {
	return domain.NewVideo(uuid.NewV4().String(), "test-video", "path")
}

func newJob(db *gorm.DB, output string, status string) (*domain.Job, *domain.Video, error) {
	video := newVideo()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	job, err := domain.NewJob(video)

	return job, video, err
}
