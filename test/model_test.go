package test

import (
	"eduX/edumodel"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestQuestionModel(t *testing.T) {
	fmt.Println("start connect mongo")
	ok := edumodel.ConnectMongo()
	if !ok {
		t.Fail()
	}

	fmt.Println("start choose database")
	ok = edumodel.ConnectDatabase(nil)
	if !ok {
		t.Fail()
	}

	fmt.Println("add question")
	var newQuestion edumodel.Question
	newQuestion.Title = "测试 Question"
	newQuestion.Text = "this is a < 测试问题》"
	newQuestion.SenderUID = "U1001"
	newQuestion.ClassName = "ts1001"
	newQuestion.SendTime = time.Now().In(utils.GlobalObject.TimeLocal)
	newQuestion.IsSolved = false
	newQuestion.IsDeleted = false

	ok = edumodel.AddQuestion(&newQuestion)
	if ok {
		fmt.Println("pass")
	} else {
		t.FailNow()
	}

	fmt.Println("get question by classname")
	questions := edumodel.GetQuestionByClassName(0, 5, false, false, "ts1001")
	if questions == nil {
		t.FailNow()
	}
	fmt.Println(*questions)

	ss := (*questions)[0].ID.Hex()
	fmt.Println(ss)
	id, _ := primitive.ObjectIDFromHex(ss)
	fmt.Println(id)
	js, _ := json.Marshal((*questions)[0])
	fmt.Println(string(js))

	fmt.Println("get questions by sender UID")
	questions = edumodel.GetQuestionBySenderUID(0, 10, false, false, "U1001")
	if questions == nil {
		t.FailNow()
	}
	fmt.Println((*questions)[0])
	fmt.Println("pass")

	fmt.Println("Answer question by inner ID")
	ok = edumodel.AnserQuestionByInnerID(((*questions)[0].ID.Hex()), "T1002", "这是a test answer, 它混合中 英 文")
	if !ok {
		t.FailNow()
	}
	fmt.Println("pass")

	fmt.Println("Delete answer by innerID")
	ok = edumodel.DeleteQuestionByInnerID((*questions)[0].ID.Hex())
	if !ok {
		t.FailNow()
	}
	fmt.Println("pass")
}

func TestClassModel(t *testing.T) {
	fmt.Println("start connect mongo")
	ok := edumodel.ConnectMongo()
	if !ok {
		return
	}

	fmt.Println("start choose database")
	ok = edumodel.ConnectDatabase(nil)
	if !ok {
		return
	}

	teacherList := []string{"T1001", "T1002"}
	studentList := []string{"U1001", "U1002"}

	fmt.Println("add class")
	class := &edumodel.Class{
		"ts1001",
		teacherList,
		studentList,
		time.Now().In(utils.GlobalObject.TimeLocal),
	}
	if edumodel.AddClass(class) {
		fmt.Println("success")
	}

	fmt.Println("get class by name")
	aclass := edumodel.GetClassByName("ts1001")
	if aclass != nil {
		fmt.Println(aclass)
		fmt.Println("success")
	}

	aclass = nil

	fmt.Println("get class by teacher uid")
	aclass = edumodel.GetClassByUID("T1001", "teacher")
	if aclass != nil {
		fmt.Println(aclass)
		fmt.Println("success")
	}

	aclass = nil

	fmt.Println("get class by student uid")
	aclass = edumodel.GetClassByUID("U1001", "student")
	if aclass != nil {
		fmt.Println(aclass)
		fmt.Println("success")
	}

	fmt.Println("add new student")
	if edumodel.UpdateClassStudentByUID("ts1001", []string{"U1001", "U1003"}) {
		fmt.Println("success")
	}

	fmt.Println("add new teacher")
	if edumodel.UpdateClassTeacherByUID("ts1001", []string{"T1001", "T1003"}) {
		fmt.Println("success")
	}

	fmt.Println("delete student")
	if edumodel.DeleteClassStudentByUID("ts1001", []string{"U1002"}) {
		fmt.Println("success")
	}

	fmt.Println("delete teacher")
	if edumodel.DeleteClassTeacherByUID("ts1001", []string{"T1002"}) {
		fmt.Println("success")
	}

	fmt.Println("delete class")
	if edumodel.DeleteClassByName("ts1001") {
		fmt.Println("success")
	}

}

func TestUserModel(t *testing.T) {
	fmt.Println("start connect mongo")
	ok := edumodel.ConnectMongo()
	if !ok {
		return
	}

	fmt.Println("start choose database")
	ok = edumodel.ConnectDatabase(nil)
	if !ok {
		return
	}

	fmt.Println("add first user")
	user := &edumodel.User{
		"test",
		"U1000",
		"stu",
		"ds1233",
		1,
		"1982-12-22",
		"群众",
		"123456789",
		false,
		"",
		false,
		"Hubei",
		true,
	}
	edumodel.AddUser(user)

	fmt.Println("get user by id")
	user = edumodel.GetUserByUID("U1000")
	fmt.Println(user)

	fmt.Println("update user by id")

	edumodel.UpdateUserByID(user)
	user = edumodel.GetUserByUID("U1000")
	fmt.Println(user)

	fmt.Println("delete user by id")
	res := edumodel.DeleteUserByUID("U1000")
	if res {
		fmt.Println("success")
	}
}
