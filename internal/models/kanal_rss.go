import (
	"encoding/xml"
	"fmt"
	"log"
	"strings"
)

func ScrapRSS(w Website) (a ArticleRender, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(w.Body))
	if err != nil {
		return a, err
	}

	a.Website = w

	title := htmlquery.FindOne(doc, "//h1[@class='article-title']")
	lead := htmlquery.FindOne(doc, "//p[@class='article-page__lead']")
	rest := htmlquery.Find(doc, "//div[@class='article-page__content article_speakable']/p")
	dateEllement := htmlquery.FindOne(doc, "//div[@class='article-date']/meta")

	var dateString string
	if dateEllement != nil {
		dateString = htmlquery.SelectAttr(dateEllement, "content")
	}

	a.Website.CreatedAt, err = time.Parse("2006-01-02T15:04:05", dateString)
	if err != nil {
		a.Website.CreatedAt = w.CreatedAt
	}

	if title != nil {
		a.Website.Title = htmlquery.InnerText(title)
	}
	if lead != nil {
		a.Lead = htmlquery.InnerText(lead)
	}

	for i := range rest {
		text := htmlquery.InnerText(rest[i])
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		a.Content = append(a.Content, text)
	}
	a.Body = ""

	return a, nil
}