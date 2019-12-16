package edumodel

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"fmt"
	"time"
	"github.com/google/uuid"
)

var fileCollection *mongo.Collection

type File struct{
	FileName 			string
	UUID					uuid.UUID
	Size					uint64
	UpdaterUID		string
	UpdateTime		time.Time
}

func checkFileCollection()  {
	if fileCollection == nil{
		fileCollection = GetCollection("file")
	}
}

func AddFile(fileName string,size uint64,updaterUID string) (string,bool) {
	checkFileCollection()

	uuid := uuid.Must(uuid.NewUUID())

	if fileName == "" || size == 0 || updaterUID=="" {
		return "",false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updateTime := time.Now()
	file := &File{
		fileName,
		uuid,
		size,
		updaterUID,
		updateTime,
	}

	_,err := fileCollection.InsertOne(ctx,file)
	if err != nil {
		fmt.Println(err)
		return "",false
	}

	return uuid.String(),true
}

func GetFileByUUID(uuidInString string) *File {
	checkFileCollection()

	uuid,err := uuid.Parse(uuidInString)

	if err!=nil {
		fmt.Println(err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"uuid":uuid}

	var result File
	err = fileCollection.FindOne(ctx,filter).Decode(&result)
	if err != nil{
		fmt.Println(err)
		return nil
	}

	return &result
}

func DeleteFileByUUID(uuidInString string) bool {
	checkFileCollection()

	uuid,err := uuid.Parse(uuidInString)

	if err!=nil {
		fmt.Println(err)
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"uuid":uuid}

	_,err = fileCollection.DeleteOne(ctx,filter)
	if err!=nil {
		fmt.Println(err)
		return false
	}

	return true
}
