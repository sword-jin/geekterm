package geekhub

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

var (
	httpClient = resty.New()
	GeekHub    IGeekHub
)

const (
	HomePage       = "https://geekhub.com/"
	PostsURI       = "posts"
	SecondHandsURI = "second_hands"
	AuctionsURI    = "auctions"
	MoleculesURI   = "molecules"
	GroupBuysURI   = "group_buys"
	ActivitiesURI  = "activities"
	PostCommentURI = "comments"
	GbitOrderURI   = "gbit_orders"
	CheckinsURI    = "checkins"
	SignURI        = "checkins/start"
)

type ConfigAttrSelectors struct {
	PostList                  string
	Post                      string
	PostHref                  string
	TextUser                  string //用户名（纯文字的）
	PostCommentCount          string
	AuthStatus                string //右上角登录状态
	PostBody                  string
	PostPageCommentCount      string //post页面的评论数量标签
	CommentList               string
	CommentContent            string
	PostPageAuthor            string
	PostPageTitle             string
	ActivityCount             string
	AuthUsername              string
	CommentAuthor             string
	CommentTime               string
	CommentFloor              string
	Activities                string
	ActivityTargetUri         string
	ActivityTargetTitle       string
	ActivityReplyContent      string
	ActivityBody              string //整个动态中间内容区域
	ActivityUser              string
	ActivityTime              string
	ActivityGetMoleculesLink  string
	ActivityYourMoleculesLink string
	ReplyToken                string
	CommentParentUser         string
	CommentParentContent      string
	CommentCurPage            string
	CheckInButton             string
	MoleculeSeconds           string
}

func GetDefaultConfigAttrSelectors() *ConfigAttrSelectors {
	return &ConfigAttrSelectors{
		PostList:                  "article",
		Post:                      "h3",
		PostHref:                  "a",
		TextUser:                  "a",
		PostCommentCount:          "a.badge.py-2px.sub",
		AuthStatus:                "header .flex.items-center.ml-5",
		PostBody:                  "main>.box .story",
		PostPageCommentCount:      "main .mt-5.box .flex.items-center.justify-between.p-3",
		CommentList:               ".comment-list",
		PostPageAuthor:            ".mr-2.text-sm.font-bold.text-primary-600 > a",
		PostPageTitle:             "main .heading",
		CommentContent:            ".break-all.max-h-screen.overflow-y-auto",
		ActivityCount:             "a.inline-flex.items-center.mr-5:nth-of-type(1)",
		AuthUsername:              "a.inline-flex.items-center.mr-5:nth-of-type(2)",
		CommentAuthor:             ".mr-2.font-semibold>a",
		CommentTime:               "div.inline-flex.items-center:nth-of-type(1)>span:nth-of-type(5)",
		CommentFloor:              ".inline-flex.items-center:nth-of-type(2)>span:nth-of-type(1)",
		CommentParentUser:         ".mt-2.text-primary-700 .meta",
		CommentParentContent:      ".mt-2.text-primary-700 .block",
		Activities:                "main .flex.items-center.px-3.py-4.border-t.border-color",
		ActivityTargetUri:         "div:nth-of-type(2) div:nth-of-type(1) a:nth-of-type(2)",
		ActivityTargetTitle:       "div:nth-of-type(2) div:nth-of-type(1) a:nth-of-type(2)",
		ActivityReplyContent:      "div:nth-of-type(2) div:nth-of-type(2) p",
		ActivityUser:              "div:nth-of-type(2) div:nth-of-type(1) a:nth-of-type(1)",
		ActivityBody:              ".flex-1",
		ActivityGetMoleculesLink:  ".flex-1 a",
		ActivityYourMoleculesLink: ".flex-1 a",
		ActivityTime:              "div:nth-of-type(3)",
		ReplyToken:                "form#comment-box-form>input:nth-of-type(1)",
		CommentCurPage:            "nav .px-2.py-px.rounded.bg-primary-300",
		CheckInButton:             ".btn.btn-blue.btn-big.block.text-center",
		MoleculeSeconds:           "var seconds = ",
	}
}

type UserDetail struct {
	Star  string
	Gbit  string
	Score string
}

type User struct {
	Username   string
	PageUri    string
	UserDetail *UserDetail
}

type PreviewPost struct {
	ID              string
	Title           string
	Uri             string
	Author          *User
	CommentCount    int64
	LatestReplyTime string
	LatestReplyUser *User
}

func (p PreviewPost) GetUrl() string {
	return HomePage + p.Uri
}

type AuthStatus int8

const (
	NotLogin AuthStatus = 0
	Logan    AuthStatus = 1
)

type AuthInfo struct {
	NotifyCount string //动态条数
	Me          *User  //用户信息
}

type BasePageInfo struct {
	AuthInfo *AuthInfo
	Status   AuthStatus
}

type PostPageResponse struct {
	*BasePageInfo
	Posts []*PreviewPost
}

type Comment struct {
	ID          string
	Author      *User
	CommentTime string
	Content     string
	Floor       string
	Parent      *Comment
}

type DetailPost struct {
	ID               string
	Title            string
	Content          string
	PostType         postType
	Uri              string
	Author           *User
	PublishTime      string
	PV               string
	CommentCount     int64
	CommentTotalPage int64
	CurCommentPage   int
	Comments         []*Comment
	ExtraInfo        interface{}

	//helper
	lastScrollRow int
	clickedDown   bool //按下
}

type MoleculesInfo struct {
	Name        string
	Price       string
	Molecule    string //分子
	Denominator string //分母
	HowToSend   string
	Contact     string
	CountDown   int //倒计时
	Floor       string
}

func (p DetailPost) GetUrl() string {
	return HomePage + p.Uri
}

type ContentPageResponse struct {
	*BasePageInfo
	Post *DetailPost
}

type ActivityType int8

const (
	ReplyPost = iota
	GbitOrder
	GetMolecules        //抢到分子
	YourMoleculesFinish //分子结束
	Unknow              //未识别的
)

type Activity struct {
	Type        ActivityType
	TargetUri   string
	TargetTitle string
	Content     string
	User        *User
	Time        string
}

func (a *Activity) GetUrl() string {
	return HomePage + a.TargetUri
}

type ActivitiesPageResponse struct {
	*BasePageInfo
	Activities []*Activity
}

type MePageResponse struct {
	*BasePageInfo
}

type IGeekHub interface {
	GetMePage(userUri string) (*MePageResponse, error) //个人中心
	GetHomePage(page int) (*PostPageResponse, error)
	GetPostsPage(page int) (*PostPageResponse, error)
	GetSecondHandsPage(page int) (*PostPageResponse, error)
	GetAuctionsPage(page int) (*PostPageResponse, error)
	GetMoleculesPage(page int) (*PostPageResponse, error)
	GetGroupBuysPage(page int) (*PostPageResponse, error)

	GetPostContent(pageUri string, page int) (*ContentPageResponse, error)
	GetActivities(page int) (*ActivitiesPageResponse, error)
	PostComment(arg *PostCommentArgs) error

	GetSignStatus() (bool, string, error) //获取签到状态
	CheckIn(token string) error           //签到
}

type geekHub struct {
	Selectors *ConfigAttrSelectors
}

func (gh *geekHub) GetHomePage(page int) (*PostPageResponse, error) {
	return gh.getPostPage(HomePage, page)
}

func (gh *geekHub) GetPostsPage(page int) (*PostPageResponse, error) {
	return gh.getPostPage(HomePage+PostsURI, page)
}

func (gh *geekHub) GetSecondHandsPage(page int) (*PostPageResponse, error) {
	return gh.getPostPage(HomePage+SecondHandsURI, page)
}

func (gh *geekHub) GetAuctionsPage(page int) (*PostPageResponse, error) {
	return gh.getPostPage(HomePage+AuctionsURI, page)
}

func (gh *geekHub) GetMoleculesPage(page int) (*PostPageResponse, error) {
	return gh.getPostPage(HomePage+MoleculesURI, page)
}

func (gh *geekHub) GetGroupBuysPage(page int) (*PostPageResponse, error) {
	return gh.getPostPage(HomePage+GroupBuysURI, page)
}

func (gh *geekHub) getPostPage(url string, page int) (*PostPageResponse, error) {
	res, err := httpClient.R().Get(url + "?page=" + fmt.Sprintf("%d", page))
	if err != nil {
		return nil, InternetError
	}
	body := res.Body()
	r := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, GoQueryError
	}

	response := &PostPageResponse{
		BasePageInfo: gh.getAuthFromHtml(doc),
		Posts:        nil,
	}

	doc.Find(gh.Selectors.PostList).Each(func(_ int, selection *goquery.Selection) {
		titleNode := selection.Find(gh.Selectors.Post)
		href, _ := titleNode.Find(gh.Selectors.PostHref).First().Attr("href")
		extraSetNode := titleNode.Siblings().First() //帖子下一条
		authorNode := extraSetNode.Find(gh.Selectors.TextUser).Eq(1)
		commentCountText := strings.TrimSpace(selection.Find(gh.Selectors.PostCommentCount).Last().Text())
		commentCount, _ := strconv.ParseInt(commentCountText, 10, 64)
		post := &PreviewPost{
			ID:              getPostIDFromUri(href),
			Title:           strings.TrimSpace(titleNode.First().Text()),
			Uri:             strings.TrimSpace(href),
			Author:          buildUserFromNode(authorNode),
			CommentCount:    commentCount,
			LatestReplyTime: extraSetNode.Find("span").Eq(3).Text(),
			LatestReplyUser: buildUserFromNode(extraSetNode.Find(gh.Selectors.TextUser).Last()),
		}
		response.Posts = append(response.Posts, post)
	})

	return response, nil
}

func getPostIDFromUri(href string) string {
	mm := strings.Split(strings.TrimSpace(href), "/")
	return mm[len(mm)-1]
}

func buildUserFromNode(node *goquery.Selection) *User {
	href, _ := node.Attr("href")
	return &User{
		Username: strings.TrimSpace(node.Text()),
		PageUri:  href,
	}
}

func (gh *geekHub) GetPostContent(pageUri string, page int) (*ContentPageResponse, error) {
	url := HomePage + pageUri
	if page > 0 {
		url += fmt.Sprintf("?page=%d", page)
	}
	doc, err := gh.getQueryDocFromUrl(url)
	if err != nil {
		return nil, err
	}

	content := Converter.Html2Md(doc.Find(gh.Selectors.PostBody))
	response := &ContentPageResponse{
		BasePageInfo: gh.getAuthFromHtml(doc),
		Post: &DetailPost{
			ID:          getPostIDFromUri(pageUri),
			Title:       strings.TrimSpace(doc.Find(gh.Selectors.PostPageTitle).Text()),
			PostType:    getPostType(pageUri),
			Content:     content,
			Uri:         pageUri,
			Author:      buildUserFromNode(doc.Find(gh.Selectors.PostPageAuthor)),
			PublishTime: strings.TrimSpace(doc.Find(".flex.items-center.mr-2").First().Text()),
			PV:          strings.TrimSpace(doc.Find(".flex.items-center.mr-2").Last().Text()),
		},
	}

	token, ok := doc.Find(gh.Selectors.ReplyToken).Attr("value")
	if ok {
		replyToken = token
	}

	commentCount := strings.TrimRight(strings.TrimSpace(doc.Find(gh.Selectors.PostPageCommentCount).Text()), " 回复")
	response.Post.CommentCount, _ = strconv.ParseInt(commentCount, 10, 64)
	response.Post.CommentTotalPage = response.Post.CommentCount/100 + 1
	curPage := strings.TrimSpace(doc.Find(gh.Selectors.CommentCurPage).Text())
	if curPage != "" {
		response.Post.CurCommentPage, _ = strconv.Atoi(curPage)
		Debugf("response curPage is %d", response.Post.CurCommentPage)
	} else {
		response.Post.CurCommentPage = 1
	}

	if response.Post.PostType == MoleculeType {
		MoleculesInfo := &MoleculesInfo{
			Name:        strings.TrimSpace(doc.Find(".flex.items-center.mb-2:nth-of-type(1) .flex-1:nth-of-type(2)").Text()),
			Price:       strings.TrimSpace(doc.Find(".flex.items-center.mb-5:nth-of-type(2) div:nth-of-type(2)").Text()),
			Molecule:    strings.TrimSpace(doc.Find(".flex-1.mt-5:nth-of-type(1) .flex.items-center:nth-of-type(4) div:nth-of-type(2)").Text()),
			Denominator: strings.TrimSpace(doc.Find(".flex-1.mt-5:nth-of-type(1) .flex.items-center:nth-of-type(5) div:nth-of-type(2)").Text()),
			HowToSend:   strings.TrimSpace(doc.Find(".flex-1.mt-5:nth-of-type(1) .flex.items-center:nth-of-type(6) div:nth-of-type(2)").Text()),
			Contact:     strings.TrimSpace(doc.Find(".flex-1.mt-5:nth-of-type(1) .flex.items-center:nth-of-type(7) div:nth-of-type(2)").Text()),
		}
		if doc.Find(".whitespace-no-wrap.mr-3").Length() == 1 {
			MoleculesInfo.Floor = strings.TrimSpace(doc.Find(".whitespace-no-wrap.mr-3").Siblings().First().Text())
		}

		script := strings.TrimSpace(doc.Find("script").Eq(3).Text())
		reg := regexp.MustCompile(`var seconds = \d+`)
		seconds := bytes.TrimLeft(reg.Find([]byte(script)), gh.Selectors.MoleculeSeconds)
		MoleculesInfo.CountDown, _ = strconv.Atoi(string(seconds))

		response.Post.ExtraInfo = MoleculesInfo
	}

	doc.Find(gh.Selectors.CommentList).Each(func(_ int, selection *goquery.Selection) {
		var parent *Comment
		parentUser := strings.TrimSpace(selection.Find(gh.Selectors.CommentParentUser).Text())
		if parentUser != "" {
			parent = &Comment{
				Author: &User{
					Username: "",
				},
				Content: strings.TrimSpace(selection.Find(gh.Selectors.CommentParentContent).Text()),
			}
		}
		response.Post.Comments = append(response.Post.Comments, &Comment{
			ID:          getCommentID(selection),
			Author:      buildUserFromNode(selection.Find(gh.Selectors.CommentAuthor)),
			CommentTime: selection.Find(gh.Selectors.CommentTime).Text(),
			Content:     strings.TrimSpace(Converter.Html2Md(selection.Find(gh.Selectors.CommentContent))),
			Floor:       strings.TrimSpace(selection.Find(gh.Selectors.CommentFloor).Text()),
			Parent:      parent,
		})
	})

	return response, nil
}

func getPostType(uri string) postType {
	if strings.Contains(uri, "posts") {
		return PostType
	} else if strings.Contains(uri, "second_hands") {
		return SecondHandType
	} else if strings.Contains(uri, "auctions") {
		return AuctionType
	} else if strings.Contains(uri, "molecules") {
		return MoleculeType
	} else if strings.Contains(uri, "group_buys") {
		return GroupBuyType
	}
	return ""
}

func getCommentID(selection *goquery.Selection) string {
	id, ok := selection.Attr("id")
	if !ok {
		return "0"
	} else {
		return strings.TrimLeft(id, "comment_")
	}
}

func (gh *geekHub) getQueryDocFromUrl(url string) (*goquery.Document, error) {
	res, err := httpClient.R().Get(url)
	if err != nil {
		return nil, InternetError
	}
	body := res.Body()

	r := bytes.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, GoQueryError
	}
	return doc, nil
}

func (gh *geekHub) getAuthFromHtml(doc *goquery.Document) *BasePageInfo {
	authStatusNode := doc.Find(gh.Selectors.AuthStatus)
	if strings.Contains(authStatusNode.Text(), "登录") {
		return &BasePageInfo{Status: NotLogin}
	} else {
		return &BasePageInfo{
			AuthInfo: &AuthInfo{
				Me:          buildUserFromNode(doc.Find(gh.Selectors.AuthUsername)),
				NotifyCount: strings.TrimSpace(doc.Find(gh.Selectors.ActivityCount).Text()),
			},
			Status: Logan,
		}
	}
}

func (gh *geekHub) GetActivities(page int) (*ActivitiesPageResponse, error) {
	doc, err := gh.getQueryDocFromUrl(HomePage + ActivitiesURI + "?page=" + fmt.Sprintf("%d", page))
	if err != nil {
		return nil, InternetError
	}

	response := &ActivitiesPageResponse{
		BasePageInfo: gh.getAuthFromHtml(doc),
		Activities:   nil,
	}

	doc.Find(gh.Selectors.Activities).Each(func(_ int, selection *goquery.Selection) {
		href, exist := selection.Find(gh.Selectors.ActivityTargetUri).First().Attr("href")
		if exist {
			a := &Activity{
				Type:        ReplyPost,
				TargetUri:   href,
				TargetTitle: strings.TrimSpace(selection.Find(gh.Selectors.ActivityTargetTitle).Text()),
				Content:     strings.TrimSpace(selection.Find(gh.Selectors.ActivityReplyContent).Text()),
				User:        buildUserFromNode(selection.Find(gh.Selectors.ActivityUser)),
				Time:        strings.TrimSpace(selection.Find(gh.Selectors.ActivityTime).Text()),
			}
			response.Activities = append(response.Activities, a)
		} else {
			content := strings.TrimSpace(selection.Find(gh.Selectors.ActivityBody).Text())
			if strings.HasPrefix(content, "您的积分订单") {
				a := &Activity{
					Type:        GbitOrder,
					TargetUri:   "",
					TargetTitle: "",
					Content:     content,
					User:        nil,
					Time:        strings.TrimSpace(selection.Find(gh.Selectors.ActivityTime).Text()),
				}
				response.Activities = append(response.Activities, a)
			} else if strings.HasPrefix(content, "抽 奖 请 进") {
				target := selection.Find(gh.Selectors.ActivityGetMoleculesLink).First()
				targetUri, _ := target.Attr("href")
				response.Activities = append(response.Activities, &Activity{
					Type:        GetMolecules,
					TargetUri:   targetUri,
					TargetTitle: content,
					Content:     "",
					User:        nil,
					Time:        strings.TrimSpace(selection.Find(gh.Selectors.ActivityTime).Text()),
				})
			} else if strings.HasPrefix(content, "您的分子") {
				target := selection.Find(gh.Selectors.ActivityGetMoleculesLink).First()
				targetUri, _ := target.Attr("href")
				response.Activities = append(response.Activities, &Activity{
					Type:        YourMoleculesFinish,
					TargetUri:   targetUri,
					TargetTitle: content,
					Content:     "",
					User:        nil,
					Time:        strings.TrimSpace(selection.Find(gh.Selectors.ActivityTime).Text()),
				})
			} else {
				response.Activities = append(response.Activities, &Activity{
					Type: Unknow,
				}) //保留一个空的，如果有新的出现，会看到明显的bug
			}
		}
	})
	return response, nil
}

type PostCommentArgs struct {
	AuthenticityToken string
	TargetType        postType
	TargetId          string
	ReplyToId         string
	Content           string
}

type postType string

const (
	PostType       postType = "Post"
	SecondHandType          = "SecondHand"
	AuctionType             = "Auction"
	MoleculeType            = "Molecule"
	GroupBuyType            = "GroupBuy"
)

func (gh *geekHub) PostComment(arg *PostCommentArgs) error {
	response, err := httpClient.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Cache-Control", "no-cache").
		SetFormData(map[string]string{
			"authenticity_token":   arg.AuthenticityToken,
			"comment[target_type]": string(arg.TargetType),
			"comment[target_id]":   arg.TargetId,
			"comment[reply_to_id]": arg.ReplyToId,
			"comment[content]":     arg.Content,
		}).
		Post(HomePage + PostCommentURI)

	Debugf("PostComment Response, status:%d", response.StatusCode())
	return err
}

func (gh *geekHub) GetSignStatus() (bool, string, error) {
	doc, err := gh.getQueryDocFromUrl(HomePage + CheckinsURI)
	if err != nil {
		return false, "", err
	}

	token, _ := doc.Find("head meta").Eq(3).Attr("content")
	content := doc.Find(gh.Selectors.CheckInButton).Text()
	if content == "签到" {
		return false, token, nil
	} else {
		return true, token, nil
	}
}

func (gh *geekHub) CheckIn(token string) error {
	response, err := httpClient.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Cache-Control", "no-cache").
		SetFormData(map[string]string{
			"_method":            "post",
			"authenticity_token": token,
		}).
		Post(HomePage + SignURI)
	Debugf("CheckIn response body:%s", string(response.Body()))
	return err
}

func (gh *geekHub) GetMePage(userUri string) (*MePageResponse, error) {
	doc, err := gh.getQueryDocFromUrl(HomePage + userUri)
	if err != nil {
		return nil, err
	}

	userDetail := &UserDetail{
		Star:  doc.Find("sidebar .box:nth-of-type(2)>div div:nth-of-type(2) div:nth-of-type(1)").Text(),
		Gbit:  doc.Find("sidebar .box:nth-of-type(2)>div div:nth-of-type(2) div:nth-of-type(2)").Text(),
		Score: doc.Find("sidebar .box:nth-of-type(2)>div div:nth-of-type(2) div:nth-of-type(3)").Text(),
	}

	response := &MePageResponse{
		BasePageInfo: gh.getAuthFromHtml(doc),
	}
	if response.BasePageInfo.AuthInfo != nil {
		response.BasePageInfo.AuthInfo.Me.UserDetail = userDetail
	}
	return response, nil
}
