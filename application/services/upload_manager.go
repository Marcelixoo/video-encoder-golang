package services

import (
	"context"
	"encoder/framework/filesystem"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"cloud.google.com/go/storage"
)

type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{
		Paths:  []string{},
		Errors: []string{},
	}
}

func (vu *VideoUpload) UploadObject(
	objectPath string,
	client *storage.Client,
	ctx context.Context,
) error {
	filename := strings.Split(objectPath, filesystem.AbsPathToLocalStorage(""))[1]

	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}
	defer f.Close()

	object := client.Bucket(vu.OutputBucket).Object(filename)
	object.If(storage.Conditions{DoesNotExist: true})

	wc := object.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err = wc.Close(); err != nil {
		return err
	}

	return nil
}

const (
	UPLOAD_STATUS_COMPLETED = "upload completed"
	UPLOAD_STATUS_RUNNING   = "upload is running"
	UPLOAD_STATUS_FAILED    = "upload failed"
)

func (vu *VideoUpload) ProcessUpload(concurrency int, doneUpload chan string) error {
	// stream the index of files to be
	// processed from the slice of paths.
	in := make(chan int, runtime.NumCPU())
	returnChannel := make(chan string)

	err := vu.loadPaths()
	if err != nil {
		return err
	}

	uploadClient, ctx, err := getClientUpload()
	if err != nil {
		return err
	}

	for process := 0; process < concurrency; process++ {
		go vu.uploadWorker(in, returnChannel, uploadClient, ctx)
	}

	go func() {
		for x := 0; x < len(vu.Paths); x++ {
			in <- x
		}
	}()

	countDoneWorker := 0
	for r := range returnChannel {
		countDoneWorker++

		if r != UPLOAD_STATUS_RUNNING {
			doneUpload <- r
			break
		}

		if countDoneWorker == len(vu.Paths) {
			close(in)
		}
	}

	return nil
}

func (vu *VideoUpload) loadPaths() error {
	return filepath.Walk(vu.VideoPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			vu.Paths = append(vu.Paths, path)
		}

		return nil
	})
}

func (vu *VideoUpload) uploadWorker(
	in chan int,
	returnChan chan string,
	uploadClient *storage.Client,
	ctx context.Context,
) {

	for x := range in {
		err := vu.UploadObject(vu.Paths[x], uploadClient, ctx)

		if err != nil {
			vu.Errors = append(vu.Errors, vu.Paths[x])
			log.Printf("error during upload of %s: %v", vu.Paths[x], err)
			returnChan <- UPLOAD_STATUS_FAILED
		}

		returnChan <- UPLOAD_STATUS_RUNNING
	}

	returnChan <- UPLOAD_STATUS_COMPLETED
}

func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}
