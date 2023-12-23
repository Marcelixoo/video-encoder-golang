package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type Job struct {
	ID        string    `json:"job_id" valid:"uuid" gorm:"type:uuid;primary_key"`
	Status    string    `json:"status" valid:"notnull"`
	Video     *Video    `json:"video" valid:"-"`
	VideoID   string    `json:"-" valid:"-" gorm:"column:video_id;type:uuid;notnull"`
	Error     string    `json:"error" valid:"-"`
	CreatedAt time.Time `json:"created_at" valid:"-"`
	UpdatedAt time.Time `json:"updated_at" valid:"-"`
}

const (
	JOB_STATUS_PENDING     = "PENDING"
	JOB_STATUS_STARTING    = "STARTING"
	JOB_STATUS_DOWNLOADING = "DOWNLOADING"
	JOB_STATUS_FRAGMENTING = "FRAGMENTING"
	JOB_STATUS_ENCODING    = "ENCODING"
	JOB_STATUS_UPLOADING   = "UPLOADING"
	JOB_STATUS_FINISHING   = "FINISHING"
	JOB_STATUS_COMPLETED   = "COMPLETED"
	JOB_STATUS_FAILED      = "FAILED"
)

func NewJob(video *Video) (*Job, error) {
	if err := video.Validate(); err != nil {
		return nil, err
	}

	job := Job{
		ID:        uuid.NewV4().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		Status: JOB_STATUS_PENDING,
		Video:  video,
	}

	if err := job.Validate(); err != nil {
		return nil, err
	}

	return &job, nil
}

func (job *Job) Validate() error {
	var err error

	_, err = govalidator.ValidateStruct(job)
	if err != nil {
		return err
	}

	return nil
}
