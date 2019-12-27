package edumodel

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var classCollection *mongo.Collection

type Class struct {
	ClassName   string    `bson:"classname"`
	TeacherList []string  `bson:"teacherlist"`
	StudentList []string  `bson:"studentlist"`
	CreateDate  time.Time `bson:"createdate"`
}

func checkClassCollection() {
	if classCollection == nil {
		classCollection = GetCollection("class")
	}
}

func AddClass(newClass *Class) bool {
	checkClassCollection()

	if newClass == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := classCollection.InsertOne(ctx, newClass)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return false
	}

	return true
}

func GetClassByOrder(skip int, limit int) *[]Class {
	checkClassCollection()
	if skip < 0 || limit <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"": ""}
	option := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	var result []Class
	cur, err := classCollection.Find(ctx, filter, option)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return nil
	}

	for cur.Next(ctx) {
		var class Class
		if err := cur.Decode(&class); err != nil {
			fmt.Println("[MODEL]", err)
			return nil
		}
		result = append(result, class)
	}

	return &result
}

func GetClassByName(className string) *Class {
	checkClassCollection()

	if className == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"classname": className}

	var result Class
	err := classCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return nil
	}

	return &result
}

func GetClassByUID(uid string, place string) *Class {
	checkClassCollection()

	if uid == "" || place == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var filter interface{}

	if place == "teacher" {
		filter = bson.M{"teacherlist": uid}
	} else if place == "student" {
		filter = bson.M{"studentlist": uid}
	}

	var result Class
	err := classCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return nil
	}

	return &result
}

func CheckUserInClass(className string, uid string, place string) bool {
	checkClassCollection()

	if uid == "" || place == "" || className == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var filter interface{}

	if place == "teacher" {
		filter = bson.M{"classname": className, "teacherlist": uid}
	} else if place == "student" {
		filter = bson.M{"classname": className, "studentlist": uid}
	}

	count, err := classCollection.CountDocuments(ctx, filter)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return false
	}

	if count == 0 {
		return false
	} else {
		return true
	}
}

func UpdateClassStudentByUID(className string, studentList []string) bool {
	checkClassCollection()

	if className == "" || len(studentList) == 0 {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"classname": className}
	update := bson.D{
		{"$addToSet", bson.D{
			{"studentlist", bson.D{
				{"$each", studentList},
			},
			}},
		}}

	_, err := classCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return false
	}

	return true
}

func UpdateClassTeacherByUID(className string, teacherList []string) bool {
	checkClassCollection()

	if className == "" || len(teacherList) == 0 {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"classname": className}
	update := bson.D{
		{"$addToSet", bson.D{
			{"teacherlist", bson.D{
				{"$each", teacherList},
			},
			}},
		}}

	_, err := classCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return false
	}

	return true
}

func DeleteClassStudentByUID(className string, studentList []string) bool {
	checkClassCollection()

	if className == "" || len(studentList) == 0 {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"classname": className}
	update := bson.D{
		{"$pullAll", bson.D{{"studentlist", studentList}}}}

	_, err := classCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return false
	}

	return true
}

func DeleteClassTeacherByUID(className string, teacherList []string) bool {
	checkClassCollection()

	if className == "" || len(teacherList) == 0 {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"classname": className}
	update := bson.D{
		{"$pullAll", bson.D{{"teacherlist", teacherList}}}}

	_, err := classCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return false || len(teacherList) == 0
	}

	return true
}

func DeleteClassByName(className string) bool {
	checkClassCollection()

	if className == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"classname": className}

	_, err := classCollection.DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return false
	}

	return true
}
