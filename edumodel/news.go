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
	NewsType   int64              `bson:"type"       json:"type"`
	Title      string             `bson:"title"      json:"title"`
	Text       string             `bson:"text"       json:"text"`
	SenderUID  string             `bson:"senduid"    json:"senduid"`
	SendTime   time.Time          `bson:"sendtime"   json:"sendtime"`
	AudientUID []string           `bson:"audiuid"    json:"audiuid"`
	TargetTime time.Time          `bson:"targettime" json:"targettime"`
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

func GetNewsByTimeOrder(skip int, limit int, newsType int64) *[]News {
	checkNewsCollection()

	if skip < 0 || limit <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"type": newsType}
	option := options.Find().SetSort(bson.M{"sendtime": -1}).SetSkip(int64(skip)).SetLimit(int64(limit))

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

func GetNewsBySenderUID(skip int, limit int, newsType int64, uid string) *[]News {
	checkNewsCollection()

	if skip < 0 || limit <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"type":    newsType,
		"senduid": uid,
	}
	option := options.Find().SetSort(bson.M{"sendtime": -1}).SetSkip(int64(skip)).SetLimit(int64(limit))

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

func GetNewsByAudientUID(skip int, limit int, newsType int64, uid string, isAdmin bool) *[]News {
	checkNewsCollection()

	if skip < 0 || limit <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var filter interface{}

	if isAdmin {
		filter = bson.M{
			"type":    newsType,
			"audiuid": bson.M{"$in": []string{uid}},
		}
	} else {
		filter = bson.M{
			"type":       newsType,
			"audiuid":    bson.M{"$in": []string{uid}},
			"targettime": bson.M{"$lt": time.Now()},
		}
	}

	option := options.Find().SetSort(bson.M{"sendtime": -1}).SetSkip(int64(skip)).SetLimit(int64(limit))

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

func GetNewsNumberBySendUID(sendUID string) int {
	checkNewsCollection()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"senduid": sendUID}

	count, err := newsCollection.CountDocuments(ctx, filter)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return -1
	}

	return int(count)
}

func GetNewsNumberByAudientUID(audiUID string) int {
	checkNewsCollection()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"audiuid": audiUID}

	count, err := newsCollection.CountDocuments(ctx, filter)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return -1
	}

	return int(count)
}

func GetNewsNumber(audiUID string, sendUID string, newsType int) int {
	checkNewsCollection()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{}

	if audiUID != "" {
		filter = append(filter, bson.E{"audiuid", audiUID})
	}
	if sendUID != "" {
		filter = append(filter, bson.E{"senduid", sendUID})
	}
	if newsType >= 1 && newsType <= 4 {
		filter = append(filter, bson.E{"type", newsType})
	}

	count, err := newsCollection.CountDocuments(ctx, filter)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return -1
	}

	return int(count)
}

func GetNewsNumberByNewsType(newsType string) int {
	checkNewsCollection()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"type": newsType}

	count, err := newsCollection.CountDocuments(ctx, filter)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return -1
	}

	return int(count)
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
