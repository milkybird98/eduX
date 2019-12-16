package edurouter

import (
	"fmt"
	"eduX/eduiface"
	"eduX/edunet"
	"eduX/edumodel"
)

var passwordData string 
var passwordCorrect bool

type LoginRouter struct {
	edunet.BaseRouter
}

func (this *LoginRouter) PreHandle(request eduiface.IRequest) {
	
}

func (this *LoginRouter) Handle(request eduiface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

func (this *LoginRouter) PostHandle(request eduiface.IRequest) {
	if passwordCorrect==true {
		request.GetConnection().SetSession("isLogined",true)
	} else {
		request.GetConnection().SetSession("isLogined",false)
	}
}