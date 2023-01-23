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
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucketName string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)

	if err != nil {
		return err
	}

	bucket := client.Bucket(bucketName)
	object := bucket.Object(v.Video.FilePath)

	r, err := object.NewReader(ctx)
	if err != nil {
		return err
	}

	defer r.Close()

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	filename := os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".mp4"
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	defer f.Close()

	log.Printf("Video %v stored", v.Video.ID)

	return nil
}

func (v *VideoService) Fragment() error {
	localStoragePath := os.Getenv("LOCAL_STORAGE_PATH")
	fragmentsDir := localStoragePath + "/" + v.Video.ID

	err := os.Mkdir(fragmentsDir, os.ModePerm)
	if err != nil {
		return err
	}

	// necessary fragmentation step to prepare for slicing
	source := localStoragePath + "/" + v.Video.ID + ".mp4"
	target := localStoragePath + "/" + v.Video.ID + ".frag"

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
