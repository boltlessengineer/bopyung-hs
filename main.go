package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/seongmin8452/school/bopyung-hs/scrapper"
)

func homePage(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, world!")
}

func main() {
	//e := echo.New()

	//e.GET("/", homePage)

	//e.Logger.Fatal(e.Start(":8080"))
	scrapper.ScrapeNotices(1)
	scrapper.ScrapeNews()
	scrapper.ScrapeEvents()
}
