package scrapper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ExtractedPost struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Date   string `json:"date"`
	Views  int    `json:"views"`
	Url    string `json:"url"`
}

// [ToDo]
// 이거 여기가 아니라 다른데(modules라던가)에서 선언되야 하는게...?

type ExtractMenu struct {
	Menu []MealMenu `json:"menu"`
}

type MealMenu struct {
	Date   string   `json:"date"`
	Lunch  []string `json:"lunch"`
	Dinner []string `json:"dinner"`
}

const baseURL = "http://www.bopyung.hs.kr/"

// ScrapeNotice scrapes notice page
func ScrapeNotices(page int) []ExtractedPost {
	fmt.Println("Scraping Notices")
	return extractPostList("100100", "15", page)
}

// ScrapeNews scrapes news page
func ScrapeNews() []ExtractedPost {
	fmt.Println("Scraping News")
	return extractPostList("100200", "16", 1)
}

// ScrapeEvents scrapes events page
func ScrapeEvents() []ExtractedPost {
	fmt.Println("Scraping Events")
	return extractPostList("100300", "67", 1)
}

// ScrapeMeal gets meal
func ScrapeMeal(date time.Time) MealMenu {
	// [ToDo]
	// 그날 급식이 없을 시, 에러를 보내거나 해서 일단 없다는 것을 알려야함
	year, month, day := date.Date()
	fmt.Println("Scraping Today's meal")

	url := fmt.Sprintf("https://schoolmenukr.ml/api/high/J100005836?allergy=hidden&year=%d&month=%d&date=%d", year, month, day)
	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	var resJson ExtractMenu
	err = json.Unmarshal(body, &resJson)
	checkErr(err)

	var meal MealMenu
	if len(resJson.Menu) == 0 {
		meal = resJson.Menu[0]
	}
	fmt.Println(meal)

	return meal
}

func extractPostList(menugrp string, masterSid string, page int) []ExtractedPost {
	// [ToDo]
	// 인자값으로 검색옵션, 검색문자 등도 포함해서 걍 이미 있는 3개까지 합쳐서 map으로
	var posts []ExtractedPost
	query := "main.php?menugrp=" + menugrp + "&master=bbs&act=list&master_sid=" + masterSid + "&Page=" + strconv.Itoa(page)
	res, err := http.Get(baseURL + query)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	table := doc.Find(".bbsContent tbody")
	table.Find("tr").Each(func(i int, p *goquery.Selection) {
		post := extractPost(p)
		posts = append(posts, post)
	})

	return posts
}

func extractPost(post *goquery.Selection) ExtractedPost {
	var id, title, author, date, url string
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
			url, _ = s.Find("a").Attr("href")
			url = baseURL + url
		case 3:
			author = s.Text()
		case 4:
			date = s.Text()
		case 5:
			views, _ = strconv.Atoi(s.Text())
		}
	})
	fmt.Println(id, title, author, date, views)
	fmt.Println(url)
	return ExtractedPost{
		Id:     id,
		Title:  title,
		Author: author,
		Date:   date,
		Views:  views,
		Url:    url,
	}
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
