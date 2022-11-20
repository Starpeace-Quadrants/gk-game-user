package consumer

import (
	"github.com/gookit/event"
	kafka "github.com/ronappleton/gk-kafka"
	"github.com/ronappleton/gk-kafka/storage"
	"github.com/ronappleton/gk-message-transport"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func ProcessMessage(event event.Event, db *mongo.Client) {
	data := storage.New()
	data.Populate(event.Data())

	message := data.GetMessage()

	serviceMessage := transport.BytesToServiceMessage(message.Value)

	switch command := serviceMessage.Command; command {
	case "get_user_profile":
		processGetUserProfile(string(message.Key), serviceMessage)
	case "create_user_company":
		processCreateUserCompany(string(message.Key), serviceMessage)
	case "list_user_companies":
		processListUserCompanies(string(message.Key), serviceMessage)
	case "remove_user_company":
		processRemoveUserCompany(string(message.Key), serviceMessage)
	case "list_companies":
		processListCompanies(string(message.Key), serviceMessage)
	case "set_user_image":
		processSetUserImage(string(message.Key), serviceMessage)
	case "create_user_alias":
		processCreateUserAlias(string(message.Key), serviceMessage)
	case "get_company_balance":
		processGetCompanyBalance(string(message.Key), serviceMessage)
	case "update_company_balance":
		processUpdateCompanyBalance(string(message.Key), serviceMessage)
	}

}

func reply(key string, message transport.ServiceMessage, command string, results map[string]interface{}) {
	reply := transport.NewClientMessage()
	reply.Command = command
	reply.Topic = message.Topic
	reply.Results = results

	replyBytes := reply.ToBytes()

	topic, _ := kafka.GetTopicByName("user", "out")

	kafka.Produce([]byte(key), replyBytes, topic, time.Now().Add(10*time.Second))
}
