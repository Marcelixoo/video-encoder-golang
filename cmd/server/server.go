package main

import (
	"encoder/application/repositories"
	"encoder/application/services"
	"encoder/framework/database"
	"encoder/framework/gcp"
	"encoder/framework/queue"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var db database.Database

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("could not load .env file")
	}

	autoMigrateDb, err := strconv.ParseBool(os.Getenv("DB_AUTO_MIGRATE_ON"))
	if err != nil {
		log.Fatalf("could not parse DB_AUTO_MIGRATE_ON value")
	}

	debug, err := strconv.ParseBool(os.Getenv("DB_DEBUG_MODE_ON"))
	if err != nil {
		log.Fatalf("could not parse DB_DEBUG_MODE_ON value")
	}

	db.AutoMigrateDb = autoMigrateDb
	db.Debug = debug
	db.DsnTest = os.Getenv("DSN_TEST") // DSN stands for Data Source Name
	db.Dsn = os.Getenv("DSN")
	db.DbTypeTest = os.Getenv("DB_TYPE_TEST")
	db.DbType = os.Getenv("DB_TYPE")
	db.Env = os.Getenv("ENVIRONMENT")
}

func main() {
	messageChannel := make(chan amqp.Delivery)
	jobReturnChannel := make(chan services.JobWorkerResult)

	dbConnection, err := db.Connect()
	if err != nil {
		log.Fatalf("could not connect to database %v", err)
	}
	defer db.Connection.Close()

	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()
	defer ch.Close()
	rabbitMQ.Consume(messageChannel)

	bucketName := os.Getenv("OUTPUT_BUCKET_NAME")
	videoStorage, err := gcp.NewCloudStorageReader(bucketName)
	if err != nil {
		log.Fatalf("could not connect to cloud bucket %s %v", bucketName, err)
	}

	jobService := services.NewJobService(
		repositories.NewJobRepository(dbConnection),
		services.NewVideoService(
			repositories.NewVideoRepository(dbConnection),
			videoStorage,
		),
	)

	jobManager := services.NewJobManager(
		dbConnection,
		rabbitMQ,
		jobReturnChannel,
		messageChannel,
		jobService,
	)

	jobManager.Start(ch)
}
