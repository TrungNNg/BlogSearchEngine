package linkparser

import (
    "fmt"
    "strings"

    "golang.org/x/net/html"
)

func Test() {
    fmt.Println("Hello")
}

type Link struct {
    Href string
    Text string
}

func ParseA(text string) []Link {
    doc, err := html.Parse(strings.NewReader(text))
    if err != nil {
        fmt.Println("error parsing html string")
        panic(err)
    }
    links := []Link{}
    tra(doc, &links)
    return links
}

func tra(n *html.Node, res *[]Link) {
    href := ""
    text := ""
    if n.Type == html.ElementNode && n.Data == "a" {
        for _, att := range n.Attr {
            if att.Key == "href" {
                href = att.Val
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            if c.Type == html.TextNode {
                text += c.Data
            }
        }
        text = strings.TrimSpace(text)
        *res = append(*res, Link{Href:href, Text:text})
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
		tra(c, res)
    }
}








