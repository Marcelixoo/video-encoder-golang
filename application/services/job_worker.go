package services

import (
	"encoder/domain"
	"encoding/json"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

type JobWorkerResult struct {
	Job     domain.Job
	Message *amqp.Delivery
	Error   error
}

func JobWorker(
	messageChannel chan amqp.Delivery,
	returnChan chan JobWorkerResult,
	jobService JobService,
	job domain.Job,
	workerID int,
) {
	/*
		{
			"resource_id": "d89f161c-a05b-11ee-8c90-0242ac120002",
			"file_path": "convite.mp4"

		}
	*/

	for message := range messageChannel {
		var (
			err   error
			video *domain.Video
		)

		err = json.Unmarshal(message.Body, video)
		if err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		// goddam do that at construction time!
		job.Video = video
		job.OutputBucketPath = os.Getenv("OUTPUT_BUCKET_NAME")
		job.ID = uuid.NewV4().String()
		job.Status = "STARTING"
		job.CreatedAt = time.Now()

		if _, err := jobService.JobRepository.Insert(&job); err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		jobService.Job = &job
		if err := jobService.Start(video); err != nil {
			returnChan <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		returnChan <- returnJobResult(job, message, nil)
	}
}

func returnJobResult(job domain.Job, message amqp.Delivery, err error) JobWorkerResult {
	return JobWorkerResult{
		Job:     job,
		Message: &message,
		Error:   err,
	}
}
