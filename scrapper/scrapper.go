package scrapper

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extrectedPost struct {
	id     string
	title  string
	author string
	date   string
	views  int
	url    string
}

const baseURL = "http://www.bopyung.hs.kr/main.php"

// ScrapeNotice scrapes notice page
func ScrapeNotices(page int) {
	fmt.Println("Scraping Notices")

	query := "?menugrp=100100&master=bbs&act=list&master_sid=15&Page=" + strconv.Itoa(page)
	res, err := http.Get(baseURL + query)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	table := doc.Find(".bbsContent tbody")
	table.Find("tr").Each(extractPost)
}

// ScrapeNews scrapes news page
func ScrapeNews() {
	fmt.Println("Scraping News")
}

// ScrapeEvents scrapes events page
func ScrapeEvents() {
	fmt.Println("Scraping Events")
}

func extractPost(i int, post *goquery.Selection) {
	var id, title, author, date string
	var views int
	post.Find("td").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0:
			id = s.Text()
		case 1:
			// [ToDo] 내용물 태그
			// image/pdf/hwp/xlsx
		case 2:
			title = strings.TrimSpace(s.Text())
		case 3:
			author = s.Text()
		case 4:
			date = s.Text()
		case 5:
			views, _ = strconv.Atoi(s.Text())
		}
	})
	fmt.Println(id, title, author, date, views)
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}
