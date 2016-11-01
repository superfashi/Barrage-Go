package configs

import "time"

// The personal data of danmaku hime.

const (
	RETRY_TIME = 5 * time.Second
	MAX_RETRY = 5
	TIME_OUT = 20 * time.Second

	TIME_FORMAT = "15:04:05"

	SEND_FORMAT = map[string]string{
		"color": "000000",
		"fontsize": "11",
		"mode": "1",
	}
)