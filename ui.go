package geekhub

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Category struct {
	name           string
	getPostHandler func(page int) (*PostPageResponse, error)
}

var categoryList []Category

var (
	mainFlex       *tview.Flex
	siderbar       *tview.Flex
	posts          *tview.List
	category       *tview.List
	contentFlex    *tview.Flex
	contentShowing bool
	authStatusView *tview.TextView

	contentView *tview.TextView
	commentList *tview.List

	activityList  *tview.List
	activityFrame *tview.Frame

	welcomePage *tview.Flex
	welcomeList *tview.List

	shutcutTable    *tview.Table
	shutcutPageFlex *tview.Flex
	shutcutPage     *tview.Frame

	newVersionContent *tview.TextView
	newVersionFlex    *tview.Flex
	newVersionPage    *tview.Frame

	pages *tview.Pages
	app   *tview.Application

	replyFlex         *tview.Flex
	replyForm         *tview.Form
	replyContentField *tview.InputField
	errorModal        *tview.Modal

	curPostsPage        = 1
	curOffset           = 0 //当前浏览的区域（分子，拼车）
	curPost             *DetailPost
	curPreviewPost      *PreviewPost
	curAuth             *AuthInfo
	curAuthUserDetail   *UserDetail
	curCommentPage      int
	curCommentTotalPage int64
	curComment          *Comment
	curComments         []*Comment
	globalHasCheckIn    bool //已签到

	replyToken string
	replyTo    int8

	setAuthChan    = make(chan struct{}, 1)
	enterGeekhubCh = make(chan struct{}, 1)
)

const (
	replyToPost int8 = iota
	replyToComment
)

func Run() error {
	return app.Run()
}

func Draw() {
	app = tview.NewApplication()

	initSetAuthStatusChan()

	initCategory(app)

	initAuthStatus(app)

	initSiderbar()

	initPosts()

	initPostContent()

	initActivity()

	initReply(app)

	initErrorModal()

	loadPosts(app, 0, 1)

	initWelcomePage(app)

	initShutcutPage(app)

	initNewVersion()

	//布局
	mainFlex = tview.NewFlex()
	mainFlex.SetTitle("terminal for Geekhub.com.")
	mainFlex.AddItem(siderbar, 0, 2, true).
		AddItem(posts, 0, 7, false)

	pages = tview.NewPages().
		AddPage("welcome", welcomePage, true, true).
		AddPage("shutcut", shutcutPage, true, false).
		AddPage("main", mainFlex, true, false).
		AddPage("activities", activityFrame, true, false).
		AddPage("replyForm", replyFlex, true, false).
		AddPage("errorModal", errorModal, true, false).
		AddPage("new-version", newVersionPage, true, false)
	app.SetRoot(pages, true)
}

func initSetAuthStatusChan() {
	go func() {
		for {
			<-setAuthChan

			if curAuth == nil {
				authStatusView.Clear()
				authStatusView.SetTitleAlign(tview.AlignCenter)
				authStatusView.SetTextColor(tcell.ColorOrangeRed)
				authStatusView.SetBorderPadding(0, 0, 0, 0)
				authStatusView.Write([]byte(" 未登录 "))
			} else {
				setLoganAuthInfo(curAuth)
			}
		}
	}()
}

func initSiderbar() {
	siderbar = tview.NewFlex().SetDirection(tview.FlexRow)
	siderbar.SetBorder(true).SetTitle(" 目录 ")
	siderbar.SetBorderPadding(0, 0, 1, 1)
	siderbar.AddItem(category, 0, 3, true)
	siderbar.AddItem(authStatusView, 0, 1, true)
}

func initAuthStatus(app *tview.Application) {
	authStatusView = tview.NewTextView()
	authStatusView.SetBorder(true)
	authStatusView.SetTitle(" 用户 ")
	timer := time.NewTicker(DefaultAuthRefreshIntervel)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				//todo
				Warnf("err: %v", err)
			}
		}()
		<-enterGeekhubCh
		for {
			if curAuth != nil {
				response, err := GeekHub.GetMePage(curAuth.Me.PageUri)
				if err != nil {
					// todo
				} else {
					setAuthInformation(response.AuthInfo)
					if response.AuthInfo.Me.UserDetail != nil {
						curAuthUserDetail = response.AuthInfo.Me.UserDetail
					}
				}
			}
			<-timer.C
		}
	}()
}

func initCategory(app *tview.Application) {
	category = tview.NewList().ShowSecondaryText(false)
	category.SetBorder(false)
	category.SetHighlightFullLine(true)
	category.SetSelectedFocusOnly(true)
	category.SetSelectedBackgroundColor(tcell.ColorLightBlue)
	for _, item := range categoryList {
		category.AddItem(item.name, "", 0, nil)
	}

	//左侧选择区域
	category.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
		app.SetFocus(posts)
		loadPosts(app, i, 1)
	})
}

func initNewVersion() {
	newVersionContent = tview.NewTextView()
	newVersionContent.SetBorder(true).SetTitleColor(tcell.ColorGreen)
	newVersionContent.SetBorderPadding(0, 0, 1, 1)

	newVersionFlex = tview.NewFlex()
	newVersionFlex.AddItem(tview.NewBox(), 0, 1, false)
	newVersionFlex.AddItem(tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 2, false).
		AddItem(newVersionContent, 0, 3, true).
		AddItem(tview.NewBox(), 0, 2, false), 0, 1, true)
	newVersionFlex.AddItem(tview.NewBox(), 0, 1, false)

	newVersionPage = tview.NewFrame(newVersionFlex)
	newVersionPage.SetBorder(true)
}

func initReply(app *tview.Application) {
	replyForm = tview.NewForm()
	replyForm.SetBorder(true)
	replyForm.SetBorderPadding(1, 1, 1, 1)
	replyForm.AddInputField("内容", "", 0, nil, nil).
		AddButton("提交", func() {
			submitReplyForm(app)
		}).
		SetTitleAlign(tview.AlignLeft)
	replyForm.AddButton("取消", func() {
		cancelReply(app)
	})
	replyContentField = replyForm.GetFormItem(0).(*tview.InputField)

	replyFlex = tview.NewFlex()
	replyFlex.SetBorder(false)
	replyFlex.AddItem(tview.NewBox(), 0, 1, false)
	replyFlex.AddItem(tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 2, false).
		AddItem(replyForm, 0, 3, true).
		AddItem(tview.NewBox(), 0, 2, false), 0, 1, true)
	replyFlex.AddItem(tview.NewBox(), 0, 1, false)
}

func initActivity() {
	activityList = tview.NewList()
	activityList.SetSecondaryTextColor(tcell.Color102)
	activityFrame = tview.NewFrame(activityList)
	activityFrame.SetBorder(true)
	activityFrame.SetBorderPadding(0, 0, 1, 1).SetTitle("  我的动态  ")
}

func initPostContent() {
	contentFlex = tview.NewFlex()
	contentFlex.SetBorder(true)
	contentFlex.SetDirection(tview.FlexRow)
	contentShowing = false

	contentView = tview.NewTextView()
	contentView.SetTitle("  内容  ")
	contentView.SetBorder(true)
	contentView.SetScrollable(true)
	contentView.SetBorderPadding(0, 0, 1, 1)

	commentList = tview.NewList()
	commentList.SetSelectedFocusOnly(true)
	commentList.SetBorder(true).SetTitle("  评论  ")
	commentList.SetBorderPadding(0, 0, 1, 0)
	commentList.SetSecondaryTextColor(tcell.Color102)
	commentList.SetChangedFunc(func(i int, _ string, _ string, _ rune) {
		curComment = curComments[i]
	})

	contentFlex.AddItem(contentView, 0, 5, true)
	contentFlex.AddItem(commentList, 0, 5, false)
}

func initPosts() {
	posts = tview.NewList().ShowSecondaryText(true)
	posts.SetSecondaryTextColor(tcell.Color102)
	posts.SetBorder(true).SetTitle(" 列表 ")
	posts.SetHighlightFullLine(true)
	posts.SetSelectedBackgroundColor(tcell.ColorLightBlue)
	posts.SetSelectedFocusOnly(true)
	posts.SetBorderPadding(0, 0, 1, 1)
}

func initWelcomePage(app *tview.Application) {
	welcomeList = tview.NewList()
	welcomeList.SetBorder(true)
	welcomeList.SetBorderAttributes(tcell.AttrUnderline)
	welcomeList.SetBorderPadding(1, 1, 2, 2)
	welcomeList.SetHighlightFullLine(true)
	welcomeList.SetTitle(fmt.Sprintf("  Welcome to geekterm %s  ", Version))
	welcomeList.SetTitleColor(tcell.ColorPink)
	welcomeList.AddItem("进入", "", 0, func() {
		enterGeekhub(app)
	})
	welcomeList.AddItem("快捷键", "", 0, func() {
		pages.SwitchToPage("shutcut")
		app.SetFocus(shutcutTable)
	})
	welcomeList.AddItem("🐞提交BUG", "", 0, func() {
		OpenChrome(NewOpenableUrl("https://github.com/rrylee/geekterm"))
	})
	welcomeList.AddItem("⭐️Github", "", 0, func() {
		OpenChrome(NewOpenableUrl("https://github.com/rrylee/geekterm/issues"))
	})

	welcomePage = tview.NewFlex()
	welcomePage.SetBorder(false)
	welcomePage.AddItem(tview.NewBox(), 0, 1, false)
	welcomePage.AddItem(tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 2, false).
		AddItem(welcomeList, 0, 3, true).
		AddItem(tview.NewBox(), 0, 2, false), 0, 1, true)
	welcomePage.AddItem(tview.NewBox(), 0, 1, false)
}

func enterGeekhub(app *tview.Application) {
	pages.SwitchToPage("main")
	app.SetFocus(category)

	if curAuth != nil {
		// 自动签到
		hasCheck, signToken, err := GeekHub.GetSignStatus()
		if err != nil {
			showErrorModal(app, "获取签到信息失败")
			return
		} else {
			if !hasCheck {
				err = GeekHub.CheckIn(signToken)
				if err != nil {
					globalHasCheckIn = true
				}
			} else {
				globalHasCheckIn = true
			}
		}
		enterGeekhubCh <- struct{}{}
	}
}

func initShutcutPage(app *tview.Application) {
	shutcutTable = tview.NewTable()
	shutcutTable.SetBorder(true)
	shutcutTable.SetBorderPadding(0, 0, 2, 2)
	shutcutTable.SetTitle("  快捷键列表  ")

	shutcutPageFlex = tview.NewFlex()
	shutcutPageFlex.AddItem(tview.NewBox(), 0, 1, false)
	shutcutPageFlex.AddItem(tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 2, false).
		AddItem(shutcutTable, 0, 3, true).
		AddItem(tview.NewBox(), 0, 2, false), 0, 1, true)
	shutcutPageFlex.AddItem(tview.NewBox(), 0, 1, false)

	shutcutPage = tview.NewFrame(shutcutPageFlex)

	for i, name := range ShutcutTitles {
		shutcutTable.SetCell(0, i, &tview.TableCell{Text: name, Align: tview.AlignCenter, Color: tcell.ColorDarkMagenta})
	}

	for _, key := range ShutcutKeys {
		row := shutcutTable.GetRowCount()
		shutcutTable.SetCell(row, 0, &tview.TableCell{Text: key[0], Align: tview.AlignCenter, Color: tcell.ColorWhite})
		shutcutTable.SetCell(row, 1, &tview.TableCell{Text: key[1], Align: tview.AlignCenter, Color: tcell.ColorGrey})
		shutcutTable.SetCell(row, 2, &tview.TableCell{Text: key[2], Align: tview.AlignCenter, Color: tcell.ColorGrey})
	}

	shutcutTable.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEscape:
			pages.SwitchToPage("welcome")
			app.SetFocus(welcomeList)
		}
	})
}

func initErrorModal() {
	errorModal = tview.NewModal()
	errorModal.SetTitle("  出现错误  ")
	errorModal.SetBorder(true)
	errorModal.SetTitleAlign(tview.AlignLeft)
	errorModal.SetBorderPadding(1, 1, 1, 1)
}

var (
	firstLoadPosts = true
)

func loadPosts(app *tview.Application, offset int, page int) {
	response, err := categoryList[offset].getPostHandler(page)
	if err != nil {
		// todo
	}
	posts.Clear()

	if firstLoadPosts {
		// todo，留着这个代码，说不定有用
		firstLoadPosts = false
		// 第一次加载的时候，添加一个 previewPost
		if len(response.Posts) > 0 {
			curPreviewPost = response.Posts[0]
		}
	}

	setAuthInformation(response.AuthInfo)

	for _, post := range response.Posts {
		loadPost(app, post)
	}
	curOffset = offset
	curPostsPage = page
	posts.SetTitle(fmt.Sprintf("  第%d页  ", curPostsPage))

	posts.SetChangedFunc(func(i int, _ string, _ string, _ rune) {
		curPreviewPost = response.Posts[i]
	})
}

func loadPost(app *tview.Application, post *PreviewPost) *tview.List {
	return posts.AddItem(post.Title, getPostSecondaryText(post), 0, func() {
		doLoadPost(post.Uri, -1)
	})
}

func doLoadPost(uri string, page int) {
	if !contentShowing {
		mainFlex.AddItem(contentFlex, 0, 6, false)
		contentShowing = true
	}
	postResponse := doRequestPost(uri, page)

	commentList.Clear()
	curComment = nil
	curCommentPage = postResponse.Post.CurCommentPage
	curCommentTotalPage = postResponse.Post.CommentTotalPage
	comments := reverseComments(postResponse.Post.Comments)
	if len(comments) > 0 {
		curComment = comments[0]
		curComments = comments
		commentList.SetTitle(fmt.Sprintf("  评论第%d页  ", curCommentPage))
	}
	for _, comment := range comments {
		second := fmt.Sprintf("「%s」%s", comment.Author.Username, comment.CommentTime)
		if comment.Parent != nil {
			second += comment.Parent.Author.Username + ": " + comment.Parent.Content
		}
		commentList.AddItem(comment.Floor+" "+comment.Content, second, 0, nil)
	}

	Debugf("curCommentPage is %d", curCommentPage)
}

// 这里把一些contentView的操作也做了
func doRequestPost(uri string, page int) *ContentPageResponse {
	postResponse, err := GeekHub.GetPostContent(uri, page)
	if err != nil {
		// todo 错误处理
	}

	setAuthInformation(postResponse.AuthInfo)
	curPost = postResponse.Post
	contentView.SetTitle(fmt.Sprintf("  内容(%s)  ", curPost.PV))
	contentView.Clear()
	contentView.ScrollToBeginning()
	content := bytes.NewBufferString(`  标题：` + curPost.Title + "\n")
	content.Write([]byte("  楼主：" + curPost.Author.Username + "发布于" + curPost.PublishTime + "\n\n"))

	if curPost.PostType == MoleculeType {
		moleculeInfo, ok := curPost.ExtraInfo.(*MoleculesInfo)
		if ok {
			content.Write([]byte(`  分子贴：` + moleculeInfo.Name + "\n"))
			if moleculeInfo.Floor == "" {
				content.Write([]byte(`  ` + formatSeconds(moleculeInfo.CountDown) + "\n"))
			}
			content.Write([]byte(`  价值：` + moleculeInfo.Price + "\n"))
			content.Write([]byte(`  中奖比例：` + moleculeInfo.Molecule + "/" + moleculeInfo.Denominator + "\n"))
			content.Write([]byte(`  物流：` + moleculeInfo.HowToSend + "\n"))
			content.Write([]byte(`  联系方式：` + moleculeInfo.Contact + "\n"))

			if moleculeInfo.Floor != "" {
				content.Write([]byte(`  分子楼层：` + moleculeInfo.Floor + "\n\n"))
			} else {
				content.Write([]byte("\n"))
			}
		}
	}

	content.Write([]byte(postResponse.Post.Content))
	contentView.Write(content.Bytes())
	return postResponse
}

func getPostSecondaryText(post *PreviewPost) string {
	return fmt.Sprintf("评论: %d,「%s」发布,「%s」回复于%s", post.CommentCount, post.Author.Username, post.LatestReplyUser.Username, post.LatestReplyTime)
}

func setAuthInformation(authInfo *AuthInfo) {
	curAuth = authInfo
	setAuthChan <- struct{}{}
}

func setLoganAuthInfo(authInfo *AuthInfo) {
	authStatusView.Clear()
	authStatusView.SetBorderColor(tcell.ColorGreen)
	authStatusView.SetBorderPadding(0, 0, 1, 0)

	content := bytes.NewBufferString(authInfo.Me.Username + "\n")
	content.Write([]byte("⏰: " + authInfo.NotifyCount + " 未读\n"))
	if curAuthUserDetail != nil {
		content.Write([]byte("Gbit: " + curAuthUserDetail.Gbit + "\n"))
		content.Write([]byte("Star: " + curAuthUserDetail.Star + "\n"))
		content.Write([]byte("Score: " + curAuthUserDetail.Score + "\n"))
	}

	if globalHasCheckIn {
		content.Write([]byte("已签到\n"))
	}

	Debugf(content.String())
	authStatusView.Write(content.Bytes())
	app.Draw()
}

func showActivities(app *tview.Application) {
	pages.SwitchToPage("activities")

	response, err := GeekHub.GetActivities(1)
	if err != nil {
		showErrorModal(app, err.Error())
		return
	}

	activityList.Clear()
	app.SetFocus(activityList)
	setAuthInformation(response.AuthInfo)
	for _, activity := range response.Activities {
		func(activity *Activity) {
			if activity.Type == ReplyPost || activity.Type == GetMolecules || activity.Type == YourMoleculesFinish {
				var title, content string
				if activity.Type == ReplyPost {
					content = fmt.Sprintf("%s「%s」在「%s」回复", activity.Time, activity.User.Username, activity.TargetTitle)
					title = activity.Content
				} else if activity.Type == GetMolecules {
					title = fmt.Sprintf("%s %s", activity.Time, activity.TargetTitle)
					content = "抢到分子"
				} else {
					title = fmt.Sprintf("%s %s", activity.Time, activity.TargetTitle)
					content = "你的分子出现了"
				}
				activityList.AddItem(title, content, 0, func() {
					pages.SwitchToPage("main")
					app.SetFocus(commentList)
					doLoadPost(activity.TargetUri, -1)
				})
			} else if activity.Type == GbitOrder {
				activityList.AddItem(fmt.Sprintf("%s %s", activity.Time, activity.Content), "福利订单", 0, func() {
					OpenChrome(NewOpenableUrl(HomePage + GbitOrderURI))
				})
			} else if activity.Type == OrderComment {
				activityList.AddItem(fmt.Sprintf("%s", activity.Content), "订单评价", 0, nil)
			} else if activity.Type == Unknow {
				activityList.AddItem("🐞 🐞 🐞 🐞 🐞 🐞 🐞 🐞 🐞", "出现bug了，来github提交bug吧", 0, func() {
					OpenChrome(NewOpenableUrl("https://github.com/rrylee/geekterm/issues"))
				})
			}
		}(activity)
	}
}

func showErrorModal(app *tview.Application, msg string) {
	errorModal.SetText(" 错误: " + msg)
	pages.SwitchToPage("errorModal")
	app.SetFocus(errorModal)
}

func replyPost(app *tview.Application) {
	replyTo = replyToPost
	replyForm.SetTitle(fmt.Sprintf("-  回复「%s」 ", getReplyPostTitle()))
	pages.SwitchToPage("replyForm")
	app.SetFocus(replyContentField)
}

func replyComment(app *tview.Application) {
	replyTo = replyToComment
	replyForm.SetTitle(fmt.Sprintf("  回复 %s @%s  ", curComment.Floor, curComment.Author.Username))
	pages.SwitchToPage("replyForm")
	app.SetFocus(replyContentField)
}

func getReplyPostTitle() string {
	if curPost != nil {
		if len(curPost.Title) > 24 {
			return curPost.Title[0:24] + "..."
		} else {
			return curPost.Title
		}
	} else {
		return curPreviewPost.Title
	}
}

func submitReplyForm(app *tview.Application) {
	if curPost == nil {
		//load cur post
		doRequestPost(curPreviewPost.Uri, -1)
	}

	replyArg := &PostCommentArgs{
		AuthenticityToken: replyToken,
		TargetType:        curPost.PostType,
		TargetId:          curPost.ID,
		ReplyToId:         getReplyToId(),
		Content:           replyContentField.GetText() + DefaultSign,
	}

	Debugf("Submit reply form, arg:%v", replyArg)

	err := GeekHub.PostComment(replyArg)
	if err != nil {
		Warnf("PostCommentURI error:%v", err)
		showErrorModal(app, fmt.Sprintf("评论失败, %s", err.Error()))
		replyToken = ""
		return
	}

	replyContentField.SetText("")
	replyToken = ""
	pages.SwitchToPage("main")
	doLoadPost(curPost.Uri, -1)
	app.SetFocus(commentList)
}

func getReplyToId() string {
	if replyTo == replyToComment {
		return curComment.ID
	} else {
		return "0"
	}
}

func reverseComments(comments []*Comment) []*Comment {
	l := len(comments)
	newComments := make([]*Comment, l)
	for i, c := range comments {
		newComments[l-1-i] = c
	}
	return newComments
}
