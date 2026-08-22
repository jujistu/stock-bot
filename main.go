package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ong-gtp/go-stockbot/broker"
	"github.com/ong-gtp/go-stockbot/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

var br broker.Broker

func main() {
	log.Println("Stock bot service starting ...")

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using environment variables")
	}

	rmqHost := os.Getenv("RMQ_HOST")
	rmqUserName := os.Getenv("RMQ_USERNAME")
	rmqPassword := os.Getenv("RMQ_PASSWORD")
	rmqPort := os.Getenv("RMQ_PORT")
	dsn := "amqp://" + rmqUserName + ":" + rmqPassword + "@" + rmqHost + ":" + rmqPort + "/"

	conn, err := amqp.Dial(dsn)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	br.SetUp(ch)
	go br.ReadMessages()
	select {}
}
