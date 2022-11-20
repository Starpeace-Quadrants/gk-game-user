package main

import (
	"github.com/gookit/event"
	"github.com/ronappleton/gk-game-user/consumer"
	kafka "github.com/ronappleton/gk-kafka"
	"log"
)

func main() {
	event.On("messageReceived", event.ListenerFunc(func(e event.Event) error {
		consumer.ProcessMessage(e)
		return nil
	}), event.Normal)

	go kafka.SaramaConsume("kafka:9092", "auth", "user_in")

	log.Println("User Service Running...")
}
