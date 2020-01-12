package edumodel

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var quesCollection *mongo.Collection

type Question struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Title     string             `bson:"title" json:"title"`
	Text      string             `bson:"text" json:"text"`
	SenderUID string             `bson:"senduid" json:"senduid"`
	ClassName string             `bson:"class" json:"class"`
	SendTime  time.Time          `bson:"sendtime" json:"sendtime"`
	IsSolved  bool               `bson:"issolved" json:"issolved"`
	IsDeleted bool               `bson:"isdeleted" json:"isdeleted"`
	Answer    Answerlist         `bson:"answer,omitempty" json:"answer"`
}

type QuestionAnser struct {
	AnswerUID  string    `bson:"answeruid,omitempty" json:"answeruid"`
	AnswerTime time.Time `bson:"answertime,omitempty" json:"answertime"`
	AnswerText string    `bson:"text,omitempty" json:"text"`
}

type Answerlist []QuestionAnser

func (a Answerlist) Len() int {
	return len(a)
}

func (a Answerlist) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Answerlist) Less(i, j int) bool {
	return a[i].AnswerTime.After(a[j].AnswerTime)
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
		fmt.Println("[MODEL] Add new Question into database fail, error: ", err)
		return false
	}

	return true
}

func GetQuestionByTimeOrder(skip, limit int64, detectSolved bool, isSolved bool) *[]Question {
	checkQuesCollection()

	if skip < 0 || limit <= 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var filter interface{}
	if detectSolved {
		filter = bson.M{"issolved": isSolved}
	} else {
		filter = bson.M{}
	}
	option := options.Find().SetSort(bson.M{"sendtime": -1}).SetSkip(skip).SetLimit(limit)

	var result []Question
	cur, err := quesCollection.Find(ctx, filter, option)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return nil
	}

	for cur.Next(ctx) {
		var question Question
		if err := cur.Decode(&question); err != nil {
			fmt.Println("[MODEL]", err)
			return nil
		}
		sort.Sort(question.Answer)
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
			"senduid":   uid,
			"isdeleted": false,
		}
	} else {
		filter = bson.M{
			"senduid":   uid,
			"isdeleted": false,
		}
	}

	option := options.Find().SetSort(bson.M{"sendtime": -1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	var result []Question
	cur, err := quesCollection.Find(ctx, filter, option)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return nil
	}

	for cur.Next(ctx) {
		var question Question
		if err := cur.Decode(&question); err != nil {
			fmt.Println("[MODEL]", err)
			return nil
		}
		sort.Sort(question.Answer)
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
	option := options.Find().SetSort(bson.M{"sendtime": -1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	var result []Question
	cur, err := quesCollection.Find(ctx, filter, option)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return nil
	}

	for cur.Next(ctx) {
		var question Question
		if err := cur.Decode(&question); err != nil {
			fmt.Println("[MODEL]", err)
			return nil
		}
		sort.Sort(question.Answer)
		result = append(result, question)
	}

	return &result
}

func GetQuestionByClassName(skip int, limit int, detectSolved bool, isSolved bool, className string) *[]Question {
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
			"class":     className,
			"isdeleted": false,
		}
	} else {
		filter = bson.M{
			"class":     className,
			"isdeleted": false,
		}
	}

	option := options.Find().SetSort(bson.M{"sendtime": -1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	var result []Question
	cur, err := quesCollection.Find(ctx, filter, option)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return nil
	}

	for cur.Next(ctx) {
		var question Question
		if err := cur.Decode(&question); err != nil {
			fmt.Println("[MODEL]", err)
			return nil
		}
		sort.Sort(question.Answer)
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
		fmt.Println("[MODEL]", err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id, "isdeleted": false}

	var question Question
	err = quesCollection.FindOne(ctx, filter).Decode(&question)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return nil
	}

	sort.Reverse(question.Answer)
	return &question
}

func AnserQuestionByInnerID(idInString string, UID string, answer string) bool {
	checkQuesCollection()

	if idInString == "" {
		return false
	}

	id, err := primitive.ObjectIDFromHex(idInString)

	if err != nil {
		fmt.Println("[MODEL]", err)
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id, "isdeleted": false}
	update := bson.M{"$set": bson.M{"issolved": true}, "$push": bson.M{"answer": bson.M{"answeruid": UID, "text": answer, "answertime": time.Now()}}}

	_, err = quesCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return false
	}

	return true
}

func GetQuestionNumber(className, sendUID string, isSolved bool, targetDate *time.Time) int64 {
	checkQuesCollection()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var targetDateInDay, targetNextDateInDay time.Time

	if targetDate != nil {
		targetDateInDay = time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, time.Local)
		targetNextDateInDay = targetDateInDay.Add(time.Hour * 24)
	}

	filter := bson.D{}
	if className != "" {
		filter = append(filter, bson.E{"class", className})
	}
	if sendUID != "" {
		filter = append(filter, bson.E{"senduid", sendUID})
	}
	if isSolved {
		filter = append(filter, bson.E{"issolved", true})
	}
	if targetDate != nil {
		filter = append(filter, bson.E{"sendtime", bson.M{"$gt": targetDateInDay, "$lt": targetNextDateInDay}})
	}

	filter = append(filter, bson.E{"isdeleted", false})

	fmt.Println(filter)

	count, err := quesCollection.CountDocuments(ctx, filter)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return -1
	}

	return count
}

func DeleteQuestionByInnerID(idInString string) bool {
	checkQuesCollection()

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

	filter := bson.M{"_id": id, "isdeleted": false}
	update := bson.M{"$set": bson.M{"isdeleted": true}}

	_, err = quesCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("[MODEL]", err)
		return false
	}

	return true
}
