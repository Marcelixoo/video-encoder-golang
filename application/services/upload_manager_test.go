package services_test

import (
	"encoder/application/services"
	"encoder/domain"
	"encoder/framework/filesystem"
	"encoder/framework/gcp"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestVideoUploadManager(t *testing.T) {
	video, videoService := encodeVideo(t)

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = "video-encoder-golang-test"
	videoUpload.VideoPath = filesystem.AbsPathToLocalStorage(video.ID)

	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(50, doneUpload)

	result := <-doneUpload
	require.Equal(t, result, "upload completed") // move this into a const, i.e., STATUS_COMPLETED

	err := videoService.Finish(video)
	require.Nil(t, err)
}

func encodeVideo(t *testing.T) (*domain.Video, services.VideoService) {
	t.Helper()

	var err error

	videoStorage, err := gcp.NewCloudStorageReader("video-encoder-golang-test")
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

	return video, videoService
}

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("error loading .env file")
	}
}
