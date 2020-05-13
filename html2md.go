package geekhub

import (
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

type Convert interface {
	Html2Md(*goquery.Selection) string
}

type converter struct {
	engine *md.Converter
}

var Converter *converter

func (c *converter) Html2Md(s *goquery.Selection) string {
	return c.engine.Convert(s)
}



