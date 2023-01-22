package database

import (
	"encoder/domain"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Db          *gorm.DB
	Dsn         string
	DsnTest     string
	DbType      string
	DbTypeTest  string
	Debug       bool
	AutoMigrate bool
	Env         string
}

func NewDb() *Database {
	return &Database{}
}

func NewDbTest() *gorm.DB {
	dbInstance := NewDb()

	dbInstance.Env = "test"
	dbInstance.DbTypeTest = "sqlite3"
	dbInstance.DsnTest = ":memory:"
	dbInstance.AutoMigrate = true
	dbInstance.Debug = true

	conn, err := dbInstance.Connect()
	if err != nil {
		log.Fatalf("Test DB error: %v", err)
	}

	return conn
}

func (d *Database) Connect() (*gorm.DB, error) {
	var err error

	d.Db, err = d.openEnvironmentAwareConnection()

	if err != nil {
		return nil, err
	}

	if d.Debug {
		d.Db.Logger.LogMode(logger.Info)
	}

	if d.AutoMigrate {
		d.Db.AutoMigrate(&domain.Video{}, &domain.Job{})
	}

	return d.Db, nil
}

func (d *Database) openEnvironmentAwareConnection() (*gorm.DB, error) {
	if d.Env == "test" {
		return gorm.Open(sqlite.Open(d.DsnTest))
	}
	return gorm.Open(postgres.Open(d.Dsn))
}
