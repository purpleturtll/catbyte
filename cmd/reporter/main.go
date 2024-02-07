package main

import (
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type Message struct {
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Timestamp string `json:"timestamp"`
	Content   string `json:"message"`
}

func main() {
	r := gin.Default()

	r.GET("/message/list", func(c *gin.Context) {
		sender := c.Query("sender")
		receiver := c.Query("receiver")

		redisClient := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})

		val, err := redisClient.HGetAll(sender + "-" + receiver).Result()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		var messages []Message
		for timestamp, content := range val {
			messages = append(messages, Message{Sender: sender, Receiver: receiver, Timestamp: timestamp, Content: content})
		}

		// Sort messages in reverse chronological order
		sort.SliceStable(messages, func(i, j int) bool {
			ti, _ := time.Parse(time.RFC3339, messages[i].Timestamp)
			tj, _ := time.Parse(time.RFC3339, messages[j].Timestamp)
			return ti.After(tj)
		})

		c.JSON(200, messages)
	})

	r.Run("localhost:8081")
}
