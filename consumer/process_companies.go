package consumer

import (
	"context"
	"github.com/kamva/mgm/v3"
	"github.com/ronappleton/gk-game-user/storage/mongo"
	transport "github.com/ronappleton/gk-message-transport"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

func processCreateUserCompany(key string, message transport.ServiceMessage) {
	var results map[string]interface{}
	company := mongo.Company{}
	_ = mgm.Coll(&company).First(bson.M{"name": message.ArgumentStore.GetString("name")}, &company)

	if len(company.Name) > 0 {
		results["result"] = false
		results["message"] = "company already exists"
		reply(key, message, "create_user_company", results)

		return
	}

	id, _ := primitive.ObjectIDFromHex(message.UserId)
	comp := mongo.NewCompany(id, message.ArgumentStore.GetString("name"))

	_ = mgm.Coll(comp).Create(comp)

	results["result"] = true

	reply(key, message, "create_user_company", results)
}

func processListUserCompanies(key string, message transport.ServiceMessage) {
	var results map[string]interface{}
	var companies []mongo.Company

	userId, _ := primitive.ObjectIDFromHex(message.UserId)
	err := mgm.Coll(&mongo.Company{}).SimpleFind(companies, bson.M{"user_id": userId})
	if err != nil {
		log.Printf("unable to list user companies, userId: %s", message.UserId)

		return
	}

	results["companies"] = companies

	reply(key, message, "list_user_companies", results)
}

func processRemoveUserCompany(key string, message transport.ServiceMessage) {
	company := mongo.Company{}

	companyId, _ := primitive.ObjectIDFromHex(message.ArgumentStore.GetString("company_id"))
	userId, _ := primitive.ObjectIDFromHex(message.UserId)

	_, err := mgm.Coll(&company).DeleteOne(context.TODO(), bson.M{"_id": companyId, "user_id": userId})
	if err != nil {
		log.Printf("failed to delete company, userId: %s companyId: %s", message.UserId, message.ArgumentStore.GetString("company_id"))
	}
}

func processListCompanies(key string, message transport.ServiceMessage) {
	var companies []mongo.Company

	_ = mgm.Coll(&mongo.Company{}).SimpleFind(companies, bson.M{"first_char": message.ArgumentStore.GetString("first_char")})

	chunks := mongo.ChunkCompanies(companies, 100)

	for _, chunk := range chunks {
		var results map[string]interface{}
		results["companies"] = chunk
		results["first_char"] = message.ArgumentStore.GetString("first_char")

		reply(key, message, "list_companies", results)
	}
}

func processGetCompanyBalance(key string, message transport.ServiceMessage) {
	var results map[string]interface{}
	company := mongo.Company{}
	err := mgm.Coll(&company).First(bson.M{"_id": message.ArgumentStore.GetString("companyId")}, &company)
	if err != nil {
		results["balance"] = false
		results["message"] = "unable to obtain company balance"

		reply(key, message, "get_company_balance", results)
	}

	results["balance"] = company.Balance

	reply(key, message, "get_company_balance", results)
}

func processUpdateCompanyBalance(key string, message transport.ServiceMessage) {
	company := mongo.Company{}
	err := mgm.Coll(&company).First(bson.M{"_id": message.ArgumentStore.GetString("companyId")}, &company)

	if err == nil {
		log.Printf("unable to update company: %v balance")

		return
	}

	company.Balance += message.ArgumentStore.GetInt64("amount")

	err = mgm.Coll(&company).Update(&company)

	if err != nil {
		log.Printf("unable to update company: %v balance")
	}
}
