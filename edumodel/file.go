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

var fileCollection *mongo.Collection

type File struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	FileName   string             `bson:"filename"`
	FileTag    []string           `bson:"filetag"`
	ClassName  string             `bson:"classname"`
	Size       uint64             `bson:"size"`
	UpdaterUID string             `bson:"updateruid"`
	UpdateTime time.Time          `bson:"updatetime"`
}

func checkFileCollection() {
	if fileCollection == nil {
		fileCollection = GetCollection("file")
	}
}

func AddFile(newFile *File) bool {
	checkFileCollection()

	if newFile == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := fileCollection.InsertOne(ctx, newFile)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return false
	}

	return true
}

func GetFileByTags(skip int, limit int, Tag []string, ClassName string) *[]File {
	checkFileCollection()

	if len(Tag) <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"filetag": Tag, "classname": ClassName}
	option := options.Find().SetSort(bson.M{"updatetime": 1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	var result []File
	cur, err := fileCollection.Find(ctx, filter, option)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return nil
	}

	for cur.Next(ctx) {
		var file File
		if err := cur.Decode(&file); err != nil {
			fmt.Println("[MODEL]", err)
			return nil
		}
		result = append(result, file)
	}

	return &result
}

func GetFileBySenderUID(skip int, limit int, SenderUID string) *[]File {
	checkFileCollection()

	if SenderUID == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"updateruid": SenderUID}
	option := options.Find().SetSort(bson.M{"updatetime": 1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	var result []File
	cur, err := fileCollection.Find(ctx, filter, option)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return nil
	}

	for cur.Next(ctx) {
		var file File
		if err := cur.Decode(&file); err != nil {
			fmt.Println("[MODEL]", err)
			return nil
		}
		result = append(result, file)
	}

	return &result
}

func GetFileByClassName(skip int, limit int, ClassName string) *[]File {
	checkFileCollection()

	if ClassName == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"classname": ClassName}
	option := options.Find().SetSort(bson.M{"updatetime": 1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	var result []File
	cur, err := fileCollection.Find(ctx, filter, option)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return nil
	}

	for cur.Next(ctx) {
		var file File
		if err := cur.Decode(&file); err != nil {
			fmt.Println("[MODEL]", err)
			return nil
		}
		result = append(result, file)
	}

	return &result
}

func GetFileByUUID(uuidInString string) *File {
	checkFileCollection()

	id, err := primitive.ObjectIDFromHex(uuidInString)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	var result File
	err = fileCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return nil
	}

	return &result
}

func GetFileNumberAll(className string) int {
	checkFileCollection()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"classname": className}

	count, err := fileCollection.CountDocuments(ctx, filter)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return -1
	}

	return int(count)
}

func GetFileNumberByDate(className string, targetDate time.Time) int {
	checkFileCollection()

	if targetDate.IsZero() || targetDate.After(time.Now().Add(time.Hour*24)) {
		fmt.Println("[MODEL] time out of range")
		return -1
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	targetDateInDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, time.Local)
	targetNextDateInDay := targetDateInDay.Add(time.Hour * 24)

	filter := bson.M{"classname": className,
		"updatetime": bson.M{"$gt": targetDateInDay, "$lt": targetNextDateInDay}}

	count, err := fileCollection.CountDocuments(ctx, filter)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return -1
	}

	return int(count)
}

func DeleteFileByUUID(uuidInString string) bool {
	checkFileCollection()

	id, err := primitive.ObjectIDFromHex(uuidInString)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	_, err = fileCollection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return false
	}

	return true
}
