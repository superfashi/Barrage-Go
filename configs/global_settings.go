package configs

// The global data of danmaku hime.

const (
	// RECEIVER PART
	SERVER_URL = "dm.live.bilibili.com"
	PLAYER_API = "http://live.bilibili.com/api/player"
	HANDSHAKE_STR = "%08x001000010000000700000001"
	HEARTBEAT_BYTE = []byte("00000010001000010000000200000001")
	HEARTBEAT_KEEP_TIME = 30
	// LOGIN PART
	LOGIN_URL = "https://passport.bilibili.com/ajax/miniLogin/minilogin"
	LOGIN_HEADER = map[string]string{
		"Host": "passport.bilibili.com",
		"Referer": "https://passport.bilibili.com/ajax/miniLogin/minilogin",
		"Content-Type": "application/x-www-form-urlencoded",
		"Connection": "keep-alive",
		"Cache-Control": "max-age=0",
		"Origin": "https://passport.bilibili.com",
	}
	// SENDING PART
	SEND_URL = "http://live.bilibili.com/msg/send"
	// REFERER
	LIVE_ROOM = "http://live.bilibili.com/%d"
)