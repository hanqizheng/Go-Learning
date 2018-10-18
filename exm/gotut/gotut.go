package main

import (
		"fmt"
		"net/http"
		"html/template"
		"io/ioutil"
		"encoding/xml"
		"sync"
)

var wg sync.WaitGroup

// SitemapIndex to 要抓取的内容结构
type SitemapIndex struct {
	Locations []string `xml:"sitemap>loc"`
}

// News to 新闻本身结构体
type News struct {
	Titles []string `xml:"url>news>title"`
	Keywords []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}

// NewsMap to 一堆新闻的结构体
type NewsMap struct {
	Keyword string
	Location string
}

// NewsAggPage to 新闻模板结构提=体
type NewsAggPage struct {
	Title string
	News map[string]NewsMap
}
// newsRoutine to
func newsRoutine(c chan News, Location string) {
		defer wg.Done()
		var n News
		resp, _ := http.Get(Location)
		bytes, _ := ioutil.ReadAll(resp.Body)		
		xml.Unmarshal(bytes, &n)
		resp.Body.Close()

		c <- n
}

// NewsAggHandler to 新闻模板路由的Controller
func newsAggHandler(w http.ResponseWriter, r *http.Request) {
		var s SitemapIndex

		newsMap := make(map[string]NewsMap)
		resp, _ := http.Get("https://www.washingtonpost.com/news-sitemap-index.xml")
		bytes, _ := ioutil.ReadAll(resp.Body)		
		xml.Unmarshal(bytes, &s)

		resp.Body.Close()

		queue := make(chan News, 30)

		for _, Location := range s.Locations {
				wg.Add(1)
				go newsRoutine(queue, Location)
		}
		
		wg.Wait()
		close(queue)

		for elem := range queue {
				for idx := range elem.Titles {
					newsMap[elem.Titles[idx]] = NewsMap{elem.Keywords[idx], elem.Locations[idx]}
				}
		}

		p := NewsAggPage{ Title: "A News Template", News: newsMap }
		t, _ := template.ParseFiles("basictemplating.html")
		fmt.Println(t.Execute(w, p))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<h1>hello</h1>")
}

func main()  {
		http.HandleFunc("/", indexHandler)
		http.HandleFunc("/agg/", newsAggHandler)
		http.ListenAndServe(":8000", nil)
}