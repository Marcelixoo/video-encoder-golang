package services_test

import (
	"encoder/application/repositories"
	"encoder/application/services"
	"encoder/domain"
	"encoder/framework/database"
	"encoder/framework/gcp"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestVideoServiceDownload(t *testing.T) {
	var err error

	bucketName := os.Getenv("OUTPUT_BUCKET_NAME")

	videoStorage, err := gcp.NewCloudStorageReader(bucketName)
	if err != nil {
		t.Fatal(err)
	}

	video, videoRepository := prepare()
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

func prepare() (*domain.Video, *repositories.VideoRepositoryDb) {
	db, err := database.NewDbTest().Connect()
	if err != nil {
		panic(fmt.Errorf("could not establish connection to db %v", err))
	}
	defer db.Close()

	video := newVideo()
	repository := newVideoRepository(db)

	return video, repository
}

func newVideo() *domain.Video {
	video := domain.NewVideo(uuid.NewV4().String(), "resource-id", "convite.mp4")

	fmt.Printf("testing with video %v", video)

	return video
}

func newVideoRepository(db *gorm.DB) *repositories.VideoRepositoryDb {
	return &repositories.VideoRepositoryDb{Db: db}
}
