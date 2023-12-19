package services

import (
	"context"
	"encoder/application/repositories"
	"encoder/domain"
	"fmt"
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

// Encode packages a local .mp4 video file
// into a series of fragments using the
// command line tool `bento4`.
func (v *VideoService) Encode(video *domain.Video) error {
	intermediaryFilename := fmt.Sprintf("%s.frag", video.ID)
	destinationFolder := fmt.Sprint(video.ID)

	cmdArgs := []string{}

	cmdArgs = append(cmdArgs, absPathToLocalStorage(intermediaryFilename))
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, absPathToLocalStorage(destinationFolder))
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")

	cmd := exec.Command("mp4dash", cmdArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

// Finish does a cleanup of files created
// during intermediary steps of the final
// encoding a .mp4 video.
func (v *VideoService) Finish(video *domain.Video) error {
	var (
		err error

		initialMP4Filename string = video.ID + ".mp4"
		fragmentedFilename string = video.ID + ".frag"
		outputFolder       string = video.ID
	)

	err = os.Remove(absPathToLocalStorage(initialMP4Filename))
	if err != nil {
		log.Printf("error removing initial .mp4 file %s \n", initialMP4Filename)
		return err
	}

	err = os.Remove(absPathToLocalStorage(fragmentedFilename))
	if err != nil {
		log.Printf("error removing fragmented file %s \n", fragmentedFilename)
		return err
	}

	err = os.RemoveAll(absPathToLocalStorage(outputFolder))
	if err != nil {
		log.Printf("error removing folder %s \n", outputFolder)
		return err
	}

	log.Printf("all intermediary files generated for video %s have been removed \n", video.ID)

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
