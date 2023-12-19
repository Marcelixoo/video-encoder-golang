package repositories_test

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/database"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestVideoRepositoryDbInsert(t *testing.T) {
	db, err := database.NewDbTest().Connect()
	if err != nil {
		t.Fatalf("could not establish connection to db %v", err)
	}
	defer db.Close()

	repo := repositories.VideoRepositoryDb{Db: db}
	video := domain.NewVideo(uuid.NewV4().String(), "resource-id", "path")
	repo.Insert(video)

	v, err := repo.Find(video.ID)

	require.NotEmpty(t, v.ID)
	require.Nil(t, err)
	require.Equal(t, v.ID, video.ID)
}
