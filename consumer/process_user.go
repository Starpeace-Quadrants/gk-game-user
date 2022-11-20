package consumer

import (
	"github.com/kamva/mgm/v3"
	"github.com/ronappleton/gk-game-user/storage/mongo"
	transport "github.com/ronappleton/gk-message-transport"
	"go.mongodb.org/mongo-driver/bson"
)

func processGetUserProfile(key string, message transport.ServiceMessage) {
	profile := &mongo.UserProfile{}
	_ = mgm.Coll(profile).First(bson.M{"user_id": message.UserId}, profile)

	if len(profile.Alias) > 0 {
		var companies []mongo.Company
		_ = mgm.Coll(&mongo.Company{}).SimpleFind(&companies, bson.M{"_id": message.UserId})

		profile.Companies = companies
	}

	var results map[string]interface{}
	results["profile"] = profile

	reply(key, message, "get_user_profile", results)
}

func processSetUserImage(key string, message transport.ServiceMessage) {
	profile := &mongo.UserProfile{}
	_ = mgm.Coll(profile).FindByID(message.UserId, profile)

	profile.ImagePath = message.ArgumentStore.GetString("image_path")
	options := mgm.UpsertTrueOption()
	err := mgm.Coll(profile).Update(profile, options)

	var results map[string]interface{}
	results["result"] = err != nil

	reply(key, message, "set_user_image", results)

}

func processCreateUserAlias(key string, message transport.ServiceMessage) {
	var results map[string]interface{}

	profile := &mongo.UserProfile{}
	err := mgm.Coll(profile).FindByID(message.UserId, profile)

	if err != nil {
		results["result"] = false
		results["message"] = "user profile does not exist"
		reply(key, message, "create_user_alias", results)

		return
	}

	if len(profile.Alias) > 0 {
		results["result"] = false
		results["message"] = "player aliases can not be change"
		reply(key, message, "create_user_alias", results)

		return
	}

	results["result"] = true

	reply(key, message, "create_user_alias", results)
}
