package edumodel

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"context"
	"fmt"
	"time"
)

var quesCollection *mongo.Collection

type Question struct{
	_id					primitive.ObjectID 	`bson:"_id"`
	Title				string
	Text				string
	SenderUID		string
	QueserUID		string
	ClassName		string
	SendTime		time.Time
	IsSolved		bool
	Answer			string
}

func checkQuesCollection()  {
	if quesCollection == nil{
		quesCollection = GetCollection("question")
	}
}

func AddQuestion(newQuestion *Question) (bool) {
	checkQuesCollection()

	if newQuestion == nil{
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_,err := quesCollection.InsertOne(ctx,newQuestion)
	if err != nil {
		fmt.Println("Add new Question into database fail, error: ",err)
		return false
	}

	return true
}

func GetQuestionByTimeOrder(skip int, limit int, isSolved bool) *[]*News {
	checkQuesCollection()

	if skip <0 || limit <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"issolved":isSolved}
	option := options.Find().SetSort(bson.M{"sendtime":1}).SetSkip(int64(skip)).SetLimit(int64(limit))
	
	var result []*News
	cur,err := quesCollection.Find(ctx,filter,option)
	if err!=nil{
		return nil
	}

	for cur.Next(ctx){
		var news News
		if err := cur.Decode(&news); err != nil {
			return nil
		}
		result = append(result,&news)
	}

	return &result
}

func GetQuestionBySenderUID(skip int, limit int, isSolved bool, uid string) *[]*News {
	checkQuesCollection()

	if skip <0 || limit <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"issolved":isSolved,
		"senderuid":uid,
	}
	option := options.Find().SetSort(bson.M{"sendtime":1}).SetSkip(int64(skip)).SetLimit(int64(limit))
	
	var result []*News
	cur,err := quesCollection.Find(ctx,filter,option)
	if err!=nil{
		return nil
	}

	for cur.Next(ctx){
		var news News
		if err := cur.Decode(&news); err != nil {
			return nil
		}
		result = append(result,&news)
	}

	return &result
}

func GetQuestionByQueserUID(skip int, limit int, isSolved bool, uid string) *[]*News {
	checkQuesCollection()

	if skip <0 || limit <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"issolved":isSolved,
		"queseruid":uid,
	}
	option := options.Find().SetSort(bson.M{"sendtime":1}).SetSkip(int64(skip)).SetLimit(int64(limit))
	
	var result []*News
	cur,err := quesCollection.Find(ctx,filter,option)
	if err!=nil{
		return nil
	}

	for cur.Next(ctx){
		var news News
		if err := cur.Decode(&news); err != nil {
			return nil
		}
		result = append(result,&news)
	}

	return &result
}

func GetQuestionByClassname(skip int, limit int, isSolved bool, className string) *[]*News {
	checkQuesCollection()

	if skip <0 || limit <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"issolved":isSolved,
		"classname":className,
	}
	option := options.Find().SetSort(bson.M{"sendtime":1}).SetSkip(int64(skip)).SetLimit(int64(limit))
	
	var result []*News
	cur,err := quesCollection.Find(ctx,filter,option)
	if err!=nil{
		return nil
	}

	for cur.Next(ctx){
		var news News
		if err := cur.Decode(&news); err != nil {
			return nil
		}
		result = append(result,&news)
	}

	return &result
}

func DeleteQuestionByInnerID(id primitive.ObjectID) bool {
	checkQuesCollection()

	if id.IsZero() {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id":id}

	_,err := quesCollection.DeleteOne(ctx,filter)
	if err!=nil {
		fmt.Println(err)
		return false
	}

	return true
}