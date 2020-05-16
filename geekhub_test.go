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
	spew.Dump(geekhub.GeekHub.GetPostContent("/molecules/70", 1))
}

func (GeekHubTestSuite) TestGetActivities() {
	spew.Dump(geekhub.GeekHub.GetActivities(1))
}

func (GeekHubTestSuite) TestGetUserPage() {
	spew.Dump(geekhub.GeekHub.GetMePage("/u/php"))
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
		Cookie: "_ga=GA1.2.2070442420.1588999418; _gid=GA1.2.919297677.1588999418; _session_id=CHgfVjgjyCY4WPHMKrN2DcaDH4eebgWq8xgvXEoWHOPasGhHYoczalfH3lTfIayYQUauUqL4QWNV0Pmm%2Fh2M9%2FFdcxohO2%2BUnpL0Ok1%2BdxnDfEUk9YLyxew8%2BPrK%2FcH%2F5Z7jkiVRQHjnFOTcplhvpm5pdDh6iFb%2BCFrmoWz%2Fp00Og%2FDwTmXwz%2F179ghejPuYsrR1th1kIpnZHYblkE18P5O7%2FgQGzIldts3L2gPY71lmlc2%2BrFZDKXXlF3bQmfWoqSSAvIw1wLnSTk6jJD0JCRGcbQTvD7Jl2lRay%2FXOQAKorxlHeEQPQ2xv1KY7%2BtxcwFFfZh0AKsRNIF2IUugOIzXdjFxns2Gl4zuUCE8cg5Tj0ttc3DUYAV%2FZ7nwv7kC4Qqg75BVpISA%2Bdu2vfUPlcie%2BaLEy4N40rEJZdxI%3D--mMaj6itFcow1%2FTqq--mkiJbIdPhXC0oyD2AP9wZg%3D%3D",
	})
}
