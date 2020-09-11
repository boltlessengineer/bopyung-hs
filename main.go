package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/seongmin8452/school/bopyung-hs/scrapper"
)

type test struct {
	id   string
	text string
}

type tmp struct {
	Notices []scrapper.ExtractedPost `json:"notices"`
	News    []scrapper.ExtractedPost `json:"news"`
	Events  []scrapper.ExtractedPost `json:"events"`
	Meal    scrapper.MealMenu        `json:"meal"`
}

func homePage(c echo.Context) error {
	notices := scrapper.ScrapeNotices(1)
	news := scrapper.ScrapeNews()
	events := scrapper.ScrapeEvents()
	meal := scrapper.ScrapeMeal(time.Now())
	data := &tmp{
		Notices: notices,
		News:    news,
		Events:  events,
		Meal:    meal,
	}
	return c.JSON(http.StatusOK, data)
}

func main() {
	e := echo.New()

	e.GET("/", homePage)

	e.Logger.Fatal(e.Start(":8080"))
}
