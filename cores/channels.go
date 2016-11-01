package cores

import (
	"net"
	"time"
	"net/http"
	"net/http/cookiejar"
	"github.com/hanbang-wang/Barrage-Go/configs"
	"io/ioutil"
	"encoding/xml"
	"log"
	"errors"
	"fmt"
	"encoding/hex"
	"encoding/binary"
	"github.com/bitly/go-simplejson"
)

type Connector struct {
	server         string
	session        *http.Client
	iSocket        net.Conn
	isConnect      bool
	roomID, userID int
}

func NewConnector(roomid, userid int) (ret *Connector) {
	ret = &Connector{roomID: roomid}
	var err error
	if userid < 0 {
		ret.userID = RandInt(1e5, 4e7)
	} else {
		ret.userID = userid
	}
	if ret.server, err = ret.getServerLink(); err != nil {
		ret.server = configs.SERVER_URL
	}
	jar, _ := cookiejar.New(nil)
	ret.session = &http.Client{Jar: jar}
	ret.connect()
	go ret.connect()
	go ret.receive()
	return
}

func (s *Connector) init() bool {
	var err error
	if s.iSocket, err = net.DialTimeout("tcp", s.server + ":788", configs.TIME_OUT); err == nil {
		if isFirst {
			log.Println("开始接收弹幕。(Ctrl + C 退出)")
			isFirst = false
		}
		s.isConnect = true
		return true
	}
	return false
}

func (s *Connector) connect() {
	if s.isConnect {
		s.iSocket.Close()
	}
	s.isConnect = false
	retryTime := 0
	for !s.init() || !s.handshake() {
		if retryTime >= configs.MAX_RETRY {
			log.Fatal(errors.New("重试请求过多，服务中止！"))
		}
		log.Println("服务器连接失败……")
		time.Sleep(configs.RETRY_TIME)
		retryTime++
	}
}

func (s *Connector) getServerLink() (link string, err error) {
	ret := new(Root)
	var resp *http.Response
	if resp, err = Network(s.session, configs.PLAYER_API, "GET", Fm("id=cid:%d", s.roomID), Fm(configs.LIVE_ROOM, s.roomID)); err == nil {
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			body = append(append([]byte(`<Root>`), body...), []byte(`</Root>`)...)
			if err := xml.Unmarshal(body, &ret); err == nil {
				link = ret.Server
			}
		}
	}
	return
}

func (s *Connector) handshake() bool {
	cdata := fmt.Sprintf(`{"roomid":%d,"uid":%d}`, s.roomID, s.userID)
	handshake := Fm(configs.HANDSHAKE_STR, len(cdata) + 16)
	buf := make([]byte, len(handshake) >> 1)
	hex.Decode(buf, []byte(handshake))
	if _, err := s.iSocket.Write(append(buf, []byte(cdata)...)); err != nil {
		return false
	}
	for {
		buf := make([]byte, 16)
		hex.Decode(buf, configs.HEARTBEAT_BYTE)
		if _, err := s.iSocket.Write(buf); err != nil {
			go s.connect()
			return true
		}
		time.Sleep(30 * time.Second)
	} // should never end
}

func (s *Connector) receive() {
	for {
		if s.isConnect {
			buffer := make([]byte, 4)
			if _, err := s.iSocket.Read(buffer); err != nil {
				continue
			}
			repr := binary.BigEndian.Uint32(buffer)
			if _, err := s.iSocket.Read(buffer); err != nil {
				continue
			}
			if _, err := s.iSocket.Read(buffer); err != nil {
				continue
			}
			typ := binary.BigEndian.Uint32(buffer)
			if _, err := s.iSocket.Read(buffer); err != nil {
				continue
			}
			buffer = make([]byte, repr-16)
			if _, err := s.iSocket.Read(buffer); err != nil {
				continue
			}
		}
	}
}