package edutest

import(
	"eduX/edumodel"
	"fmt"
	"testing"
)

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