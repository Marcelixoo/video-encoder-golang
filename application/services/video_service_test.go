package services_test

import (
	"encoder/application/repositories"
	"encoder/application/services"
	"encoder/domain"
	"encoder/framework/database"
	"log"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestVideoServiceDownload(t *testing.T) {
	db := prepare()
	video := newVideo()
	videoRepository := newVideoRepository(db)

	videoService := services.NewVideoService(
		"video-encoder-golang-test",
		videoRepository,
	)

	err := videoService.Download(video)

	require.Nil(t, err)

	err = videoService.Fragment(video)
	require.Nil(t, err)
}

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Could not load .env file.")
	}
}

func prepare() *gorm.DB {
	db := database.NewDbTest()

	defer db.Close()

	return db
}

func newVideo() *domain.Video {
	video := domain.NewVideo()

	video.ID = uuid.NewV4().String()
	video.FilePath = "convite.mp4"
	video.CreatedAt = time.Now()

	return video
}

func newVideoRepository(db *gorm.DB) *repositories.VideoRepositoryDb {
	return &repositories.VideoRepositoryDb{Db: db}
}
