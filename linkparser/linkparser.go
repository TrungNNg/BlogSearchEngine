// linkparse package will parse all the links in <a> tag of an html file
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

// given an html page in string, this function will return a list of all <a> tag
func ParseAnchorTag(html_text string) []Link {
    //Parse returns the parse tree for the HTML from the given Reader. (*Node)
    node, err := html.Parse(strings.NewReader(html_text))
    if err != nil {
        fmt.Println("error parsing html text")
        panic(err)
    }
    links := []Link{}
    traverse(node, &links)
    return links
}

func traverse(node *html.Node, links *[]Link) {
    href := ""
    text := ""
    // check if current node is <a> tag
    if node.Type == html.ElementNode && node.Data == "a" {
        for _, att := range node.Attr {
            if att.Key == "href" { // save href value
                href = att.Val
            }
        }
        for c := node.FirstChild; c != nil; c = c.NextSibling {
            if c.Type == html.TextNode { // save text within <a> tag
                text += c.Data
            }
        }
        text = strings.TrimSpace(text)
        *links = append(*links, Link{Href:href, Text:text}) // save <a> tag content
    }
    // traverse the tree using DFS 
    for c := node.FirstChild; c != nil; c = c.NextSibling {
		traverse(c, links)
    }
}








