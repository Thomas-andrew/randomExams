package database

import (
	"database/sql"
	"fmt"

	"github.com/Twintat/randomExams/data"
)

// get chapter info from a given chapterID
func GetChapter(chapterID int) (data.Chapter, error) {
	result := data.Chapter{}
	if chapterID == 0 {
		return result, nil
	}
	db, err := openDB()
	if err != nil {
		return result, fmt.Errorf("[GetChapter] %w", err)
	}
	defer db.Close()
	var bookID, number int
	var name string
	if err := db.QueryRow("SELECT bookID, number, name FROM chapters WHERE chapterID = ?",
		chapterID).Scan(&bookID, &number, &name); err != nil {
		if err == sql.ErrNoRows {
			return result, fmt.Errorf("cant retrive chapter %d: unknown chapter", chapterID)
		}
		return result, fmt.Errorf("cant retrive chapter %d: %v", chapterID, err)
	}
	result.Id = chapterID
	result.BookID = bookID
	result.Name = name
	result.Num = number
	result.GenerateInfo()
	return result, nil
}

// get the chapter of a given bookID from the db
func ListChapters(bookID int) (data.Chapters, error) {
	if bookID == 0 {
		return nil, nil
	}
	db, err := openDB()
	if err != nil {
		return nil, fmt.Errorf("[listChapters] %w", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT chapterID, number, name FROM chapters WHERE bookID = ?", bookID)
	if err != nil {
		return nil, fmt.Errorf("[listChapters] error retriving chapters of a book: %w", err)
	}
	defer rows.Close()

	chapters := make(data.Chapters, 0)

	for rows.Next() {
		var chapterID int
		var chapterNum int
		var chapterName string

		err = rows.Scan(&chapterID, &chapterNum, &chapterName)
		if err != nil {
			return nil, fmt.Errorf("[listChapters] error scanning chapterNum, chapterName: %w", err)
		}
		chapter := data.Chapter{Id: chapterID, Num: chapterNum, Name: chapterName}
		chapter.GenerateInfo()
		chapters = append(chapters, chapter)
	}
	return chapters, nil
}
