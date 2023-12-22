package services

import (
	"encoder/domain"
	"encoder/framework/queue"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
)

type JobManager struct {
	Db               *gorm.DB
	Domain           domain.Job
	MessageChannel   chan amqp.Delivery
	JobReturnChannel chan JobWorkerResult
	RabbitMQ         *queue.RabbitMQ
	JobService       JobService
}

type JobNotificationError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewJobManager(
	db *gorm.DB,
	rabbitMQ *queue.RabbitMQ,
	jobReturnChannel chan JobWorkerResult,
	messageChannel chan amqp.Delivery,
	jobService JobService,
) *JobManager {
	return &JobManager{
		Db:               db,
		Domain:           domain.Job{}, // wut?
		MessageChannel:   messageChannel,
		JobReturnChannel: jobReturnChannel,
		RabbitMQ:         rabbitMQ,
		JobService:       jobService,
	}
}

func (j *JobManager) Start(ch *amqp.Channel) {
	var concurrency int

	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_WORKERS"))
	if err != nil {
		log.Default().Println("could not load CONCURRENCY_WORKERS; using default value of 1")
		concurrency = 1
	}

	for qtdProcesses := 0; qtdProcesses < concurrency; qtdProcesses++ {
		go JobWorker(
			j.MessageChannel,
			j.JobReturnChannel,
			j.JobService,
			j.Domain,
			qtdProcesses,
		)
	}

	for jobResult := range j.JobReturnChannel {
		if jobResult.Error != nil {
			j.checkParseErrors(jobResult)
		}

		var requeue = false

		if err != nil {
			jobResult.Message.Reject(requeue)
		}

		if err := j.notifySuccess(jobResult, ch); err != nil {
			jobResult.Message.Reject(requeue)
		}
	}
}

func (j *JobManager) notifySuccess(jobResult JobWorkerResult, ch *amqp.Channel) error {
	jobJson, err := json.Marshal(jobResult.Job)
	if err != nil {
		return err
	}

	if err := j.notify(jobJson); err != nil {
		return err
	}

	multiple := false
	return jobResult.Message.Ack(multiple)
}

// this method is unnecessarily complex
func (j *JobManager) checkParseErrors(jobResult JobWorkerResult) error {
	if jobResult.Job.ID != "" {
		log.Printf("")
	} else {
		log.Printf("")
	}

	errorMsg := JobNotificationError{
		Message: string(jobResult.Message.Body),
		Error:   jobResult.Error.Error(),
	}
	jobJson, err := json.Marshal(errorMsg)
	if err != nil {
		return err
	}

	if err := j.notify(jobJson); err != nil {
		return err
	}

	requeue := false
	return jobResult.Message.Reject(requeue)
}

func (j *JobManager) notify(jobJson []byte) error {
	return j.RabbitMQ.Notify(
		string(jobJson),
		"application/json",
		os.Getenv("RABBITMQ_NOTIFICATION_EX"),
		os.Getenv("RABBITMQ_NOTIFICATION_ROUTING_KEY"),
	)
}
