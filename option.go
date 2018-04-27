package alexander

import "time"

type Option struct {
	HTTPVersion     string
	ConnectTimeout  time.Duration
	ReadTimeout     time.Duration
	IdleConnTimeout time.Duration
	KeepAlive       time.Duration
	MaxIdleConns    int
}
