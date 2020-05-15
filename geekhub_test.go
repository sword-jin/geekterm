package geekhub_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	geekhub "github.com/rrylee/geekterm"
	"github.com/stretchr/testify/suite"
)

func TestGEekHubSuite(t *testing.T) {
	suite.Run(t, &GeekHubTestSuite{})
}

type GeekHubTestSuite struct {
	suite.Suite
}

func (GeekHubTestSuite) TestGetHomePage() {
	resp, err := geekhub.GeekHub.GetHomePage(1)
	spew.Dump(resp, err)
}

func (GeekHubTestSuite) TestGetMoleculesPage() {
	resp, err := geekhub.GeekHub.GetMoleculesPage(1)
	spew.Dump(resp, err)
}

func (GeekHubTestSuite) TestGetPostContent() {
	//spew.Dump(geekhub.GeekHub.GetPostContent("/posts/715"))
	spew.Dump(geekhub.GeekHub.GetPostContent("/molecules/58", 1))
}

func (GeekHubTestSuite) TestGetActivities() {
	spew.Dump(geekhub.GeekHub.GetActivities(1))
}

func (GeekHubTestSuite) TestPostComment() {
	spew.Dump(geekhub.GeekHub.PostComment(&geekhub.PostCommentArgs{
		AuthenticityToken: "0mqX5AVBWdd9pG5Lv4Ww3sIOsCRRLkDDC+I/QjMYM5P/89u1+XL22zJYZzWEjDzO8kN1tvY4PkL+68TIcJYvQQ==",
		TargetType:        "Molecule",
		TargetId:          "49",
		ReplyToId:         "0",
		Content:           "+1",
	}))
}

func (GeekHubTestSuite) TestGetSignStatus() {
	spew.Dump(geekhub.GeekHub.GetSignStatus())
}

func (GeekHubTestSuite) SetupSuite() {
	geekhub.Setup(&geekhub.Config{
	})
}
