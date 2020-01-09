package test

import (
	"eduX/edumodel"
	"eduX/edunet"
	"eduX/edurouter"
	"testing"
)

func TestServerEmpty(t *testing.T) {
	//创建一个server句柄
	edumodel.ConnectMongo()

	edumodel.ConnectDatabase(nil)

	s := edunet.NewServer()
	sFile := edunet.NewServer()

	s.AddRouter(0, &edurouter.LoginRouter{})
	s.AddRouter(27, &edurouter.FileAddRouter{})
	s.AddRouter(23, &edurouter.FileDownloadRouter{})

	go s.Serve()
	go sFile.ServeFile()

	for a := 2; a < 10; a++ {
		a = 1
	}
}
