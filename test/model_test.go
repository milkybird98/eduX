package edutest

import(
	"eduX/edumodel"
	"fmt"
	"testing"
)

func TestClassModel(t *testing.T){
	fmt.Println("start connect mongo")
	edumodel.ConnectMongo()

	fmt.Println("start choose database")
	edumodel.ConnectDatabase(nil)

	teacherList := []string{"T1001","T1002"}
	studentList := []string{"U1001","U1002"}


	fmt.Println("add class")
	class := &edumodel.Class{
		"ts1001",
		teacherList,
		studentList,
	}
	if edumodel.AddClass(class) {
		fmt.Println("success")
	}

	fmt.Println("get class by name")
	aclass := edumodel.GetClassByName("ts1001")
	if aclass != nil{
		fmt.Println(aclass)
		fmt.Println("success")
	}

	aclass = nil

	fmt.Println("get class by teacher uid")
	aclass = edumodel.GetClassByUID("T1001","teacher")
	if aclass != nil{
		fmt.Println(aclass)
		fmt.Println("success")
	}

	aclass = nil

	fmt.Println("get class by student uid")
	aclass = edumodel.GetClassByUID("U1001","student")
	if aclass != nil{
		fmt.Println(aclass)
		fmt.Println("success")
	}

	fmt.Println("add new student")
	if edumodel.UpdateClassStudentByUID("ts1001",[]string{"U1001","U1003"}){
		fmt.Println("success")
	}

	fmt.Println("add new teacher")
	if edumodel.UpdateClassTeacherByUID("ts1001",[]string{"T1001","T1003"}){
		fmt.Println("success")
	}

	
	fmt.Println("delete student")
	if edumodel.DeleteClassStudentByUID("ts1001",[]string{"U1002"}){
		fmt.Println("success")
	}

	fmt.Println("delete teacher")
	if edumodel.DeleteClassTeacherByUID("ts1001",[]string{"T1002"}){
		fmt.Println("success")
	}

	
	fmt.Println("delete class")
	if edumodel.DeleteClassByName("ts1001"){
		fmt.Println("success")
	}
	
}

func TestUserModel(t *testing.T){
	fmt.Println("start connect mongo")
	edumodel.ConnectMongo()

	fmt.Println("start choose database")
	edumodel.ConnectDatabase(nil)

	fmt.Println("add first user")
	user := &edumodel.User{
		"test",
		"U1000",
		"123123",
		"stu",
		"ds1233",
		1,
	}
	edumodel.AddUser(user)
	
	fmt.Println("get user by id")
	user = edumodel.GetUserByUID("U1000")
	fmt.Println(user)

	fmt.Println("update user by id")
	edumodel.UpdateUserByID("U1000","ts1002","","",0)
	user = edumodel.GetUserByUID("U1000")
	fmt.Println(user)

	fmt.Println("delete user by id")
	res := edumodel.DeleteUserByUID("U1000")
	if res {
		fmt.Println("success")
	}
}
