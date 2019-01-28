package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type list []string

func checkError(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
}

func (s list) join(c string) string {
	return strings.Join(s, c)
}

func sendLinks(w http.ResponseWriter, r *http.Request) {
	links := list{}
	d, e := os.Getwd()
	checkError(e)
	file, err := os.Open(d)
	checkError(err)
	fn, _ := file.Readdirnames(0)
	ch := make(chan list)
	for _, f := range fn {
		if filepath.Ext(f) == ".html" {
			go retrieveLinks(f, ch)
			l := list(<-ch)
			if len(l) > 0 {
				links = append(links, l.join("<br/>"))
			}
		}
	}
	fmt.Fprintf(w, links.join("<br/>"))
}

func retrieveLinks(file string, c chan list) {
	r, _ := ioutil.ReadFile(file)
	expression, e := regexp.Compile("<a[^>]*>([^<]+)<\\/a>")
	checkError(e)
	c <- list(expression.FindAllString(string(r), -1))
}

func handleRequests() {
	http.HandleFunc("/", sendLinks)
	http.ListenAndServe(":8081", nil)
}

func main() {
	handleRequests()
}
