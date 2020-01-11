package main

import (
	"eduX/edumodel"
	"eduX/edunet"
	"eduX/edurouter"
)

func main() {
	edumodel.ConnectMongo()

	edumodel.ConnectDatabase(nil)

	s := edunet.NewServer()
	sFile := edunet.NewServer()

	s.AddRouter(0, &edurouter.LoginRouter{})
	s.AddRouter(1, &edurouter.LogoutRouter{})
	s.AddRouter(2, &edurouter.RegisterRouter{})
	s.AddRouter(3, &edurouter.PingRouter{})

	s.AddRouter(11, &edurouter.ClassAddRouter{})
	s.AddRouter(12, &edurouter.ClassDelRouter{})
	s.AddRouter(13, &edurouter.ClassJoinInGetRouter{})
	s.AddRouter(14, &edurouter.ClassListGetRouter{})
	s.AddRouter(15, &edurouter.ClassStudentAddRouter{})
	s.AddRouter(16, &edurouter.ClassStudentDelRouter{})
	s.AddRouter(17, &edurouter.ClassSetAlterNameRouter{})
	s.AddRouter(18, &edurouter.ClassCountRouter{})

	s.AddRouter(21, &edurouter.FileCountRouter{})
	s.AddRouter(22, &edurouter.FileDeleteRouter{})
	s.AddRouter(23, &edurouter.FileDownloadRouter{})
	s.AddRouter(24, &edurouter.FileGetByClassNameRouter{})
	s.AddRouter(25, &edurouter.FileGetBySenderUIDRouter{})
	s.AddRouter(26, &edurouter.FileGetByTagsRouter{})
	s.AddRouter(27, &edurouter.FileAddRouter{})

	s.AddRouter(31, &edurouter.NewsAddRouter{})
	s.AddRouter(32, &edurouter.NewsDeleteRouter{})
	s.AddRouter(33, &edurouter.NewsGetByAudientUIDRouter{})
	s.AddRouter(34, &edurouter.NewsGetBySenderUIDRouter{})
	s.AddRouter(35, &edurouter.NewsGetByTimeOrderRouter{})
	s.AddRouter(36, &edurouter.NewsCountRouter{})

	s.AddRouter(41, &edurouter.PersonAddRouter{})
	s.AddRouter(42, &edurouter.PersonInfoGetByClassRouter{})
	s.AddRouter(43, &edurouter.PersonInfoGetRouter{})
	s.AddRouter(44, &edurouter.PersonInfoPutRouter{})
	s.AddRouter(45, &edurouter.PersonCountRouter{})
	s.AddRouter(46	, &edurouter.PersonInfoGetAllRouter{})

	s.AddRouter(51, &edurouter.PwdGetQuestionRouter{})
	s.AddRouter(52, &edurouter.PwdForgetRouter{})
	s.AddRouter(53, &edurouter.PwdResetRouter{})
	s.AddRouter(54, &edurouter.PwdSetQuestionRouter{})

	s.AddRouter(61, &edurouter.QuestionAddRouter{})
	s.AddRouter(62, &edurouter.QuestionAnswerRouter{})
	s.AddRouter(63, &edurouter.QuestionCountRouter{})
	s.AddRouter(64, &edurouter.QuestionDeleteRouter{})
	s.AddRouter(65, &edurouter.QuestionGetByClassNameRouter{})
	s.AddRouter(66, &edurouter.QuestionGetBySenderUIDRouter{})

	go s.Serve()
	go sFile.ServeFile()

	select {}
}
