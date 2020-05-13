package geekhub

import (
	"os/exec"
)

type Openable interface {
	GetUrl() string
}

type openableUrl struct {
	url string
}

func (o openableUrl) GetUrl() string {
	return o.url
}

func NewOpenableUrl(url string) Openable {
	return &openableUrl{url: url}
}

func OpenChrome(openable Openable) error {
	cmd := exec.Command("open", openable.GetUrl())
	return cmd.Run()
}
