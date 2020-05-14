package geekhub

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func Keybinds(app *tview.Application) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		Debugf(spew.Sdump(event))

		curFocus := app.GetFocus()
		if curFocus == replyContentField { //编辑模式下，禁用快捷键
			//if event.Key() == tcell.KeyESC {
			//	return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModMask(0))
			//} else {
			//}
			return event
		}

		switch event.Rune() {
		case 'j':
			if app.GetFocus() == contentView {
				r, _ := contentView.GetScrollOffset()
				if !curPost.clickedDown {
					curPost.clickedDown = true
				} else {
					if r == curPost.lastScrollRow {
						app.SetFocus(commentList)
					}
				}
				curPost.lastScrollRow = r
			}

			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModMask(0))
		case 'q':
			app.Stop()
		case 'k':
			if app.GetFocus() == contentView {
				r, _ := contentView.GetScrollOffset()
				if curPost.clickedDown && r == 0 {
					curPost.clickedDown = false
				}
			}
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModMask(0))
		case 'h': //<-
			handleLeft(app)
		case 'l': //->
			handleRight(app)
		case 'm': //next page
			if app.GetFocus() == posts {
				Debugf("start load next page")
				loadPosts(app, curOffset, curPostsPage+1)
			}
		case 'n': //previous page
			if curFocus == posts && curPostsPage > 1 {
				loadPosts(app, curOffset, curPostsPage-1)
			}
		case 'M':
			if curFocus == posts || curFocus == contentView {
				app.SetFocus(commentList)
			}
		case 'o':
			if curFocus == contentView || curFocus == posts || curFocus == commentList {
				// todo，错误处理
				if curPost != nil {
					OpenChrome(curPost)
				} else {
					OpenChrome(curPreviewPost)
				}
			}
		case 'i':
			//打开个人未读
			if curAuth == nil {
				showErrorModal(app, "请先登录")
				return event
			}
			showActivities(app)
		case 'r':
			Debugf("curAuth:%v", curAuth)
			if curAuth == nil {
				showErrorModal(app, "请先登录")
				return event
			}
			if curFocus == contentView || curFocus == posts {
				replyPost(app)
			} else if curFocus == commentList {
				replyComment(app)
			}
		case 'R':
			if curFocus == contentView || curFocus == posts || curFocus == commentList {
				replyPost(app)
			}
		}

		switch event.Key() {
		case 260: // <-键
			handleLeft(app)
		case 259: // ->键
			handleRight(app)
		case tcell.KeyESC: // ESC
			pages.SwitchToPage("main")
		}

		return event
	})
}

func cancelReply(app *tview.Application) {
	pages.SwitchToPage("main")
	replyContentField.SetText("")
	app.SetFocus(posts)
}

func handleRight(app *tview.Application) {
	if app.GetFocus() == category {
		app.SetFocus(posts)
	} else if app.GetFocus() == posts {
		if contentShowing { //修复，在contentView没有显示的时候，造成的错误
			app.SetFocus(contentView)
			contentView.SetBorderColor(tcell.ColorGreen)
		}
	}
}

func handleLeft(app *tview.Application) {
	if app.GetFocus() == posts {
		app.SetFocus(category)
	} else if app.GetFocus() == contentView || app.GetFocus() == commentList {
		app.SetFocus(posts)
		contentView.SetBorderColor(tcell.ColorWhite)
	}
}
