package database

import (
	"encoder/domain"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
)

type Database struct {
	Connection    *gorm.DB
	Dsn           string
	DsnTest       string
	DbType        string
	DbTypeTest    string
	Debug         bool
	AutoMigrateDb bool
	Env           string
}

func NewDb() *Database {
	return &Database{
		Debug: false,
		Env:   "production",
	}
}

func NewDbTest() *Database {
	return &Database{
		Debug:         true,
		AutoMigrateDb: true,
		Env:           "test",
		DbType:        "sqlite3",
		Dsn:           ":memory:",
	}
}

func (d *Database) Connect() (*gorm.DB, error) {
	var err error

	if d.Connection != nil {
		return d.Connection, nil
	}

	d.Connection, err = gorm.Open(d.DbType, d.Dsn)
	if err != nil {
		return nil, err
	}

	d.Connection.LogMode(d.Debug)

	if d.AutoMigrateDb {
		d.Connection.AutoMigrate(&domain.Video{}, &domain.Job{})
		d.Connection.Model(domain.Job{}).AddForeignKey("video_id", "videos (id)", "CASCADE", "CASCADE")
	}

	return d.Connection, nil

}
