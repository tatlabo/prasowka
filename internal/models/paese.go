package models

import (
	"html/template"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
)

func ParseSourceBody(w *Website) ([]Website, error) {

	doc, err := htmlquery.Parse(strings.NewReader(w.Body))
	if err != nil {
		return nil, err
	}
	list := htmlquery.Find(doc, "//div/div/article/h3")

	subpages := []Website{}

	for _, node := range list {
		subpage := Website{}

		a := htmlquery.FindOne(node, "//a")
		title := htmlquery.InnerText(a)
		title = strings.TrimSpace(title)
		link := htmlquery.SelectAttr(a, "href")

		subpage.Title = title
		subpage.URL = template.URL(link)
		subpage.CreatedAt = time.Now()

		subpage.SourceId = w.Id
		subpages = append(subpages, subpage)

	}

	return subpages, nil
}
