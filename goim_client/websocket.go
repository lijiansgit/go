package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sync/atomic"
	"time"

	log "github.com/alecthomas/log4go"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

var (
	msgTimeRegexp = regexp.MustCompile(`\d+`)
)

func websocketReadProto(conn *websocket.Conn, p *Proto) error {
	msg, _ := json.Marshal(p)
	log.Debug("%s", string(msg))
	return conn.ReadJSON(p)
}

func websocketWriteProto(conn *websocket.Conn, p *Proto) error {
	if p.Body == nil {
		p.Body = []byte("{}")
	}
	return conn.WriteJSON(p)
}

func newRoomInfo(rid, clients int) *roomInfo {
	return &roomInfo{
		rid:     rid,
		clients: clients,
	}
}

type roomInfo struct {
	rid         int
	clients     int
	Level1Count int64
	Level2Count int64
	Level3Count int64

	Level4Count int64
	Level5Count int64
	Level6Count int64

	ParseFaildCount int64

	recvCount     int64
	disConnCount  int64
	connFailCount int64
	//TODO 有规律的增加连接数、减少连接数
}

func (r *roomInfo) run() {
	for i := 0; i < r.clients; i++ {
		go r.newClient(r.rid, websocketAddr)
		time.Sleep(time.Millisecond * 1)
	}
}

func (r *roomInfo) newClient(key int, websocketAddr string) {
	var (
		params string
	)
	params = "room_id=" + fmt.Sprintf("%d", key)
	if mobile != "-1" {
		params += fmt.Sprintf("&device=%s", mobile)
	}
	if version != "-1" {
		params += fmt.Sprintf("&version=%s", version)
	}
	// params += fmt.Sprintf("&packageId=0")
	u := url.URL{Scheme: "ws", Host: websocketAddr, Path: "/", RawQuery: params}
	var (
		cookieStr string
		conn      *websocket.Conn
		err       error
	)

	wsHeader := http.Header{
		"Origin": {websocketAddr},
		"Cookie": {cookieStr},
	}
	websocket.DefaultDialer.HandshakeTimeout = time.Second * 5
	// websocket.DefaultDialer.Subprotocols = []string{"binary"}

Loop:

	log.Info("connecting: %s", u.String())
	conn, _, err = websocket.DefaultDialer.Dial(u.String(), wsHeader)
	if err != nil {
		log.Error("websocket.Dial(\"%s\") error(%v)", websocketAddr, err)
		log.Warn("websocket.Dial(\"%s\") reconnect after 1 seconds", websocketAddr)
		time.Sleep(1e9)
		goto Loop
	}

	log.Info("connect ok: %s", u.String())

	if heartbeat {
		go writeHeartbeat(conn)
	}

	for {
		op, m, err := conn.ReadMessage()
		if err != nil {
			log.Error("conn.ReadMessage() err(%v)", err)
			log.Warn("websocket.Dial(\"%s\") reconnect after 1 seconds", websocketAddr)
			time.Sleep(1e9)
			goto Loop
		}

		log.Debug("recv ws, op:%d, size:%d, err:%v, msg:%s", op, len(m), err, m)
		if op == int(OP_HEARTBEAT_REPLY) {
			log.Debug("----------------recieve heartbeat")
		}

		mType := gjson.Get(string(m), "type").String()
		if mType == "chat" {
			atomic.AddInt64(&countDown, 1)
		}
	}
}

func writeHeartbeat(conn *websocket.Conn) {
	proto1 := new(Proto)
	seqID := int32(0)
	for {
		// heartbeat
		proto1.Operation = OP_HEARTBEAT
		proto1.SeqId = seqID
		proto1.Body = nil
		if err := websocketWriteProto(conn, proto1); err != nil {
			log.Error("write heartbeat error(%v)", err)
			return
		} else {
			log.Debug("----------------send heartbeat %v", seqID)
		}
		seqID++

		time.Sleep(50000 * time.Millisecond)
	}
}
