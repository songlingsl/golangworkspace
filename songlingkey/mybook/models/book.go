package models

type Book struct {
	BookId         string `gorm:"PRIMARY_KEY:book_Id"`
	Summary        string
	ImageUrl       string
	Recommendation string
	BookName       string
	MediaUrl       string
	ShowTime       string
	ReadCount      int
	BookType       string
}

func GetBooks() *[]Book {
	var list []Book
	//db.Limit(3).Order("read_count desc").Find(&list)
	db.Order("read_count desc").Find(&list)
	return &list
}

func GetContent(bookId string) string {
	type result struct {
		Content string
	}
	var r result
	db.Raw("SELECT content FROM book WHERE book_id = ?", bookId).Scan(&r)
	return r.Content
}
