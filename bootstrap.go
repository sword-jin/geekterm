package geekhub

import (
	"os"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

func Setup(cfg *Config) {
	MyConfig = cfg

	logger = logrus.New()
	f, _ := os.OpenFile(MyConfig.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	logger.Out = f
	logger.SetLevel(logrus.Level(MyConfig.LogLevel))

	httpClient = resty.New()
	httpClient.SetHeader("User-Agent", DefaultUserAgent)
	httpClient.SetHeader("Cookie", cfg.Cookie)
	httpClient.SetContentLength(true)
	GeekHub = &geekHub{
		Selectors: GetDefaultConfigAttrSelectors(),
	}

	Converter = &converter{
		engine: md.NewConverter("", true, nil),
	}

	categoryList = []Category{
		{
			"  首页  ",
			GeekHub.GetHomePage,
		},
		{
			"  话题区  ",
			GeekHub.GetPostsPage,
		},
		{
			"  二手  ",
			GeekHub.GetSecondHandsPage,
		},
		{
			"  拍卖  ",
			GeekHub.GetAuctionsPage,
		},
		{
			"  分子  ",
			GeekHub.GetMoleculesPage,
		},
		{
			"  拼车  ",
			GeekHub.GetGroupBuysPage,
		},
	}
}
