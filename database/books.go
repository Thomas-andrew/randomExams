package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Twintat/randomExams/data"
)

type bookInfoAtomDB struct {
	bookID    int
	typeField string
	content   string
}

func GetBooks() (data.BookInfos, error) {
	slog.Info("retriving books from db")
	db, err := openDB()
	if err != nil {
		return nil, fmt.Errorf("[getBooks]->%w", err)
	}
	defer db.Close()

	// first get all the books info
	DBentries, err := getBookInfo(db)
	if err != nil {
		return nil, fmt.Errorf("[getBooks]->%w", err)
	}

	// separete for bookID
	bookInfosDB := make(map[int]map[string]string)

	for _, entry := range DBentries {
		if book, exists := bookInfosDB[entry.bookID]; exists {
			if _, exists = book[entry.typeField]; !exists {
				bookInfosDB[entry.bookID][entry.typeField] = entry.content
			}
		} else {
			bookInfosDB[entry.bookID] = make(map[string]string)
			bookInfosDB[entry.bookID][entry.typeField] = entry.content
		}
	}

	// concatenate string of info
	books := data.NewBookInfos()
	//  key, val
	for id, infoMap := range bookInfosDB {
		var str string = ""
		str += infoMap["title"] + ", "
		str += infoMap["author"] + ", "
		str += infoMap["volume"] + ", "
		str += infoMap["edition"] + ", "
		str += infoMap["publisher"] + ", "
		str += infoMap["year"]

		bk := data.BookInfo{
			Id:   id,
			Info: str,

			Title:     infoMap["title"],
			Author:    infoMap["author"],
			Volume:    infoMap["volume"],
			Edition:   infoMap["edition"],
			Publisher: infoMap["publisher"],
			Year:      infoMap["year"],
		}

		books = append(books, bk)

	}
	return books, nil
}

func getBookInfo(db *sql.DB) ([]bookInfoAtomDB, error) {
	query := "SELECT bookID, typeField, content FROM bookInfo"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("[getBookInfo] db query %w", err)
	}
	defer rows.Close()

	var bookInfosDB []bookInfoAtomDB

	for rows.Next() {
		var info bookInfoAtomDB
		err := rows.Scan(&info.bookID, &info.typeField, &info.content)
		if err != nil {
			return nil, fmt.Errorf("[getBookInfo] book info atom scan error: %w", err)
		}
		bookInfosDB = append(bookInfosDB, info)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("[getBookInfo] rows error: %w", err)
	}

	return bookInfosDB, nil
}

func GetBook(bookID int) (data.BookInfo, error) {
	result := data.BookInfo{}
	// retrive bookInfos of a bookID
	db, err := openDB()
	if err != nil {
		return result, fmt.Errorf("[GetBook] could not open db %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT typeField, content FROM bookInfo WHERE bookID = ?", bookID)
	if err != nil {
		return result, fmt.Errorf("[GetBook] %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var field, content string
		err := rows.Scan(&field, &content)
		if err != nil {
			return result, fmt.Errorf("[GetBook] %v", err)
		}
		switch field {
		case "title":
			result.Title = content
		case "author":
			result.Author = content
		case "volume":
			result.Volume = content
		case "publisher":
			result.Publisher = content
		case "edition":
			result.Publisher = content
		case "year":
			result.Year = content
		default:
			return data.BookInfo{}, fmt.Errorf(
				"[GetBook] typeField retrive doesn't match any existing field: %v", field)
		}
	}
	result.Id = bookID
	result.GenerateInfo()

	return result, nil
}
