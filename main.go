package main

import (
	"github.com/gookit/event"
	"github.com/ronappleton/gk-game-user/consumer"
	"github.com/ronappleton/gk-game-user/storage/mongo"
	kafka "github.com/ronappleton/gk-kafka"
	"log"
)

func main() {
	db, err := mongo.NewDatabase()
	if err != nil {
		panic(err.Error())
	}

	db.Start()

	event.On("messageReceived", event.ListenerFunc(func(e event.Event) error {
		consumer.ProcessMessage(e, db.Client)
		return nil
	}), event.Normal)

	go kafka.SaramaConsume("kafka:9092", "auth", "user_in")

	log.Println("User Service Running...")
}
