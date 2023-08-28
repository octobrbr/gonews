package rss

import (
	"testing"
)

func TestParse(t *testing.T) {
	feed, err := Parse("https://feed.infoq.com/")
	if err != nil {
		t.Fatal(err)
	}
	if len(feed) == 0 {
		t.Fatal("данные не раскодированы")
	}
	t.Logf("получено %d новостей\n%+v", len(feed), feed)

}
