package edumodel

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"fmt"
	"time"
)

var collection *mongo.Collection

type User struct{
	Name 	string
	UID  	string
	Pwd		string
	Plcae string
	Class string
	Gender		int
}

func checkCollection()  {
	if collection == nil{
		collection = GetCollection("user")
	}
}

func AddUser(newUser *User) bool {
	checkCollection()

	if newUser == nil{
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_,err := collection.InsertOne(ctx,newUser)
	if err != nil {
		fmt.Println("Add new user into database fail, error: ",err)
		return false
	}

	return true
}

func GetUserByUID(uid string) *User {
	checkCollection()

	if uid == ""{
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"uid",uid}}

	var result User
	err := collection.FindOne(ctx,filter).Decode(&result)
	if err != nil{
		fmt.Println(err)
		return nil
	}

	return &result
}

func GetUserByClass(className string) *[]*User {
	checkCollection()

	if className == ""{
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"class",className}}

	var result []*User
	cur,err := collection.Find(ctx,filter)
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
	checkCollection()

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
    {"$set", bson.D{
				{"name", name},
				{"pwd", pwd},
				{"class",class},
				{"gender", gender},
    }},
	}

	_,err := collection.UpdateOne(ctx,filter,update)
	if err != nil{
		fmt.Println(err)
		return false
	}

	return true
}

func DeleteUserByUID(uid string) bool {
	checkCollection()

	if uid == ""{
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"uid",uid}}

	_,err := collection.DeleteOne(ctx,filter)
	if err!=nil {
		fmt.Println(err)
		return false
	}

	return true
}