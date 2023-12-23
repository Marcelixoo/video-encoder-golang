package services

import (
	"encoder/application/repositories"
	"encoder/domain"
	"errors"
	"os"
	"strconv"
)

type JobService struct {
	JobRepository repositories.JobRepository
	VideoService  VideoService
}

func NewJobService(jobRepository repositories.JobRepository, videoService VideoService) *JobService {
	return &JobService{
		JobRepository: jobRepository,
		VideoService:  videoService,
	}
}

func (j *JobService) Start(job *domain.Job) error {
	err := j.changeJobStatus(job, domain.JOB_STATUS_DOWNLOADING)
	if err != nil {
		return j.failJob(job, err)
	}

	err = j.VideoService.Download(job.Video)
	if err != nil {
		return j.failJob(job, err)
	}

	err = j.changeJobStatus(job, domain.JOB_STATUS_FRAGMENTING)
	if err != nil {
		return j.failJob(job, err)
	}

	err = j.VideoService.Fragment(job.Video)
	if err != nil {
		return j.failJob(job, err)
	}

	err = j.changeJobStatus(job, domain.JOB_STATUS_ENCODING)
	if err != nil {
		return j.failJob(job, err)
	}

	err = j.VideoService.Encode(job.Video)
	if err != nil {
		return j.failJob(job, err)
	}

	err = j.changeJobStatus(job, domain.JOB_STATUS_UPLOADING)
	if err != nil {
		return j.failJob(job, err)
	}

	err = j.performUpload(job)
	if err != nil {
		return j.failJob(job, err)
	}

	err = j.changeJobStatus(job, domain.JOB_STATUS_FINISHING)
	if err != nil {
		return j.failJob(job, err)
	}

	err = j.VideoService.Finish(job.Video)
	if err != nil {
		return j.failJob(job, err)
	}

	err = j.changeJobStatus(job, domain.JOB_STATUS_COMPLETED)
	if err != nil {
		return j.failJob(job, err)
	}

	return nil
}

func (j *JobService) performUpload(job *domain.Job) error {
	videoUpload := NewVideoUpload()
	videoUpload.OutputBucket = os.Getenv("OUTPUT_BUCKET_NAME")
	videoUpload.VideoPath = os.Getenv("LOCAL_STORAGE_PATH") + "/" + job.Video.ID

	concurrency, _ := strconv.Atoi(os.Getenv("CONCURRENCY_UPLOAD"))
	doneUpload := make(chan string)

	go videoUpload.ProcessUpload(concurrency, doneUpload)

	uploadResult := <-doneUpload
	if uploadResult != UPLOAD_STATUS_COMPLETED {
		return j.failJob(job, errors.New(uploadResult))
	}
	return nil
}

func (j *JobService) changeJobStatus(job *domain.Job, newStatus string) error {
	var err error

	job.Status = newStatus

	job, err = j.JobRepository.Update(job)
	if err != nil {
		return j.failJob(job, err)
	}

	return nil
}

func (j *JobService) failJob(job *domain.Job, prevErr error) error {
	job.Status = domain.JOB_STATUS_FAILED
	job.Error = prevErr.Error()

	_, err := j.JobRepository.Update(job)
	if err != nil {
		return err
	}

	return prevErr
}
