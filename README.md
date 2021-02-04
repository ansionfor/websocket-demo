# websocket-demo
基于beego的websocket点对点私聊demo

### 1、下载

```sql
git clone https://github.com/ansionfor/websocket-demo.git
```

### 2、运行
```sql
go run main.go

2021/02/04 19:57:10.407 [I] [main.go:9]  demoIM 1.0 
2021/02/04 19:57:10.428 [I] [asm_amd64.s:1374]  http server Running on http://:8080

```
### 3、连接
```sql
var ws = new WebSocket("ws://ip:8080/ws?sessionId=1");

ws.onopen = function(evt) { 
  console.log("Connection open ..."); 
  ws.send("ping");
};

ws.onmessage = function(evt) {
  console.log( "Received Message: " + evt.data);
  ws.close();
};

ws.onclose = function(evt) {
  console.log("Connection closed.");
}; 
```
