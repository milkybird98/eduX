package edumodel

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"context"
	"fmt"
	"time"
)

var userCollection *mongo.Collection

type User struct{
	_id					primitive.ObjectID 	`bson:"_id"`
	Name 				string
	UID  				string
	Pwd					string
	Plcae 			string
	Class 			string
	Gender			int
}

func checkUserCollection()  {
	if userCollection == nil{
		userCollection = GetCollection("user")
	}
}

func AddUser(newUser *User) bool {
	checkUserCollection()

	if newUser == nil{
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_,err := userCollection.InsertOne(ctx,newUser)
	if err != nil {
		fmt.Println("Add new user into database fail, error: ",err)
		return false
	}

	return true
}

func GetUserByUID(uid string) *User {
	checkUserCollection()

	if uid == ""{
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"uid":uid}

	var result User
	err := userCollection.FindOne(ctx,filter).Decode(&result)
	if err != nil{
		fmt.Println(err)
		return nil
	}

	return &result
}

func GetUserByClass(className string) *[]*User {
	checkUserCollection()

	if className == ""{
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"class":className}

	var result []*User
	cur,err := userCollection.Find(ctx,filter)
	if err != nil{
		return nil
	}
	defer cur.Close(ctx)

	for cur.Next(ctx){
		var user User
		if err := cur.Decode(&user); err != nil {
			return nil
		}
		result = append(result,&user)
	}

	return &result
}

func UpdateUserByID(uid string,class string, name string, pwd string, gender int) bool {
	checkUserCollection()

	if uid == ""{
		return false
	}
	
	originData := GetUserByUID(uid)

	if class == "" {
		class = originData.Class
	}
	if name == "" {
		name = originData.Name
	}
	if pwd == "" {
		pwd = originData.Pwd
	}
	if (gender != 1 && gender != 2) {
		gender = originData.Gender
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"uid",uid}}
	update := bson.D{
    {"$set", bson.M{
				"name": name,
				"pwd": pwd,
				"class": class,
				"gender": gender,
    }},
	}

	_,err := userCollection.UpdateOne(ctx,filter,update)
	if err != nil{
		fmt.Println(err)
		return false
	}

	return true
}

func DeleteUserByUID(uid string) bool {
	checkUserCollection()

	if uid == ""{
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"uid":uid}

	_,err := userCollection.DeleteOne(ctx,filter)
	if err!=nil {
		fmt.Println(err)
		return false
	}

	return true
}
