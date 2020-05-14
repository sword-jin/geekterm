package geekhub

import (
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

	pages *tview.Pages

	replyFlex         *tview.Flex
	replyForm         *tview.Form
	replyContentField *tview.InputField
	errorModal        *tview.Modal

	curPostsPage   = 1
	curOffset      = 0 //当前浏览的区域（分子，拼车）
	curPost        *DetailPost
	curPreviewPost *PreviewPost
	curAuth        *AuthInfo
	curComment     *Comment
	curComments    []*Comment

	replyToken string
	replyTo    int8
)

const (
	replyToPost int8 = iota
	replyToComment
)

func Draw(app *tview.Application) {
	siderbar = tview.NewFlex().SetDirection(tview.FlexRow)
	siderbar.SetBorder(true).SetTitle(" 目录 ")
	siderbar.SetBorderPadding(0, 0, 1, 1)

	category = tview.NewList().ShowSecondaryText(false)
	category.SetBorder(false)
	category.SetHighlightFullLine(true)
	category.SetSelectedFocusOnly(true)
	category.SetSelectedBackgroundColor(tcell.ColorLightBlue)

	authStatusView = tview.NewTextView()
	authStatusView.SetBorder(true)
	authStatusView.SetTitle(" 用户 ")
	timer := time.NewTicker(DefaultAuthRefreshIntervel)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				//todo
			}
		}()
		for {
			<-timer.C
			Debugf("Refresh Auth Information.")
			response, err := GeekHub.GetHomePage(1)
			if err != nil {
				// todo
			} else {
				setAuthInformation(response.AuthInfo)
			}
		}
	}()

	siderbar.AddItem(category, 0, 4, true)
	siderbar.AddItem(authStatusView, 0, 1, true)

	posts = tview.NewList().ShowSecondaryText(true)
	posts.SetSecondaryTextColor(tcell.Color102)
	posts.SetBorder(true).SetTitle(" 列表 ")
	posts.SetHighlightFullLine(true)
	posts.SetSelectedBackgroundColor(tcell.ColorLightBlue)
	posts.SetSelectedFocusOnly(true)
	posts.SetBorderPadding(0, 0, 1, 1)

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
	commentList.SetBorder(true).SetTitle("  留言  ")
	commentList.SetBorderPadding(0, 0, 1, 0)
	commentList.SetSecondaryTextColor(tcell.Color102)
	commentList.SetChangedFunc(func(i int, _ string, _ string, _ rune) {
		curComment = curComments[i]
	})

	contentFlex.AddItem(contentView, 0, 5, true)
	contentFlex.AddItem(commentList, 0, 5, false)

	activityList = tview.NewList()
	activityList.SetSecondaryTextColor(tcell.Color102)
	activityFrame = tview.NewFrame(activityList)
	activityFrame.SetBorder(true)
	activityFrame.SetBorderPadding(0, 0, 1, 1).SetTitle("  我的动态  ")

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

	initErrorModal()

	//布局
	mainFlex = tview.NewFlex()
	mainFlex.SetTitle("terminal for Geekhub.com.")
	mainFlex.AddItem(siderbar, 0, 1, true).
		AddItem(posts, 0, 5, false)

	for _, item := range categoryList {
		category.AddItem(item.name, "", 0, nil)
	}

	//左侧选择区域
	category.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
		app.SetFocus(posts)
		loadPosts(app, i, 1)
	})
	loadPosts(app, 0, 1)

	welcomeList = tview.NewList()
	welcomeList.SetBorder(true)
	welcomeList.SetBorderAttributes(tcell.AttrUnderline)
	welcomeList.SetBorderPadding(1, 1, 2, 2)
	welcomeList.SetHighlightFullLine(true)
	welcomeList.SetTitle(fmt.Sprintf("  Welcome to geekterm %s  ", version))
	welcomeList.AddItem("进入", "", 0, func() {
		pages.SwitchToPage("main")
		app.SetFocus(category)
	})
	welcomeList.AddItem("快捷键", "", 0, func() {
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

	pages = tview.NewPages().
		AddPage("welcome", welcomePage, true, true).
		AddPage("main", mainFlex, true, false).
		AddPage("activities", activityFrame, true, false).
		AddPage("replyForm", replyFlex, true, false).
		AddPage("errorModal", errorModal, true, false)
	app.SetRoot(pages, true)
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
		doLoadPost(post.Uri)
	})
}

func doLoadPost(uri string) {
	if !contentShowing {
		mainFlex.AddItem(contentFlex, 0, 5, false)
		contentShowing = true
	}
	postResponse := doRequestPost(uri)

	commentList.Clear()
	curComment = nil
	comments := reverseComments(postResponse.Post.Comments)
	if len(comments) > 0 {
		curComment = comments[0]
		curComments = comments
	}
	for _, comment := range comments {
		commentList.AddItem(comment.Floor+" "+comment.Content, fmt.Sprintf("「%s」%s", comment.Author.Username, comment.CommentTime), 0, nil)
	}
}

// 这里把一些contentView的操作也做了
func doRequestPost(uri string) *ContentPageResponse {
	postResponse, err := GeekHub.GetPostContent(uri)
	if err != nil {
		// todo 错误处理
	}

	setAuthInformation(postResponse.AuthInfo)
	curPost = postResponse.Post
	contentView.SetTitle(fmt.Sprintf("  内容(%s)  ", curPost.PV))
	contentView.Clear()
	contentView.ScrollToBeginning()
	contentView.Write([]byte(`  标题：` + curPost.Title + "\n\n"))
	contentView.Write([]byte(postResponse.Post.Content))
	return postResponse
}

func getPostSecondaryText(post *PreviewPost) string {
	return fmt.Sprintf("评论: %d,「%s」发布,「%s」回复于%s", post.CommentCount, post.Author.Username, post.LatestReplyUser.Username, post.LatestReplyTime)
}

func setAuthInformation(authInfo *AuthInfo) {
	if authInfo == nil {
		authStatusView.Clear()
		authStatusView.SetTitleAlign(tview.AlignCenter)
		authStatusView.SetTextColor(tcell.ColorOrangeRed)
		authStatusView.SetBorderPadding(0, 0, 0, 0)
		authStatusView.Write([]byte(" 未登录 "))
	} else {
		setLoganAuthInfo(authInfo)
	}
}

func setLoganAuthInfo(authInfo *AuthInfo) {
	curAuth = authInfo

	authStatusView.Clear()
	authStatusView.SetBorderColor(tcell.ColorGreen)
	authStatusView.SetBorderPadding(0, 0, 1, 0)
	authStatusView.Write([]byte(authInfo.Me.Username + "\n"))
	authStatusView.Write([]byte("⏰: " + authInfo.NotifyCount + " 未读\n"))
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
					doLoadPost(activity.TargetUri)
				})
			} else if activity.Type == GbitOrder {
				activityList.AddItem(fmt.Sprintf("%s %s", activity.Time, activity.Content), "福利订单", 0, func() {
					OpenChrome(NewOpenableUrl(HomePage + GbitOrderURI))
				})
			} else if activity.Type == Unknow {
				activityList.AddItem("未适配，按 enter 进入 bug 提交页面", "福利订单", 0, func() {
					// todo 打开 github issue 页面
					//OpenChrome(NewOpenableUrl(HomePage + GbitOrderURI))
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
		doRequestPost(curPreviewPost.Uri)
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
	doLoadPost(curPost.Uri)
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
