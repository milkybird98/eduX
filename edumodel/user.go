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
	Name          string `bson:"name" json:"name"`
	UID           string `bson:"uid" json:"uid"`
	Place         string `bson:"place" json:"place"`
	Class         string `bson:"class" json:"class"`
	Gender        int    `bson:"gender" json:"gender"`
	Birth         string `bson:"birth" json:"birth"`
	Political     int    `bson:"political" json:"political"`
	Contact       string `bson:"contact" json:"contact"`
	IsContactPub  bool   `bson:"iscontactpub" json:"iscontactpub"`
	Email         string `bson:"email" json:"email"`
	IsEmailPub    bool   `bson:"isemailpub" json:"isemailpub"`
	Localion      string `bson:"localion" json:"localion"`
	IsLocalionPub bool   `bson:"islocalionpub" json:"islocalionpub"`
	Job           string `bson:"job" json:"job"`
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
		fmt.Println("[MODEL] Add new user into database fail, error: ", err)
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
		fmt.Println("[MODEL]", err)
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
		fmt.Println("[MODEL]", err)
		return nil
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var user User
		if err := cur.Decode(&user); err != nil {
			fmt.Println("[MODEL]", err)
			return nil
		}
		result = append(result, &user)
	}

	return &result
}

func GetUserNumber() int {
	checkClassCollection()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}

	count, err := userCollection.CountDocuments(ctx, filter)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return -1
	}

	return int(count)
}

func UpdateUserByID(newUserData *User) bool {
	checkUserCollection()

	if newUserData == nil {
		return false
	}

	originData := GetUserByUID(newUserData.UID)

	if newUserData.Name != "" {
		originData.Name = newUserData.Name
	}
	if newUserData.Gender != 0 {
		originData.Gender = newUserData.Gender
	}
	if newUserData.Birth != "" {
		originData.Birth = newUserData.Birth
	}
	if newUserData.Political != 0 {
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
	if newUserData.Localion != "" {
		originData.Localion = newUserData.Localion
	}
	if newUserData.IsLocalionPub != originData.IsLocalionPub {
		originData.IsLocalionPub = newUserData.IsLocalionPub
	}
	if newUserData.Job != originData.Job {
		originData.Job = newUserData.Job
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"uid", newUserData.UID}}
	update := bson.D{
		{"$set", bson.M{
			"name":          originData.Name,
			"class":         originData.Class,
			"gender":        originData.Gender,
			"bitrh":         originData.Birth,
			"political":     originData.Political,
			"contact":       originData.Contact,
			"iscontactpub":  originData.IsContactPub,
			"email":         originData.Email,
			"isemailpub":    originData.IsEmailPub,
			"localion":      originData.Localion,
			"islocalionpub": originData.IsLocalionPub,
			"job":           originData.Job,
		}},
	}

	_, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("[MODEL]", err)
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

	filter := bson.M{"uid": bson.M{"$in": uidList}, "class": ""}
	update := bson.D{
		{"$set", bson.M{
			"class": ClassName,
		}},
	}

	_, err := userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		fmt.Println("[MODEL]", err)
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

	filter := bson.M{"uid": bson.M{"$in": uidList}, "class": ClassName}
	update := bson.D{
		{"$set", bson.M{
			"class": "",
		}},
	}

	_, err := userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		fmt.Println("[MODEL]", err)
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
		fmt.Println("[MODEL]", err)
		return false
	}

	return true
}
