// Пакет для работы с БД приложения GoNews.
package storage

import (
	"context"
	"errors"
	"module37/gonews/pkg/models"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func New() (*DB, error) {

	connstr := os.Getenv("CONNSTR")

	if connstr == "" {
		return nil, errors.New("не указано подключение к БД")
	}
	pool, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	db := DB{
		pool: pool,
	}
	return &db, nil
}

// запись новостей в БД
func (db *DB) StoreNews(news []models.Post) error {
	for _, post := range news {
		_, err := db.pool.Exec(context.Background(), `
		INSERT INTO news(title, content, pub_time, link)
		VALUES ($1, $2, $3, $4)`,
			post.Title,
			post.Content,
			post.PubTime,
			post.Link,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// чтение n-последних новостей из БД
func (db *DB) News(n int) ([]models.Post, error) {
	if n == 0 {
		n = 10
	}
	rows, err := db.pool.Query(context.Background(), `
	SELECT id, title, content, pub_time, link FROM news
	ORDER BY pub_time DESC
	LIMIT $1
	`,
		n,
	)
	if err != nil {
		return nil, err
	}
	var news []models.Post
	for rows.Next() {
		var p models.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		news = append(news, p)
	}
	return news, rows.Err()
}
