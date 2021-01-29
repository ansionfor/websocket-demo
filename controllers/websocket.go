package controllers

import (
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/astaxie/beego"
	"demoIM/services"
	"encoding/json"
)

type WsController struct {
	BaseController
}

func (this *WsController) Connect() {
	sessionId := this.GetString("sessionId")
	if len(sessionId) == 0 {
		this.Redirect("/", 302)
		return
	}
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "not a ws handshake", 400)
		return
	} else if err != nil {
		beego.Error("can not setup ws connection:", err)
		return
	}

	// 判断用户信息
	userId := services.GetUserIdBySessionId(sessionId)
	if userId == 0 {
		data := services.ResponseMsg{Code:1,Msg:"用户不存在"}
		mJson, _ := json.Marshal(data)
		ws.WriteMessage(websocket.TextMessage, mJson)
		ws.Close()
		return
	}

	services.AddToOnlineList(userId, ws)
	// 连接断开时调用
	defer services.RemoveFromOnlineList(userId)

	// 处理消息
	services.HandleMsg(userId, ws)
}