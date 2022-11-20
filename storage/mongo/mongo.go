package mongo

import (
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strings"
)

func init() {
	host := os.Getenv("MONGO_SERVER_HOST")
	port := os.Getenv("MONGO_SERVER_PORT")
	database := os.Getenv("MONGO_SERVER_DATABASE")

	if err := mgm.SetDefaultConfig(nil, database, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port))); err != nil {
		panic(err)
	}
}

type UserProfile struct {
	mgm.DefaultModel
	UserId    primitive.ObjectID
	Alias     string    `bson:"alias" json:"alias"`
	ImagePath string    `bson:"image_path" json:"image_path"`
	Companies []Company `bson:"-" json:"companies"`
}

func NewUserProfile(userId primitive.ObjectID, alias string) *UserProfile {
	return &UserProfile{
		UserId: userId,
		Alias:  alias,
	}
}

type Company struct {
	mgm.DefaultModel
	UserId    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name      string             `bson:"name" json:"name"`
	FirstChar string             `bson:"first_char" json:"first_char"`
	Balance   int64              `bson:"balance" json:"balance"`
}

func NewCompany(userId primitive.ObjectID, name string) *Company {
	return &Company{
		UserId:    userId,
		Name:      name,
		FirstChar: strings.ToLower(name[0:1]),
	}
}

func ChunkCompanies(items []Company, chunkSize int) (chunks [][]Company) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}

	return append(chunks, items)
}
