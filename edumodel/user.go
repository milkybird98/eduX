package edumodel

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	Com1A         string `bson:"com1a" json:"com1a"`
	Com1B         string `bson:"com1b" json:"com1b"`
	Com2A         string `bson:"com2a" json:"com2a"`
	Com2B         string `bson:"com2b" json:"com2b"`
	Com3A         string `bson:"com3a" json:"com3a"`
	Com3B         string `bson:"com3b" json:"com3b"`
	Com4A         string `bson:"com4a" json:"com4a"`
	Com4B         string `bson:"com4b" json:"com4b"`
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

func GetUserSimpleAll() *[]*User {
	checkClassCollection()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	option := options.Find().SetProjection(bson.M{"name": 1, "uid": 1})

	var result []*User
	cur, err := userCollection.Find(ctx, filter, option)
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
	if newUserData.Com1A != originData.Com1A {
		originData.Com1A = newUserData.Com1A
	}
	if newUserData.Com1B != originData.Com1B {
		originData.Com1B = newUserData.Com1B
	}
	if newUserData.Com2A != originData.Com2A {
		originData.Com2A = newUserData.Com2A
	}
	if newUserData.Com2B != originData.Com2B {
		originData.Com2B = newUserData.Com2B
	}
	if newUserData.Com3A != originData.Com3A {
		originData.Com3A = newUserData.Com3A
	}
	if newUserData.Com3B != originData.Com3B {
		originData.Com3B = newUserData.Com3B
	}
	if newUserData.Com4A != originData.Com4A {
		originData.Com4A = newUserData.Com4A
	}
	if newUserData.Com4B != originData.Com4B {
		originData.Com4B = newUserData.Com4B
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"uid", newUserData.UID}}
	update := bson.D{
		{"$set", bson.M{
			"name":          originData.Name,
			"class":         originData.Class,
			"gender":        originData.Gender,
			"birth":         originData.Birth,
			"political":     originData.Political,
			"contact":       originData.Contact,
			"iscontactpub":  originData.IsContactPub,
			"email":         originData.Email,
			"isemailpub":    originData.IsEmailPub,
			"localion":      originData.Localion,
			"islocalionpub": originData.IsLocalionPub,
			"job":           originData.Job,
			"com1a":         originData.Com1A,
			"com1b":         originData.Com1B,
			"com2a":         originData.Com2A,
			"com2b":         originData.Com2B,
			"com3a":         originData.Com3A,
			"com3b":         originData.Com3B,
			"com4a":         originData.Com4A,
			"com4b":         originData.Com4B,
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
