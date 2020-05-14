package geekhub_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	geekhub "github.com/rrylee/geekterm"
)

func TestCheckNewVersion(t *testing.T) {
	spew.Dump(geekhub.CheckNewVersion())
}
