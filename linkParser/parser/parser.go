package parser

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type linkDef struct {
	Href string
	Text string
}

func getHref(attrs []html.Attribute) string {
	for _, el := range attrs {
		if el.Key == "href" {
			return el.Val
		}
	}

	return ""
}

func assembleText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}

	if node.Type != html.ElementNode {
		return ""
	}

	var t string

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		t += assembleText(c)
	}

	return strings.Join(strings.Fields(t), " ")
}

type Reader interface {
	io.Reader
	Cleanup()
}

func ParseHTML(r io.Reader) ([]linkDef, error) {
	doc, err := html.Parse(r)

	if err != nil {
		return nil, err
	}

	var result []linkDef

	var fn func(*html.Node)

	fn = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			text := assembleText(n)
			href := getHref(n.Attr)
			result = append(result, linkDef{Text: text, Href: href})
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fn(c)
		}

	}

	fn(doc)

	return result, nil
}
