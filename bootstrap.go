package geekhub

import (
	"fmt"
	"os"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

func Setup(cfg *Config) {
	MyConfig = cfg

	logger = logrus.New()
	f, _ := os.OpenFile(MyConfig.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	logger.Out = f
	logger.SetLevel(logrus.Level(MyConfig.LogLevel))

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

func WatchUpgrade(app *tview.Application) {
	go func() {
		firstTimer := time.NewTimer(3 * time.Second)
		<-firstTimer.C
		handleVersionCheck(app)

		timer := time.NewTicker(1 * time.Hour)

		for range timer.C {
			handleVersionCheck(app)
		}
	}()
}

func handleVersionCheck(app *tview.Application) {
	hasNewVersion, newVersion, err := CheckNewVersion()
	Infof("WatchUpgrade %v-%v-%v", hasNewVersion, newVersion, err)
	if err != nil {
		Warnf("WatchUpgrade error:%v", err)
	} else {
		if hasNewVersion {
			newVersionContent.SetTitle(fmt.Sprintf("新版本 %s 发布", newVersion.t))
			newVersionContent.Clear()
			newVersionContent.Write([]byte(newVersion.s))

			pages.SwitchToPage("new-version")
			app.SetFocus(newVersionContent)
		}
	}
}
