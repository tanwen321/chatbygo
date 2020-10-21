package gchat

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	// "strconv"
	// "strings"
	"sync"
	"time"
)

const (
	pnghlen = 22
	firtype = 1
	grptype = 2
	ontype  = 3
	offtype = 4
)

// http升级websocket协议的配置
// var wsUpgrader = websocket.Upgrader{
// 允许所有CORS跨域请求
var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 客户端读写消息
type wsMessage struct {
	messageType int
	data        []byte
}

// 客户端连接
type wsConnection struct {
	wsSocket *websocket.Conn // 底层websocket
	inChan   chan *wsMessage // 读队列
	outChan  chan *wsMessage // 写队列

	mutex     sync.Mutex // 避免重复关闭管道
	isClosed  bool
	closeChan chan byte // 关闭通知
	clientid  int
}

//消息通讯json
type jtype struct {
	Id int    `json:"id,string"`
	Na string `json:"na"`
	Im string `json:"im"`
	Wo string `json:"wo"`
	Ti string `json:"ti"`
}

type revmsg struct {
	Id    int    `json:"id,string"`
	Type  int    `json:"type"`
	Mdata string `json:"mdata"`
}

type sendmsg struct {
	Id     int    `json:"id,string"`       //群组或者好友id
	Sendid int    `json:"senderId,string"` //发送者id
	Type   int    `json:"type"`            //消息类型1：朋友/2：群组
	Mdata  string `json:"payload"`         //消息内容
}

type oafuser struct {
	Id     int    `json:"id,string"` //发送id
	Type   int    `json:"type"`      //消息类型3:上线消息/4:下线消息
	Name   string `json:"name"`      //发送名字
	Avatar string `json:"avatar"`    //发送头像
}

type juser struct {
	Uuid   int    `json:"uuid,string"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func (wsConn *wsConnection) wsReadLoop() {
	for {
		// 读一个message
		msgType, data, err := wsConn.wsSocket.ReadMessage()
		if err != nil {
			goto error
		}
		req := &wsMessage{
			msgType,
			data,
		}
		// 放入请求队列
		select {
		case wsConn.inChan <- req:
		case <-wsConn.closeChan:
			goto closed
		}
	}
error:
	off_user(wsConn.clientid)
	wsConn.wsClose()
closed:
}

func (wsConn *wsConnection) wsWriteLoop() {
	for {
		select {
		// 取一个应答
		case msg := <-wsConn.outChan:
			// 写给websocket
			if err := wsConn.wsSocket.WriteMessage(msg.messageType, msg.data); err != nil {
				goto error
			}
		case <-wsConn.closeChan:
			goto closed
		}
	}
error:
	off_user(wsConn.clientid)
	wsConn.wsClose()
closed:
}

func (wsConn *wsConnection) procLoop() {

	//发送初始信息
	msg := init_msg(wsConn.clientid)
	if err := wsConn.wsWrite(websocket.TextMessage, msg); err != nil {
		fmt.Println("write fail 1")
		off_user(wsConn.clientid)
		wsConn.wsClose()
	}
	// 这是一个同步处理模型（只是一个例子），如果希望并行处理可以每个请求一个gorutine，注意控制并发goroutine的数量!!!
	for {
		msg, err := wsConn.wsRead()
		if err != nil {
			fmt.Println("read fail 1")
			off_user(wsConn.clientid)
			break
		}

		if string(msg.data) == "ping" {
			err = wsConn.wsWrite(msg.messageType, []byte("pong"))
			if err != nil {
				fmt.Println("write fail 2")
				off_user(wsConn.clientid)
				break
			}
		} else {
			mdata := &revmsg{}
			err := json.Unmarshal(msg.data, mdata)
			if err != nil {
				fmt.Println("write fail 3")
				off_user(wsConn.clientid)
				break
			}
			send_msg(wsConn.clientid, msg.messageType, mdata)
		}
	}
}

func wsHandler(resp http.ResponseWriter, req *http.Request) {
	// 应答客户端告知升级连接为websocket
	//	resp.Header().Set("content-type", "application/json")
	wsSocket, err := wsUpgrader.Upgrade(resp, req, nil)
	if err != nil {
		return
	}
	c, err := req.Cookie("mmssid")
	if err != nil || c == nil {
		return
	}
	id := chnti[c.Value]
	client, ok := userdb[id]
	if !ok || client.ucook == nil || c == nil || client.ucook.Name != c.Name {
		return
	}
	wsConn := &wsConnection{
		wsSocket:  wsSocket,
		inChan:    make(chan *wsMessage, 1024),
		outChan:   make(chan *wsMessage, 1024),
		closeChan: make(chan byte),
		isClosed:  false,
		clientid:  id,
	}
	client.uws = wsConn
	userdb[id] = client
	onuch <- client

	// 处理器
	go wsConn.procLoop()
	// 读协程
	go wsConn.wsReadLoop()
	// 写协程
	go wsConn.wsWriteLoop()
}

func (wsConn *wsConnection) wsWrite(messageType int, data []byte) error {
	select {
	case wsConn.outChan <- &wsMessage{messageType, data}:
	case <-wsConn.closeChan:
		return errors.New("websocket closed")
	}
	return nil
}

func (wsConn *wsConnection) wsRead() (*wsMessage, error) {
	select {
	case msg := <-wsConn.inChan:
		return msg, nil
	case <-wsConn.closeChan:
	}
	return nil, errors.New("websocket closed")
}

func (wsConn *wsConnection) wsClose() {
	wsConn.wsSocket.Close()

	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.closeChan)
	}
}

func get_group_info(glist []int) gdb {
	var ngdb = make(gdb)
	for _, gid := range glist {
		g, ok := groupdb[gid]
		if ok {
			ngdb[gid] = g
		}
	}
	return ngdb
}

func get_user_info(flist []int) udb {
	var nudb = make(udb)
	for _, fid := range flist {
		if onuser[fid] {
			u, _ := userdb[fid]
			nudb[fid] = u
		}
	}
	return nudb
}

func init_msg(id int) []byte {
	client, _ := userdb[id]
	var sendjson = make(map[string](map[string]map[int]juser))
	var fmlist = make(map[string]map[int]juser)
	var j juser
	var ms = make(map[int]juser)
	j.Uuid = client.mid
	j.Avatar = client.uimage
	j.Name = client.unkname
	ms[0] = j
	fmlist["currentUser"] = ms

	glist := get_group_info(client.group)
	var gs = make(map[int]juser)
	for _, g := range glist {
		j.Uuid = g.gid
		j.Avatar = g.gimage
		j.Name = g.gname
		gs[g.gid] = j
	}
	fmlist["groups"] = gs
	flist := get_user_info(client.friends)
	var fs = make(map[int]juser)
	for _, f := range flist {
		if f.uws != nil {
			var ou oafuser
			ou.Id = id
			ou.Type = ontype
			ou.Name = client.unkname
			ou.Avatar = client.uimage
			data, err := json.Marshal(ou)
			if err != nil {
				data = []byte("pong")
			}
			err = f.uws.wsWrite(websocket.TextMessage, data)
			if err != nil {
				fmt.Println("write fail 6", err)
				onuser[f.mid] = false
				continue
			}
			j.Uuid = f.mid
			j.Avatar = f.uimage
			j.Name = f.unkname
			fs[f.mid] = j
		}
	}
	fmlist["friends"] = fs

	sendjson["init"] = fmlist
	data, err := json.Marshal(sendjson)
	if err != nil {
		return []byte("pong")
	}
	return data
}

func now_time_s() string {
	ds := time.Now().Format("01-02, 15:04")
	return string(ds)
}

func send_msg(cid, t int, rdata *revmsg) {
	if rdata.Type == grptype {
		ulist, ok := ongtou[rdata.Id]
		if !ok {
			fmt.Printf("online group error:%v\n", rdata.Id)
		}
		for _, u := range ulist {
			var sdata sendmsg
			sdata.Id = rdata.Id
			sdata.Mdata = rdata.Mdata
			sdata.Type = rdata.Type
			sdata.Sendid = cid
			data, _ := json.Marshal(sdata)
			err := u.uws.wsWrite(t, data)
			if err != nil {
				fmt.Println("write fail 4", err)
				onuser[u.mid] = false
				//				off_user(u.mid)
				continue
			}
		}
	} else if rdata.Type == firtype {
		if onuser[rdata.Id] {
			f := userdb[rdata.Id]
			var sdata sendmsg
			sdata.Id = cid
			sdata.Mdata = rdata.Mdata
			sdata.Type = rdata.Type
			sdata.Sendid = cid
			data, _ := json.Marshal(sdata)
			err := f.uws.wsWrite(t, data)
			if err != nil {
				fmt.Println("write fail 7", err)
				onuser[f.mid] = false
			}
			my := userdb[cid]
			sdata.Id = rdata.Id
			data, _ = json.Marshal(sdata)
			err = my.uws.wsWrite(t, data)
			if err != nil {
				fmt.Println("write fail 8", err)
				onuser[my.mid] = false
			}
		}
	}
}

func get_data(t, s string) string {
	if t == "3" {
		ms := md5v(s)
		if f, ok := pngtmp[ms]; ok {
			return f
		}
		bin := []byte(s)
		pngtype, data := get_type(bin)
		ds := time.Now().Format("15_04_05")
		ifname := "/png_tmp/" + string(ds) + "." + pngtype
		data2, err := base64.StdEncoding.DecodeString(string(data))
		if err != nil {
			return "pong"
		}
		err = ioutil.WriteFile(Whome+ifname, data2, 0666)
		if err != nil {
			return "pong"
		}
		pngtmp[ms] = ifname
		return ifname
	}
	return s
}

func get_type(b []byte) (string, []byte) {
	var f = false
	var t []byte
	for i := 0; i < 64; i++ {
		if !f && b[i] == 47 {
			f = true
		} else if f && b[i] == 59 {
			f = false
		} else if !f && b[i] == 44 {
			i++
			return string(t), b[i:]
		} else if f {
			t = append(t, b[i])
		}
	}
	return "png", b
}

func md5v(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func send_exit_msg(id int, fl []int) {
	flist := get_user_info(fl)
	for _, f := range flist {
		if f.uws != nil {
			var ou oafuser
			ou.Id = id
			ou.Type = offtype
			data, err := json.Marshal(ou)
			if err != nil {
				data = []byte("pong")
			}
			err = f.uws.wsWrite(websocket.TextMessage, data)
			if err != nil {
				fmt.Println("write fail 6", err)
				onuser[f.mid] = false
				continue
			}
		}
	}
}

func off_user(id int) {
	client, _ := userdb[id]
	client.ucook = nil
	client.uws = nil
	userdb[id] = client
	if onuser[id] {
		send_exit_msg(id, client.friends)
		offuch <- client
	}
}
