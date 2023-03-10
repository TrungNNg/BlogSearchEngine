package main

import (
    "strings"
    "fmt"
    "net/http"
    "log"
    "html/template"
    "net/url"
    "io"
    "regexp"
    "os"

    "golang.org/x/net/html"
    "github.com/TrungNNg/BlogSearchEngine/linkparser"
)

var host string
var port string
func init() {
    b, err := os.ReadFile(".env")
    if err != nil {
        panic("can not read config file")
    }
    l := strings.Split(string(b), "\n")
    host = l[0]
    port = l[1]
}

type Data struct {
    Text string
    Url string
    Keyword string
    Valid_url bool
    All_urls map[string]bool
    Url_contents map[string]string
    Url_contain_keyword []string

    Host string
}

func handler(w http.ResponseWriter, r *http.Request) {
    //fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
    tpl, _ := template.ParseFiles("index.html")
    data := Data{Text:"Blog Keyword Search Engine", Host:host}
    if r.Method == http.MethodPost {
        data.Url = r.FormValue("root_url")
        data.Keyword = r.FormValue("keyword")
        u, err := url.ParseRequestURI(data.Url)
        if err != nil {
            data.Valid_url = false
            tpl.Execute(w, data)
            return
        }
        data.Valid_url = true
        // url is valid, crawl the url and return all link with given keyword
        // only crawl with depth of 1, so it only check links that in the html of root_url
        data.All_urls = map[string]bool{}
        //fmt.Println("all link BEFORE crawl", data.All_urls)
        crawl(data.Url, 1, data.All_urls, u.Hostname(), u.Scheme) // we only find links with depth of 1
        //fmt.Println("all link AFTER crawl", data.All_urls)

        // build a map of [url] -> text in html of url
        data.Url_contents = map[string]string{}
        //fmt.Println("Urls_content BEFORE", data.Url_contents)
        for k := range data.All_urls {
            getText(k, data.Url_contents)
        }
        //fmt.Println("Urls_content AFTER ", data.Url_contents)

        // for each text of each URL use regexp to match keyword
        for k, v := range data.Url_contents {
            match, _ := regexp.MatchString(strings.ToLower(data.Keyword), v)
            if match {
                data.Url_contain_keyword = append(data.Url_contain_keyword, k)
            }
        }


        //fmt.Println("HERE",u.Hostname(), u.Scheme)

        tpl.Execute(w, data)
        return
    }
    tpl.Execute(w, data)
}

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Server listening on port",port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

// crawl will traverse all links using BFS from the root url and save all valid links to data.All_urls
// the m map[string]int use to keep track of url's depth
// the default is 1 depth, so we only save links in the given root url
func crawl(root_url string, depth int, All_urls map[string]bool, host, scheme string) {
    // queue use to BFS
    url_q := []string{root_url}
    checked := map[string]bool{root_url:true}
    
    // keep track of depth
    m := map[string]int{root_url:1}

    // save root url 
    All_urls[root_url] = true

    for len(url_q) != 0 {
        //p("checking", url_q[0])
        curr := url_q[0]    // get next url in queue
        url_q = url_q[1:]   // deque

        // if current url's depth > depth(1) break
        if m[curr] > depth {
            break
        }

        // get string from html from current url
        b, _ := GetHTML(curr)
        // get all link from html 
        // links type is []linkparser.Link
        links := linkparser.ParseAnchorTag(string(b))

        // there are duplicate links or links that point to different sites which need to filter
        filtered_links := filterLinks(links, host, scheme)
            
        // for all new links we found add to url_q to recursive crawl it
        for _, link := range filtered_links {
            if !checked[link] {
                url_q = append(url_q, link)

                // link in curr url will be 1 layer lower than current url
                m[link] = m[curr] + 1

                checked[link] = true
                
                // save valid url
                All_urls[link] = true
            }
        }
    }
}

// given a []linkparser.Link and hostname, scheme
// filter out link that are invalid, duplicate, different domain
func filterLinks(links []linkparser.Link, hostname, scheme string) []string {
    // added use to not add duplicate
    added := map[string]bool{}
    res := []string{}
    // remove dub and include only same hostname link
    for _, l := range links {
        u, err := url.Parse(l.Href)
        if err != nil {
            continue
        }
        // if same domain and not added then add to res
        if u.Hostname() == hostname && !added[l.Href] {
            fmt.Println("THIS SHOULD HIT")
            res = append(res, l.Href)
            fmt.Println(res)
            added[l.Href] = true
        } else if u.Hostname() == "" && !added[l.Href] {
            // this is for relative url
            // Forexample: href="/foo/bar.html" is valid url
            // however it need to change to scheme://hostname/foo/bar.html
            res = append(res, scheme + "://" + hostname + l.Href)
            added[l.Href] = true
        }
    }
    return res
}

// return html text with given url
func GetHTML(url string) ([]byte, error) {
    res, err := http.Get(url)
    if err != nil {
        //p("can not GET url")
        return []byte{}, err
    }
    b, err := io.ReadAll(res.Body)
    if err != nil {
        //p("can not read res.Body")
        return []byte{}, err
    }
    return b, nil
}

// getText build the Url_contents map of [url] -> text
func getText(url string, Url_contents map[string]string) {
    b, _ := GetHTML(url)
    if len(b) == 0 {
        return
    }
    node, err := html.Parse(strings.NewReader(string(b)))
    if err != nil {
        return
    }
    text := ""
    traverseText(node, &text)
    Url_contents[url] = text
}

func traverseText(node *html.Node, text *string) {
    if node.Type == html.TextNode {
        *text += strings.TrimSpace(strings.ToLower(node.Data))
    }
    // traverse the tree using DFS 
    for c := node.FirstChild; c != nil; c = c.NextSibling {
		traverseText(c, text)
    }
}



