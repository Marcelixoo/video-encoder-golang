package repositories

import (
	"encoder/domain"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type JobRepository interface {
	Insert(job *domain.Job) (*domain.Job, error)
	Find(id string) (*domain.Job, error)
	Update(job *domain.Job) (*domain.Job, error)
}

type JobRepositoryDb struct {
	Db *gorm.DB
}

func (repo JobRepositoryDb) Insert(job *domain.Job) (*domain.Job, error) {
	if job.ID == "" {
		job.ID = uuid.NewV4().String()
	}

	err := repo.Db.Create(job).Error
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (repo JobRepositoryDb) Find(id string) (*domain.Job, error) {
	var job domain.Job

	repo.Db.Preload("Video").First(&job, "id = ?", id)
	if job.ID == "" {
		return nil, fmt.Errorf("Could not find job with ID %q", id)
	}

	return &job, nil
}

func (repo JobRepositoryDb) Update(job *domain.Job) (*domain.Job, error) {
	err := repo.Db.Save(&job).Error
	if err != nil {
		return nil, err
	}

	return job, err
}
