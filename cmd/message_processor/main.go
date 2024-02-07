package main

import (
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://user:password@localhost:7001")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"message_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Default().Printf("Received a message: %s", d.Body)
			msgParts := strings.SplitN(string(d.Body), ",", 3)
			if len(msgParts) == 3 {
				log.Default().Printf("Storing message between %s and %s in Redis: %s", msgParts[0], msgParts[1], msgParts[2])
				err := redisClient.HSet(msgParts[0]+"-"+msgParts[1], time.Now().Format(time.RFC3339), msgParts[2]).Err()
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
