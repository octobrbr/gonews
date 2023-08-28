package rss

import (
	"encoding/xml"
	"io"
	"module37/gonews/pkg/models"

	"net/http"
	"strings"
	"time"

	strip "github.com/grokify/html-strip-tags-go"
)

// получение массива новостей из rss
func Parse(url string) ([]models.Post, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var feed models.Feed

	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	var data []models.Post
	for _, item := range feed.Chanel.Items {
		var p models.Post
		p.Title = item.Title
		p.Content = item.Description
		p.Content = strip.StripTags(p.Content)
		p.Link = item.Link
		item.PubDate = strings.ReplaceAll(item.PubDate, ",", "")

		t, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			t, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", item.PubDate)
		}
		if err == nil {
			p.PubTime = t.Unix()
		}

		data = append(data, p)
	}
	return data, nil
}
