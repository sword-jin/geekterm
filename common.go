package geekhub

import (
	"errors"
	"time"
)

const DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.82 Safari/537.36"
const DefaultSign = "\n\n 「来自 geekterm」"
const DefaultAuthRefreshIntervel = 60 * time.Second

var (
	InternetError = errors.New("Internet Error.")
	GoQueryError  = errors.New("GoQuery error.")
)
