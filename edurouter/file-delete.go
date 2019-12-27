package edurouter

import (
	"eduX/eduiface"
	"eduX/edumodel"
	"eduX/edunet"
	"fmt"
	"os"
	"time"

	"github.com/tidwall/gjson"
)

// FileDeleteRouter 负责文件删除业务的路由
type FileDeleteRouter struct {
	edunet.BaseRouter
}

// FileDeleteData 用于客户端指定需要删除的文件
type FileDeleteData struct {
	FileID string `json:"id"`
}

var filedeleteReplyData string

// PreHandle 负责进行数据验证,权限验证和删除`数据库操作
func (router *FileDeleteRouter) PreHandle(request eduiface.IRequest) {
	var reqMsgInJSON *ReqMsg
	var ok bool
	reqMsgInJSON, filedeleteReplyData, ok = CheckMsgFormat(request)
	if ok != true {
		return
	}

	filedeleteReplyData, ok = CheckConnectionLogin(request, reqMsgInJSON.UID)
	if ok != true {
		return
	}

	if !gjson.Valid(string(reqMsgInJSON.Data)) {
		filedeleteReplyData = "data_format_error"
		return
	}

	innerIDData := gjson.GetBytes(reqMsgInJSON.Data, "id")
	if !innerIDData.Exists() {
		filedeleteReplyData = "fileid_cannot_be_empty"
		return
	}

	innerIDString := innerIDData.String()

	//权限检查
	c := request.GetConnection()
	sessionUID, err := c.GetSession("UID")
	if err != nil {
		filedeleteReplyData = "session_error"
		return
	}

	UID, ok := sessionUID.(string)
	if !ok {
		filedeleteReplyData = "session_error"
		return
	}

	sessionPlace, err := c.GetSession("place")
	if err != nil {
		filedeleteReplyData = "session_error"
		return
	}

	fileData := edumodel.GetFileByUUID(innerIDString)
	if fileData == nil {
		filedeleteReplyData = "file_not_found"
		return
	}

	if sessionPlace != "manager" && UID != fileData.UpdaterUID {
		filedeleteReplyData = "permission_error"
		return
	}

	//
	path := "./file/" + innerIDString
	err = os.Remove(path)
	if err != nil {
		ok := edumodel.DeleteFileByUUID(innerIDString)
		if !ok {
			filedeleteReplyData = "model_fail"
		} else {
			filedeleteReplyData = "success"
		}
	} else {
		filedeleteReplyData = "remove_file_error"
	}
}

// Handle 负责将处理结果发回客户端
func (router *FileDeleteRouter) Handle(request eduiface.IRequest) {
	fmt.Println("[ROUTER] ",time.Now().Format("2006-01-01 Jan 2 15:04:05"), ", Client Address: ", request.GetConnection().GetTCPConnection().RemoteAddr(), ", FileDeleteRouter: ", filedeleteReplyData)
	jsonMsg, err := CombineReplyMsg(filedeleteReplyData, nil)
	if err != nil {
		fmt.Println("FileDeleteRouter: ", err)
		return
	}

	c := request.GetConnection()
	c.SendMsg(request.GetMsgID(), jsonMsg)
}
