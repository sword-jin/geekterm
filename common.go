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

var ShutcutTitles = []string{"键位", "功能", "备注"}

var ShutcutKeys = [][3]string{
	{"上下左右", "正常移动", "左右可以切换不同的窗口"},
	{"j k h l", "对应以上", ""},
	{"i", "查看个人动态", ""},
	{"o", "打开到浏览器", "选中帖子，打开帖子，在动态，打开动态"},
	{"n", "上一页", "帖子列表，评论列表"},
	{"m", "下一页", "帖子列表，评论列表"},
	{"r", "回帖", "如果选中某条评论，会回复留言"},
	{"R", "回帖", "直接回帖"},
	{"M", "直接进入评论列表", ""},
	{"enter", "各种确定操作", "加载帖子，提交评论"},
	{"esc", "退出", "个人动态页使用"},
	{"tab", "切换", "在评论弹窗切换选中区"},
}
