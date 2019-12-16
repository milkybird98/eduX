package edurouter

import (
	"fmt"
	"eduX/eduiface"
	"eduX/edunet"
)

var passwordData string 
var passwordCorrect bool

type LoginRouter struct {
	edunet.BaseRouter
}

//Test PreHandle
func (this *LoginRouter) PreHandle(request eduiface.IRequest) {

}

//Test Handle
func (this *LoginRouter) Handle(request eduiface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

//Test PostHandle
func (this *LoginRouter) PostHandle(request eduiface.IRequest) {
	if passwordCorrect==true {
		request.GetConnection().SetSession("isLogined",true)
	} else {
		request.GetConnection().SetSession("isLogined",false)
	}
}