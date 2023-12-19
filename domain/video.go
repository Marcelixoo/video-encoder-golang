package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type Video struct {
	ID         string    `json:"encoded_video_folder" valid:"uuid" gorm:"type:uuid;primary_key"`
	ResourceID string    `json:"resource_id" valid:"notnull" gorm:"type:varchar(255)"`
	FilePath   string    `json:"file_path" valid:"notnull" gorm:"type:varchar(255)"`
	CreatedAt  time.Time `json:"-" valid:"-"`
	Jobs       []*Job    `json:"-" valid:"-" gorm:"ForeignKey:VideoID"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func NewVideo(id string, resourceID string, filepath string) *Video {
	return &Video{
		ID:        id,
		CreatedAt: time.Now(),

		FilePath:   filepath,
		ResourceID: resourceID,
	}
}

func (video *Video) Validate() error {
	_, err := govalidator.ValidateStruct(video)

	if err != nil {
		return err
	}

	return nil
}
