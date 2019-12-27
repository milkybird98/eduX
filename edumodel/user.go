package edumodel

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

type User struct {
	Name          string
	UID           string
	Pwd           string
	Place         string
	Class         string
	Gender        int
	Birth         string
	Political     string
	Contact       string
	IsContactPub  bool
	Email         string
	IsEmailPub    bool
	Location      string
	IsLocationPub bool
}

func checkUserCollection() {
	if userCollection == nil {
		userCollection = GetCollection("user")
	}
}

func AddUser(newUser *User) bool {
	checkUserCollection()

	if newUser == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		fmt.Println("Add new user into database fail, error: ", err)
		return false
	}

	return true
}

func GetUserByUID(uid string) *User {
	checkUserCollection()

	if uid == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"uid": uid}

	var result User
	err := userCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &result
}

func GetUserByClass(className string) *[]*User {
	checkUserCollection()

	if className == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"class": className}

	var result []*User
	cur, err := userCollection.Find(ctx, filter)
	if err != nil {
		return nil
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var user User
		if err := cur.Decode(&user); err != nil {
			return nil
		}
		result = append(result, &user)
	}

	return &result
}

func UpdateUserByID(newUserData *User) bool {
	checkUserCollection()

	if newUserData == nil {
		return false
	}

	originData := GetUserByUID(newUserData.UID)

	if newUserData.Pwd != "" {
		originData.Pwd = newUserData.Pwd
	}
	if newUserData.Name != "" {
		originData.Name = newUserData.Name
	}
	if newUserData.Gender != 0 {
		originData.Gender = newUserData.Gender
	}
	if newUserData.Birth != "" {
		originData.Birth = newUserData.Birth
	}
	if newUserData.Political != "" {
		originData.Political = newUserData.Political
	}
	if newUserData.Contact != "" {
		originData.Contact = newUserData.Contact
	}
	if newUserData.IsContactPub != originData.IsContactPub {
		originData.IsContactPub = newUserData.IsContactPub
	}
	if newUserData.Email != "" {
		originData.Email = newUserData.Email
	}
	if newUserData.IsEmailPub != originData.IsEmailPub {
		originData.IsEmailPub = newUserData.IsEmailPub
	}
	if newUserData.Location != "" {
		originData.Location = newUserData.Location
	}
	if newUserData.IsLocationPub != originData.IsLocationPub {
		originData.IsLocationPub = newUserData.IsLocationPub
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"uid", newUserData.UID}}
	update := bson.D{
		{"$set", bson.M{
			"name":          originData.Name,
			"pwd":           originData.Pwd,
			"class":         originData.Class,
			"gender":        originData.Gender,
			"bitrh":         originData.Birth,
			"political":     originData.Political,
			"contact":       originData.Contact,
			"iscontactpub":  originData.IsContactPub,
			"email":         originData.Email,
			"isemailpub":    originData.IsEmailPub,
			"location":      originData.Location,
			"islocationpub": originData.IsLocationPub,
		}},
	}

	_, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func AddUserToClassByUID(uidList []string, ClassName string) bool {
	checkUserCollection()

	if uidList == nil || len(uidList) <= 0 {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"uid": uidList, "class": ""}
	update := bson.D{
		{"$set", bson.M{
			"class": ClassName,
		}},
	}

	_, err := userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func DeleteUserFromClassByUID(uidList []string, ClassName string) bool {
	checkUserCollection()

	if uidList == nil || len(uidList) <= 0 {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"uid": uidList, "class": ClassName}
	update := bson.D{
		{"$set", bson.M{
			"class": "",
		}},
	}

	_, err := userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func DeleteUserByUID(uid string) bool {
	checkUserCollection()

	if uid == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"uid": uid}

	_, err := userCollection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
