package geekhub

var MyConfig *Config

type Config struct {
	Cookie    string //登录cookie
	ReplySign string //回复签名
	LogFile   string
	LogLevel  int
}
