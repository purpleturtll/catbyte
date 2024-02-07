package main

import (
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type Message struct {
	Sender   string `json:"sender" binding:"required"`
	Receiver string `json:"receiver" binding:"required"`
	Content  string `json:"message" binding:"required"`
}

func main() {
	r := gin.Default()

	r.POST("/message", func(c *gin.Context) {
		var msg Message
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		conn, err := amqp.Dial("amqp://user:password@localhost:7001/")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
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
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		body := msg.Sender + "," + msg.Receiver + "," + msg.Content
		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"status": "OK"})
	})

	r.Run("localhost:8080")
}
