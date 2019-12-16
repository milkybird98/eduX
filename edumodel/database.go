package edumodel

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"eduX/utils"
	"context"
	"fmt"
	"time"
)

var Client *mongo.Client
var Database *mongo.Database


/*
	创建客户端，连接MongoDB服务器
*/
func ConnectMongo()  {
	//创建MongoDB客户端
	client, err := mongo.NewClient(options.Client().ApplyURI(utils.GlobalObject.DataBaseUrl))
	if err!=nil {
		fmt.Println("Create MongoDB client failed, error: ",err)
		return
	}

	//客户端尝试连接数据库
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println("Try connecting Mongodb Server at ",utils.GlobalObject.DataBaseUrl)
	err = client.Connect(ctx)

	//连接失败处理
	if err!=nil {
		fmt.Println("Client connects MongoDB server failed, error: ",err)
		return
	}

	//连接成功,保存客户端对象
	Client = client
	fmt.Println("Connect MongoDB server successfully")
	return
}

/*
	客户端连接数据库
*/
func ConnectDatabase(databaseName *string)  {
	//如未指定数据库名称则从GlobalObject中读取默认值
	if databaseName == nil {
		databaseName = &utils.GlobalObject.DataBaseName
		fmt.Println("Use default database name: ",utils.GlobalObject.DataBaseName)
	}

	//连接数据库
	fmt.Println("Try choosing database named ",*databaseName)
	db := Client.Database(*databaseName)

	//连接数据库失败处理
	if db==nil {
		fmt.Println("Choose database failed, target name: ",*databaseName)
	}

	//保存数据库对象
	Database = db
	fmt.Println("Choose database successfully")
	return
}

/*
	根据collection名称获取的collection对象
*/
func GetCollection(collectionName string) *mongo.Collection {
	if collectionName == "" {
		fmt.Println("Get database collection failed, collection name CANNOT be empty")
		return nil
	}

	collection := Database.Collection(collectionName)
	
	if collection == nil {
		fmt.Println("Get collections failed")
		return nil
	}

	return collection
}