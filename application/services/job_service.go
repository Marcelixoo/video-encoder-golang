package services

import (
	"encoder/application/repositories"
	"encoder/domain"
	"errors"
	"os"
	"strconv"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepository
	VideoService  VideoService
}

func NewJobService(jobRepository repositories.JobRepository, videoService VideoService) *JobService {
	return &JobService{
		Job:           &domain.Job{},
		JobRepository: jobRepository,
		VideoService:  videoService,
	}
}

func (j *JobService) Start(video *domain.Video) error {
	err := j.changeJobStatus("DOWNLOADING")
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Download(video)
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("FRAGMENTING")
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Fragment(video)
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("ENCODING")
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Encode(video)
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("UPLOADING")
	if err != nil {
		return j.failJob(err)
	}

	err = j.performUpload(video)
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("FINISHING")
	if err != nil {
		return j.failJob(err)
	}

	err = j.VideoService.Finish(video)
	if err != nil {
		return j.failJob(err)
	}

	err = j.changeJobStatus("COMPLETED")
	if err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) performUpload(video *domain.Video) error {
	videoUpload := NewVideoUpload()
	videoUpload.OutputBucket = os.Getenv("OUTPUT_BUCKET_NAME")
	videoUpload.VideoPath = os.Getenv("LOCAL_STORAGE_PATH") + "/" + video.ID

	concurrency, _ := strconv.Atoi(os.Getenv("CONCURRENCY_UPLOAD"))
	doneUpload := make(chan string)

	go videoUpload.ProcessUpload(concurrency, doneUpload)

	uploadResult := <-doneUpload
	if uploadResult != UPLOAD_STATUS_COMPLETED {
		return j.failJob(errors.New(uploadResult))
	}
	return nil
}

func (j *JobService) changeJobStatus(newStatus string) error {
	var err error

	j.Job.Status = newStatus

	j.Job, err = j.JobRepository.Update(j.Job)
	if err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) failJob(prevErr error) error {
	j.Job.Status = "FAILED"
	j.Job.Error = prevErr.Error()

	_, err := j.JobRepository.Update(j.Job)
	if err != nil {
		return err
	}

	return prevErr
}
