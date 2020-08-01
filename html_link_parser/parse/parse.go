package parse

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func collectText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += collectText(c)
	}
	return strings.Join(strings.Fields(ret), " ")
}

func GetLinks(r io.Reader) []Link {
	links := make([]Link, 0)
	root, _ := html.Parse(r)

	buildLink(root, &links)

	return links
}

func buildLink(n *html.Node, links *[]Link) {
	if n.Type == html.ElementNode && n.Data == "a" {
		link := Link{}
		for _, a := range n.Attr {
			if a.Key == "href" {
				link.Href = a.Val
			}
			break
		}
		link.Text = collectText(n)
		*links = append(*links, link)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		buildLink(c, links)
	}
}

func writeToJSON(l []Link, jsonFile io.Reader) {

}

func convertToJSON(l []Link) []byte {
	return nil
}
