package model

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Crawler for crawler
type Crawler struct {
	Resp *http.Response
}

// GetJpy is crawler JPY from ESUN response
func (crawler *Crawler) GetJpy() (jpy float64, err error) {
	jpy = 1.0
	defer crawler.Resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(crawler.Resp.Body)
	s := doc.Find(os.Getenv("SELECTION"))
	jpyString := s.Text()
	jpy, err = strconv.ParseFloat(strings.TrimSpace(jpyString), 64)
	return
}
