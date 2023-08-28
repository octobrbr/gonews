// Пакет для работы с БД приложения GoNews.
package storage

import (
	"math/rand"
	"module37/gonews/pkg/models"
	"strconv"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	_, err := New()
	if err != nil {
		t.Fatal(err)
	}

}

func TestDB_StoreNews(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	posts := []models.Post{
		{
			Title: "Test Post",
			Link:  strconv.Itoa(rand.Intn(1_000_000_000)),
		},
	}
	db, err := New()
	if err != nil {
		t.Fatal(err)
	}
	err = db.StoreNews(posts)
	if err != nil {
		t.Fatal(err)
	}
	news, err := db.News(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", news)

}
