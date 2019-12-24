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

var newsCollection *mongo.Collection

type News struct {
	_id        primitive.ObjectID `bson:"_id,omitempty"`
	IsAnnoce   bool
	Title      string
	Text       string
	SenderUID  string
	SendTime   time.Time
	AudientUID []string
}

func checkNewsCollection() {
	if newsCollection == nil {
		newsCollection = GetCollection("news")
	}
}

func AddNews(newNews *News) bool {
	checkNewsCollection()

	if newNews == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := newsCollection.InsertOne(ctx, newNews)
	if err != nil {
		fmt.Println("Add new News into database fail, error: ", err)
		return false
	}

	return true
}

func GetNewsByTimeOrder(skip int, limit int, isAnnoce bool) *[]*News {
	checkNewsCollection()

	if skip < 0 || limit <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"isannoce": isAnnoce}
	option := options.Find().SetSort(bson.M{"sendtime": 1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	var result []*News
	cur, err := newsCollection.Find(ctx, filter, option)
	if err != nil {
		return nil
	}

	for cur.Next(ctx) {
		var news News
		if err := cur.Decode(&news); err != nil {
			return nil
		}
		result = append(result, &news)
	}

	return &result
}

func GetNewsBySenderUID(skip int, limit int, isAnnoce bool, uid string) *[]*News {
	checkNewsCollection()

	if skip < 0 || limit <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"isannoce":  isAnnoce,
		"senderuid": uid,
	}
	option := options.Find().SetSort(bson.M{"sendtime": 1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	var result []*News
	cur, err := newsCollection.Find(ctx, filter, option)
	if err != nil {
		return nil
	}

	for cur.Next(ctx) {
		var news News
		if err := cur.Decode(&news); err != nil {
			return nil
		}
		result = append(result, &news)
	}

	return &result
}

func GetNewsByAudientUID(skip int, limit int, isAnnoce bool, uid string) *[]*News {
	checkNewsCollection()

	if skip < 0 || limit <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"isannoce":   isAnnoce,
		"audientuid": uid,
	}
	option := options.Find().SetSort(bson.M{"sendtime": 1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	var result []*News
	cur, err := newsCollection.Find(ctx, filter, option)
	if err != nil {
		return nil
	}

	for cur.Next(ctx) {
		var news News
		if err := cur.Decode(&news); err != nil {
			return nil
		}
		result = append(result, &news)
	}

	return &result
}

func DeleteNewsByInnerID(id primitive.ObjectID) bool {
	checkNewsCollection()

	if id.IsZero() {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	_, err := newsCollection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
