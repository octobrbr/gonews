// Сервер GoNews.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"module37/gonews/pkg/api"
	"module37/gonews/pkg/models"
	"module37/gonews/pkg/rss"
	"module37/gonews/pkg/storage"
)

// структура для параметров запуска
type config struct {
	URLS   []string `json:"rss"`
	Period int      `json:"request_period"`
}

func main() {

	ip := os.Getenv("APP_IP")
	port := os.Getenv("APP_PORT")

	// создаем объект БД
	db, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	api := api.New(db)

	// читаем параметры запуска нашего приложения
	b, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config config
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err)
	}

	// для каждой новостной ленты считываем новости в отдельном потоке
	chPosts := make(chan []models.Post)
	chErrs := make(chan error)
	for _, url := range config.URLS {
		go parseURL(url, chPosts, chErrs, config.Period)
	}

	// запускаем запись считанных новостей в БД
	go func() {
		for posts := range chPosts {
			db.StoreNews(posts)
		}
	}()

	// обрабатываем ошибки, возникшие при чтении новостей
	go func() {
		for err := range chErrs {
			log.Println("ошибка:", err)
		}
	}()

	// запускаем веб-сервер
	err = http.ListenAndServe(ip+":"+port, api.Router())
	if err != nil {
		log.Fatal(err)
	}
}

// func parseURL(url string, db *storage.DB, posts chan<- []models.Post, errs chan<- error, period int) {
func parseURL(url string, posts chan<- []models.Post, errs chan<- error, period int) {
	for {
		news, err := rss.Parse(url)
		if err != nil {
			errs <- err
			continue
		}
		posts <- news
		time.Sleep(time.Minute * time.Duration(period))
	}
}
