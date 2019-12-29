package edumodel

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var authCollection *mongo.Collection

func checkAuthCollection() {
	if authCollection == nil {
		authCollection = GetCollection("authorize")
	}
}

type UserAuth struct {
	UID       string `bson:"uid"`
	Pwd       string `bson:"pwd"`
	QuestionA string `bson:"qa"`
	AnswerA   string `bson:"aa"`
	QuestionB string `bson:"qb"`
	AnswerB   string `bson:"ab"`
	QuestionC string `bson:"qc"`
	AnswerC   string `bson:"ac"`
}

func AddUserAuth(newUserAuth *UserAuth) bool {
	checkAuthCollection()

	if newUserAuth == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := authCollection.InsertOne(ctx, newUserAuth)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return false
	}

	return true
}

func GetUserAuthByUID(uid string) *UserAuth {
	checkAuthCollection()

	if uid == "" {
		fmt.Println("[MODEL]", "GetUserAuth: uid cannot be empty")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"uid": uid}

	var result UserAuth

	err := authCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return nil
	}

	return &result
}

func UpdateUserAuthByUID(uid, newPwd, Qa, Aa, Qb, Ab, Qc, Ac string) bool {
	checkAuthCollection()

	if uid == "" {
		fmt.Println("[MODEL]", "GetUserAuth: uid cannot be empty")
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"uid": uid}

	var result UserAuth

	err := authCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return false
	}

	if newPwd == "" {
		newPwd = result.Pwd
	}
	if Qa == "" {
		Qa = result.QuestionA
	}
	if Qb == "" {
		Qb = result.QuestionB
	}
	if Qc == "" {
		Qc = result.QuestionC
	}
	if Aa == "" {
		Aa = result.AnswerA
	}
	if Ab == "" {
		Ab = result.AnswerB
	}
	if Ac == "" {
		Ac = result.AnswerC
	}

	update := bson.D{
		{"$set", bson.M{
			"pwd": newPwd,
			"qa":  Qa,
			"qb":  Qb,
			"qc":  Qc,
			"aa":  Aa,
			"ab":  Ab,
			"ac":  Ac,
		}},
	}

	_, err = authCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return false
	}

	return true

}

func DeleteUserAuthByUID(uid string) bool {
	checkAuthCollection()

	if uid == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"uid": uid}

	_, err := authCollection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return false
	}

	return true
}
