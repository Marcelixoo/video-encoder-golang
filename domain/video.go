package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type Video struct {
	// unique identifier within our context
	ID string `json:"encoded_video_folder" valid:"uuid" gorm:"type:uuid;primary_key"`
	// external identifier used by collaborating systems
	ResourceID string    `json:"resource_id" valid:"notnull" gorm:"type:varchar(255)"`
	FilePath   string    `json:"file_path" valid:"notnull" gorm:"type:varchar(255)"`
	CreatedAt  time.Time `json:"-" valid:"-"`
	// slice of jobs
	Jobs []*Job `json:"-" valid:"-" gorm:"ForeignKey:VideoID"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func NewVideo() *Video {
	return &Video{}
}

// That's how you create "instance" methods
func (video *Video) Validate() error {
	// ignore first returned value
	_, err := govalidator.ValidateStruct(video)

	if err != nil {
		return err
	}

	return nil // normally returned when no error is found
}
