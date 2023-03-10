package services

import (
	"context"
	"encoder/application/repositories"
	"encoder/domain"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
)

type VideoService struct {
	VideoRepository repositories.VideoRepository
	Bucket          string
}

func NewVideoService(bucket string, repository repositories.VideoRepository) VideoService {
	return VideoService{
		VideoRepository: repository,
		Bucket:          bucket,
	}
}

func (v *VideoService) Download(video *domain.Video) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)

	if err != nil {
		return err
	}

	bucket := client.Bucket(v.Bucket)
	object := bucket.Object(video.FilePath)

	r, err := object.NewReader(ctx)
	if err != nil {
		return err
	}

	defer r.Close()

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	filename := absPathToLocalStorage(video.ID + ".mp4")
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	defer f.Close()

	log.Printf("Video %v stored", video.ID)

	return nil
}

func (v *VideoService) Fragment(video *domain.Video) error {
	fragmentsDir := absPathToLocalStorage(video.ID)

	err := os.Mkdir(fragmentsDir, os.ModePerm)
	if err != nil {
		return err
	}

	// necessary fragmentation step to prepare for slicing
	source := absPathToLocalStorage(video.ID + ".mp4")
	target := absPathToLocalStorage(video.ID + ".frag")

	cmd := exec.Command("mp4fragment", source, target)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func printOutput(output []byte) {
	if len(output) > 0 {
		log.Printf("=====> Output: %s\n", string(output))
	}
}

func absPathToLocalStorage(partToConcat string) string {
	localStoragePath := os.Getenv("LOCAL_STORAGE_PATH")
	return localStoragePath + "/" + partToConcat
}
