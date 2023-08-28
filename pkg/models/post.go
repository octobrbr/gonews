package models

// структура отдельной новости для БД
type Post struct {
	ID      int
	Title   string
	Content string
	PubTime int64
	Link    string
}
