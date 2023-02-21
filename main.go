package main

import (
    "fmt"
    "net/http"
    "log"
    "html/template"
)

type Data struct {
    Text string
    Url string
}

func handler(w http.ResponseWriter, r *http.Request) {
    //fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
    tpl, _ := template.ParseFiles("index.html")
    data := Data{Text:"Hello"}
    if r.Method == http.MethodPost {
        data.Url = r.FormValue("root_url")
    }
    tpl.Execute(w, data)
}

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("server running on 8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

