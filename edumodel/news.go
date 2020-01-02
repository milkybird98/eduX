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
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	IsAnnounce bool               `bson:"isan"`
	Title      string             `bson:"title"`
	Text       string             `bson:"text"`
	SenderUID  string             `bson:"senuid"`
	SendTime   time.Time          `bson:"sendtime"`
	AudientUID []string           `bson:"audiuid"`
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
		fmt.Println("[MODEL] Add new News into database fail, error: ", err)
		return false
	}

	return true
}

func GetNewsByInnerID(idInString string) *News {
	checkNewsCollection()

	if idInString == "" {
		return nil
	}

	id, err := primitive.ObjectIDFromHex(idInString)

	if err != nil {
		fmt.Println("[MODEL]", err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id, "isdeleted": false}

	var news News
	err = newsCollection.FindOne(ctx, filter).Decode(&news)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return nil
	}

	return &news
}

func GetNewsByTimeOrder(skip int, limit int, isAnnounce bool) *[]News {
	checkNewsCollection()

	if skip < 0 || limit <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"isannounce": isAnnounce}
	option := options.Find().SetSort(bson.M{"sendtime": 1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	var result []News
	cur, err := newsCollection.Find(ctx, filter, option)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return nil
	}

	for cur.Next(ctx) {
		var news News
		if err := cur.Decode(&news); err != nil {
			fmt.Println("[MODEL]", err)
			return nil
		}
		result = append(result, news)
	}

	return &result
}

func GetNewsBySenderUID(skip int, limit int, isAnnounce bool, uid string) *[]News {
	checkNewsCollection()

	if skip < 0 || limit <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"isannounce": isAnnounce,
		"senderuid":  uid,
	}
	option := options.Find().SetSort(bson.M{"sendtime": 1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	var result []News
	cur, err := newsCollection.Find(ctx, filter, option)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return nil
	}

	for cur.Next(ctx) {
		var news News
		if err := cur.Decode(&news); err != nil {
			fmt.Println("[MODEL]", err)
			return nil
		}
		result = append(result, news)
	}

	return &result
}

func GetNewsByAudientUID(skip int, limit int, isAnnounce bool, uid string) *[]News {
	checkNewsCollection()

	if skip < 0 || limit <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"isannounce": isAnnounce,
		"audientuid": []string{"all", uid},
	}
	option := options.Find().SetSort(bson.M{"sendtime": 1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	var result []News
	cur, err := newsCollection.Find(ctx, filter, option)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return nil
	}

	for cur.Next(ctx) {
		var news News
		if err := cur.Decode(&news); err != nil {
			fmt.Println("[MODEL]", err)
			return nil
		}
		result = append(result, news)
	}

	return &result
}

func DeleteNewsByInnerID(idInString string) bool {
	checkNewsCollection()

	id, err := primitive.ObjectIDFromHex(idInString)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return false
	}

	if id.IsZero() {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	_, err = newsCollection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return false
	}

	return true
}
