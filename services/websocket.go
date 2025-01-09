package services

import (
	"time"
	"strconv"
	"encoding/json"
	"github.com/astaxie/beego"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
)

type ResponseType int

const (
	RT_SINGLE_MSG 				= 0		// 私聊，用户收到消息时响应
	RT_FRIEND_UNREAD_MSG 			= 1		// 好友列表以及未读消息汇总，用户上线后响应
	RT_HISTORY_MSG				= 2		// 历史消息，用户触发后响应
	PING = "ping"
	PONG = "pong"
	GET_FRIEND_LIST = "getFriendList"
	SEND_SUCC = 1
	SEND_FAIL = 0
)

// 响应结构体
type ResponseMsg struct {
	Code 		int    `json:"code"`
	Msg 		string `json:"msg"`
	Timestamp 	int    `json:"timestamp"`
	DataType	int	   `json:"dataType"`
}

// 客户端发起消息格式
type ClientSingleMsg struct {
	ReqType			int 	`json:"reqType"`
	ToUserId		int		`json:"toUserId"`
	FromUserId		int		`json:"fromUserId"`
	Content			string	`json:"content"`
}

// 私聊消息格式
type ResponseSingleMsg struct {
	ResponseMsg
	Data			ClientSingleMsg  	`json:"data"`
}
// 新加入的用户
type NewOnlineUser struct {
	UserId int
	Conn *websocket.Conn
}

// 聊天记录
type History struct {
	Id 				int
	FPublisherId	int
	TPublisherId	int
	Content			string
	CTime 			int
	IsReceived		int
	IsReaded		int
}

var (
	// 新上线用户
	newUserChan = make(chan NewOnlineUser, 10)
	// 退出用户
	existUserChan = make(chan int, 10)
	// 发送消息
	sendMsgChan = make(chan ClientSingleMsg, 10)
	// 响应消息结构
	responseSingleMsg = new(ResponseSingleMsg)
	// 在线列表
	onlineUsersMap = map[int] *websocket.Conn{}
	// 数据库连接
	mysqlDb, mysqlDbErr = sql.Open("mysql", beego.AppConfig.String("dsn"))
)

func init() {
	mysqlDb.SetMaxOpenConns(10)
	go listenChannel()
}

func main() {
	if mysqlDbErr != nil {
		beego.Info(mysqlDbErr)
	}
	defer mysqlDb.Close()
}

// 监听chan消息
func listenChannel() {
	for {
		select {
		case newUser := <-newUserChan:
			onlineUsersMap[newUser.UserId] = newUser.Conn
			// 通知客户端该用户上线 TODO
		case existUserId := <-existUserChan:
			delete(onlineUsersMap, existUserId)
			// 通知客户端该用户离线 TODO
		case sendMsg := <-sendMsgChan:
			sendSingleMsg(sendMsg)
		}
	}
}

// 监听消息
func HandleMsg(userId int, ws *websocket.Conn) {
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		msg := string(p)
		beego.Info(msg)
		if msg == PING {
			// 心跳检测
			ws.WriteMessage(websocket.TextMessage, []byte(PONG))
			continue
		}
		// 集中处理不同类型消息
		handleMsgType(userId, msg)
	}
}

func handleMsgType(userId int, jsonMsg string) {
	msg := make(map[string]interface{})
	json.Unmarshal([]byte(jsonMsg), &msg)
	if msg["reqType"] == nil {
		beego.Info("reqType is empty")
		return
	}
	if msg["reqType"] == float64(RT_SINGLE_MSG) {
		jsonToSingleStruct(userId, jsonMsg)
	} else if msg["reqType"] == float64(RT_FRIEND_UNREAD_MSG) {
		getFriendAndUnreadMsg(userId, msg)
	} else if msg["reqType"] == float64(RT_HISTORY_MSG) {
		getHistoryMsg(userId, msg)
	}
}

func getFriendAndUnreadMsg(userId int, msg map[string]interface{}) {

}

func getHistoryMsg(userId int, msg map[string]interface{}) {

}

// 转发私聊消息
func sendSingleMsg(sendMsg ClientSingleMsg) {
	ws, isOnline := onlineUsersMap[sendMsg.ToUserId];
	//content, _ := json.Marshal(sendMsg)
	currentTime := int(time.Now().Unix())
	responseMsg := ResponseSingleMsg{ResponseMsg{0, "success", currentTime, RT_SINGLE_MSG}, sendMsg}
	responseJson, _ := json.Marshal(responseMsg)
	if isOnline {
		// 在线
		if ws.WriteMessage(websocket.TextMessage, responseJson) != nil {
			// 离线发送失败 TODO
			go saveMsgToMysql(sendMsg.FromUserId, sendMsg.ToUserId, sendMsg.Content, SEND_FAIL)
		} else {
			// 发送成功
			go saveMsgToMysql(sendMsg.FromUserId, sendMsg.ToUserId, sendMsg.Content, SEND_SUCC)
		}
	} else {
		// 离线 TODO
		go saveMsgToMysql(sendMsg.FromUserId, sendMsg.ToUserId, sendMsg.Content, SEND_FAIL)
	}
}

// 消息入库
func saveMsgToMysql(fromUserId int, toUserId int, content string, isReceived int) {
	prepare := "INSERT INTO bk_im_history(f_publisher_id,t_publisher_id,content,is_readed,is_received,c_time) VALUES(?,?,?,?,?,?)"
	stmt, err := mysqlDb.Prepare(prepare)
	if err != nil {
		beego.Info(err)
	}
	_, err = stmt.Exec(fromUserId, toUserId, content, 0, isReceived, int(time.Now().Unix()))
	if err != nil {
		beego.Info(err)
	}
	defer stmt.Close()
}

// 处理客户端消息
func jsonToSingleStruct(fromUserId int, jsonMsg string) {
 	var singleMsg ClientSingleMsg
	json.Unmarshal([]byte(jsonMsg), &singleMsg)
	singleMsg.FromUserId = fromUserId
	sendMsgChan <- singleMsg
}

// 根据sessionid获取userid
func GetUserIdBySessionId(sessionId string) (userId int) {
	userId, _ = strconv.Atoi(sessionId)
	return userId
}

// 加入在线列表
func AddToOnlineList(userId int, ws *websocket.Conn) {
	newUserChan <- NewOnlineUser{UserId: userId, Conn: ws}
}

// 退出在线列表
func RemoveFromOnlineList(userId int)  {
	beego.Info(userId)
}
