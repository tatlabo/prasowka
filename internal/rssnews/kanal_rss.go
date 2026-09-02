// Package rssnews to main unmarshal rss
package rssnews

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title         string `xml:"title"`
	Description   string `xml:"description"`
	Link          string `xml:"link"`
	Language      string `xml:"language"`
	PubDate       string `xml:"pubDate"`
	PubTime       time.Time
	LastBuildDate string `xml:"lastBuildDate"`
	LastBuilTime  time.Time
	Items         []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	GUID        int    `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	PubTime     time.Time
	Category    string    `xml:"category"`
	Enclosure   Enclosure `xml:"enclosure"`
}

type Enclosure struct {
	URL    string `xml:"url,attr"`
	Type   string `xml:"type,attr"`
	Length int64  `xml:"length,attr"`
}

func UnmarshalFeed(data []byte) (Feed, error) {
	var feed Feed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return Feed{}, err
	}

	return feed, nil
}

func ReadRRSS(url string) (b []byte, err error) {
	res, err := http.Get(url)
	if err != nil {
		return b, fmt.Errorf("error getting data from url")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return b, fmt.Errorf("error in data body")
	}

	return body, err
}

func CreateFeedTable(db *sql.DB) error {
	sql := `
	CREATE TABLE IF NOT EXISTS 
		channel (
		id INTEGER PRIMARY KEY,
		url TEXT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		link TEXT NOT NULL,
		language TEXT NOT NULL DEFAULT 'pl',
		pubDate TEXT NOT NULL,
		lastBuildDate TEXT NOT NULL) STRICT;`

	_, err := db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func CreateArticlesTable(db *sql.DB) error {
	sql := `
	CREATE TABLE IF NOT EXISTS 
		article (
	id INTEGER PRIMARY KEY,
	channel_id INTEGER,
	title TEXT NOT NULL, 
	description TEXT NOT NULL, 
	link TEXT NOT NULL, 
	guid INTEGER NOT NULL, 
	pubDate TEXT NOT NULL, 
	category TEXT NOT NULL, 
	enclosure TEXT NOT NULL,
	FOREIGN KEY (channel_id) REFERENCES channel(id)
	);`

	_, err := db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
