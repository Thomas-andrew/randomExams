package main

import (
	"database/sql"
	"log"
)

type bookInfoAtomDB struct {
	bookID    int
	typeField string
	content   string
}

func openDB() (*sql.DB, error) {
	return sql.Open("sqlite3", "./randomEx.db")
}

func getBooks() (bookInfos, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// first get all the books info
	DBentries, err := getBookInfo(db)
	if err != nil {
		return nil, err
	}

	// separete for bookID
	bookInfosDB := make(map[int]map[string]string)

	for _, entry := range DBentries {
		if book, exists := bookInfosDB[entry.bookID]; exists {
			if _, exists = book[entry.typeField]; !exists {
				bookInfosDB[entry.bookID][entry.typeField] = entry.content
			} else {
				log.Printf(
					"[getBooks] typeField already exists:\n\t{bookID: %v, typeField: '%v', content: '%v'}",
					entry.bookID,
					entry.typeField,
					entry.content,
				)
			}
		} else {
			bookInfosDB[entry.bookID] = make(map[string]string)
			bookInfosDB[entry.bookID][entry.typeField] = entry.content
		}
	}

	// concatenate string of info
	books := newBookInfos()
	//  key, val
	for id, infoMap := range bookInfosDB {
		var str string = ""
		str += infoMap["titulo"] + ", "
		str += infoMap["autor"] + ", "
		str += infoMap["volume"] + ", "
		str += infoMap["edição"] + ", "
		str += infoMap["editora"] + ", "
		str += infoMap["ano"]

		bk := bookInfo{
			id:   id,
			info: str,
		}

		books = append(books, bk)

	}
	return books, nil
}

func getBookInfo(db *sql.DB) ([]bookInfoAtomDB, error) {
	query := "SELECT bookID, typeField, content FROM bookInfo"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookInfosDB []bookInfoAtomDB

	for rows.Next() {
		var info bookInfoAtomDB
		err := rows.Scan(&info.bookID, &info.typeField, &info.content)
		if err != nil {
			return nil, err
		}
		bookInfosDB = append(bookInfosDB, info)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bookInfosDB, nil
}

func listChapters(bookID int) (map[int]string, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT number, name FROM chapters WHERE bookID = ?", bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chapters := make(map[int]string)

	for rows.Next() {
		var chapterNum int
		var chapterName string

		err = rows.Scan(&chapterNum, &chapterName)
		if err != nil {
			return nil, err
		}
		chapters[chapterNum] = chapterName
	}
	return chapters, nil
}

func addChapters(bookID, chapterNum int, name string) error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO chapters(bookID, number, name) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(bookID, chapterNum, name)
	if err != nil {
		return err
	}

	return nil
}

func insertBookInfoDB(bookID int64, typeField, content string) error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO bookInfo(bookID, typeField, content) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(bookID, typeField, content)
	if err != nil {
		return err
	}

	return nil
}

func insertbookIdDB() (int64, error) {
	db, err := openDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO bookId DEFAULT VALUES")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec()
	if err != nil {
		return 0, err
	}

	bookId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return bookId, nil
}
