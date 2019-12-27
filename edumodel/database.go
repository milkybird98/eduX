package edumodel

import (
	"context"
	"eduX/utils"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

/*
	创建客户端，连接MongoDB服务器
*/
func ConnectMongo() bool {
	//创建MongoDB客户端
	client, err := mongo.NewClient(options.Client().ApplyURI(utils.GlobalObject.DataBaseUrl))
	if err != nil {
		fmt.Println("[MODEL] Create MongoDB client failed, error: ", err)
		return false
	}

	//客户端尝试连接数据库
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println("[MODEL] Try connecting Mongodb Server at ", utils.GlobalObject.DataBaseUrl)
	err = client.Connect(ctx)

	//连接失败处理
	if err != nil {
		fmt.Println("[MODEL] Client connects MongoDB server failed, error: ", err)
		return false
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("[MODEL] Client ping MongoDB server failed, error: ", err)
		return false
	}
	//连接成功,保存客户端对象
	Client = client
	fmt.Println("[MODEL] Connect MongoDB server successfully")
	return true
}

/*
	客户端连接数据库
*/
func ConnectDatabase(databaseName *string) bool {
	//如未指定数据库名称则从GlobalObject中读取默认值
	if databaseName == nil {
		databaseName = &utils.GlobalObject.DataBaseName
		fmt.Println("[MODEL] Use default database name: ", utils.GlobalObject.DataBaseName)
	}

	//连接数据库
	fmt.Println("[MODEL] Try choosing database named ", *databaseName)
	db := Client.Database(*databaseName)

	//连接数据库失败处理
	if db == nil {
		fmt.Println("[MODEL] Choose database failed, target name: ", *databaseName)
		return false
	}

	//保存数据库对象
	Database = db
	fmt.Println("[MODEL] Choose database successfully")
	return true
}

/*
	根据collection名称获取的collection对象
*/
func GetCollection(collectionName string) *mongo.Collection {
	if collectionName == "" {
		fmt.Println("[MODEL] Get database collection failed, collection name CANNOT be empty")
		return nil
	}

	collection := Database.Collection(collectionName)

	if collection == nil {
		fmt.Println("[MODEL] Get collections failed")
		return nil
	}

	return collection
}
