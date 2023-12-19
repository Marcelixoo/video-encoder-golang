package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type Job struct {
	ID               string    `json:"job_id" valid:"uuid" gorm:"type:uuid;primary_key"`
	OutputBucketPath string    `json:"output_bucket_path" valid:"notnull"`
	Status           string    `json:"status" valid:"notnull"`
	Video            *Video    `json:"video" valid:"-"`
	VideoID          string    `json:"-" valid:"-" gorm:"column:video_id;type:uuid;notnull"`
	Error            string    `json:"error" valid:"-"`
	CreatedAt        time.Time `json:"created_at" valid:"-"`
	UpdatedAt        time.Time `json:"updated_at" valid:"-"`
}

func NewJob(output string, status string, video *Video) (*Job, error) {
	job := Job{
		ID:        uuid.NewV4().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		OutputBucketPath: output,
		Status:           status,
		Video:            video,
	}

	err := job.Validate()

	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (job *Job) Validate() error {
	_, err := govalidator.ValidateStruct(job)

	if err != nil {
		return err
	}

	return nil
}
