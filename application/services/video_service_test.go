package services_test

import (
	"encoder/application/repositories"
	"encoder/application/services"
	"encoder/domain"
	"encoder/framework/database"
	"encoder/framework/gcp"
	"fmt"
	"log"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestVideoServiceDownload(t *testing.T) {
	var err error

	db, err := prepare()
	if err != nil {
		t.Fatal(err)
	}

	video := newVideo()
	videoRepository := newVideoRepository(db)

	videoStorage, err := gcp.NewCloudStorageReader("video-encoder-golang-test")
	if err != nil {
		t.Fatal(err)
	}

	videoService := services.NewVideoService(
		videoRepository,
		videoStorage,
	)

	err = videoService.Download(video)
	require.Nil(t, err)

	err = videoService.Fragment(video)
	require.Nil(t, err)

	err = videoService.Encode(video)
	require.Nil(t, err)

	err = videoService.Finish(video)
	require.Nil(t, err)
}

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Could not load .env file.")
	}
}

func prepare() (*gorm.DB, error) {
	db, err := database.NewDbTest().Connect()
	if err != nil {
		return nil, fmt.Errorf("could not establish connection to db %v", err)
	}
	defer db.Close()

	return db, nil
}

func newVideo() *domain.Video {
	video := domain.NewVideo(uuid.NewV4().String(), "resource-id", "convite.mp4")

	fmt.Printf("testing with video %v", video)

	return video
}

func newVideoRepository(db *gorm.DB) *repositories.VideoRepositoryDb {
	return &repositories.VideoRepositoryDb{Db: db}
}
