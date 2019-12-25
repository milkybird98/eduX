package edumodel

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var quesCollection *mongo.Collection

type Question struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Title      string
	Text       string
	SenderUID  string
	AnswerUID  string `bson:"answeruid,omitempty"`
	ClassName  string
	SendTime   time.Time
	AnswerTime time.Time `bson:"answertime,omitempty"`
	IsSolved   bool
	Answer     string `bson:"answer,omitempty"`
	IsDeleted  bool   `bson:",omitempty"`
}

func checkQuesCollection() {
	if quesCollection == nil {
		quesCollection = GetCollection("question")
	}
}

func AddQuestion(newQuestion *Question) bool {
	checkQuesCollection()

	if newQuestion == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := quesCollection.InsertOne(ctx, newQuestion)
	if err != nil {
		fmt.Println("Add new Question into database fail, error: ", err)
		return false
	}

	return true
}

func GetQuestionByTimeOrder(skip int, limit int, isSolved bool) *[]Question {
	checkQuesCollection()

	if skip < 0 || limit <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"issolved": isSolved}
	option := options.Find().SetSort(bson.M{"sendtime": 1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	var result []Question
	cur, err := quesCollection.Find(ctx, filter, option)
	if err != nil {
		return nil
	}

	for cur.Next(ctx) {
		var question Question
		if err := cur.Decode(&question); err != nil {
			return nil
		}
		result = append(result, question)
	}

	return &result
}

func GetQuestionBySenderUID(skip int, limit int, detectSolved bool, isSolved bool, uid string) *[]Question {
	checkQuesCollection()

	if skip < 0 || limit <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var filter interface{}

	if detectSolved {
		filter = bson.M{
			"issolved":  isSolved,
			"senderuid": uid,
			"isdeleted": false,
		}
	} else {
		filter = bson.M{
			"senderuid": uid,
			"isdeleted": false,
		}
	}

	option := options.Find().SetSort(bson.M{"sendtime": 1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	var result []Question
	cur, err := quesCollection.Find(ctx, filter, option)
	if err != nil {
		return nil
	}

	for cur.Next(ctx) {
		var question Question
		if err := cur.Decode(&question); err != nil {
			return nil
		}
		result = append(result, question)
	}

	return &result
}

func GetQuestionByQueserUID(skip int, limit int, isSolved bool, uid string) *[]Question {
	checkQuesCollection()

	if skip < 0 || limit <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"issolved":  isSolved,
		"queseruid": uid,
		"isdeleted": false,
	}
	option := options.Find().SetSort(bson.M{"sendtime": 1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	var result []Question
	cur, err := quesCollection.Find(ctx, filter, option)
	if err != nil {
		return nil
	}

	for cur.Next(ctx) {
		var question Question
		if err := cur.Decode(&question); err != nil {
			return nil
		}
		result = append(result, question)
	}

	return &result
}

func GetQuestionByClassname(skip int, limit int, detectSolved bool, isSolved bool, className string) *[]Question {
	checkQuesCollection()

	if skip < 0 || limit <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var filter interface{}

	if detectSolved {
		filter = bson.M{
			"issolved":  isSolved,
			"classname": className,
			"isdeleted": false,
		}
	} else {
		filter = bson.M{
			"classname": className,
			"isdeleted": false,
		}
	}

	option := options.Find().SetSort(bson.M{"sendtime": 1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	var result []Question
	cur, err := quesCollection.Find(ctx, filter, option)
	if err != nil {
		return nil
	}

	for cur.Next(ctx) {
		var question Question
		if err := cur.Decode(&question); err != nil {
			return nil
		}
		result = append(result, question)
	}

	return &result
}

func GetQuestionByInnerID(idInString string) *Question {
	checkQuesCollection()

	if idInString == "" {
		return nil
	}

	id, err := primitive.ObjectIDFromHex(idInString)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id, "isdeleted": false}

	var question Question
	err = quesCollection.FindOne(ctx, filter).Decode(&question)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &question
}

func AnserQuestionByInnerID(idInString string, UID string, answer string) bool {
	checkQuesCollection()

	if idInString == "" {
		return false
	}

	id, err := primitive.ObjectIDFromHex(idInString)

	if err != nil {
		fmt.Println(err)
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id, "isdeleted": false}
	update := bson.M{"$set": bson.M{"answeruid": UID, "issolved": true, "answer": answer, "answertime": time.Now()}}

	_, err = quesCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func DeleteQuestionByInnerID(idInString string) bool {
	checkQuesCollection()

	id, err := primitive.ObjectIDFromHex(idInString)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if id.IsZero() {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id, "isdeleted": false}
	update := bson.M{"$set": bson.M{"isdeleted": true}}

	_, err = quesCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
