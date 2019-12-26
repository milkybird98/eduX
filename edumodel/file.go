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
		fmt.Println(err)
		return false
	}

	return true
}

func GetFileBySenderUID(skip int, limit int, SenderUID string) *[]File {
	checkFileCollection()

	if SenderUID == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"classname": SenderUID}
	option := options.Find().SetSort(bson.M{"updatetime": 1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	var result []File
	cur, err := fileCollection.Find(ctx, filter, option)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for cur.Next(ctx) {
		var file File
		if err := cur.Decode(&file); err != nil {
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
		fmt.Println(err)
		return nil
	}

	for cur.Next(ctx) {
		var file File
		if err := cur.Decode(&file); err != nil {
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
		fmt.Println(err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	var result File
	err = fileCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &result
}

func DeleteFileByUUID(uuidInString string) bool {
	checkFileCollection()

	id, err := primitive.ObjectIDFromHex(uuidInString)
	if err != nil {
		fmt.Println(err)
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	_, err = fileCollection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
