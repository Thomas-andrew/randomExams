package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type bookInfoAtomDB struct {
	bookID    int
	typeField string
	content   string
}

func openDB() (*sql.DB, error) {
	return sql.Open("sqlite3", "./exercises.db")
}

func getBooks() (bookInfos, error) {
	Logger.Info("retriving books from db")
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
	books := newBookInfos()
	//  key, val
	for id, infoMap := range bookInfosDB {
		var str string = ""
		str += infoMap["title"] + ", "
		str += infoMap["author"] + ", "
		str += infoMap["volume"] + ", "
		str += infoMap["edition"] + ", "
		str += infoMap["publisher"] + ", "
		str += infoMap["year"]

		bk := bookInfo{
			id:   id,
			info: str,

			title:     infoMap["title"],
			author:    infoMap["author"],
			volume:    infoMap["volume"],
			edition:   infoMap["edition"],
			publisher: infoMap["publisher"],
			year:      infoMap["year"],
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

func (d *dynamicForm) submitToDB() error {
	db, err := openDB()
	if err != nil {
		return fmt.Errorf("[submitToDB] %w", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("[submitToDB] %w", err)
	}
	Logger.Debug("open db transaction")
	// stage book info to db
	var bookID int
	if d.isNewBook {
		// try add bookId
		res, err := tx.Exec("INSERT INTO bookId DEFAULT VALUES")
		if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				return fmt.Errorf("[submitToDB] 2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
			}
			return err
		}

		resBookID, err := res.LastInsertId()
		if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				return fmt.Errorf("[submitToDB] 2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
			}
			return err
		}
		bookID = int(resBookID)
		Logger.Debug("insert new book into db", "bookID", bookID)

		// add book infos
		infos := d.book.getInfos()
		stmtBookInfo, err := tx.Prepare("INSERT INTO bookInfo (bookID, typeField, content) VALUES (?, ?, ?)")
		if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				return fmt.Errorf("[submitToDB] 2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
			}
			return err
		}
		defer stmtBookInfo.Close()
		for typeField, content := range infos {
			_, err = stmtBookInfo.Exec(bookID, typeField, content)
			if err != nil {
				errTx := tx.Rollback()
				if errTx != nil {
					return fmt.Errorf("[submitToDB] 2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
				}
				return err
			}
			Logger.Debug(
				"insert book info into db",
				"bookID", bookID,
				"typeField", typeField,
				"content", content,
			)
		}
	} else {
		bookID = d.book.id
	}
	// stage chapter info to db
	var chapterID int
	if d.isNewChapter {
		res, err := tx.Exec(
			"INSERT INTO chapters (bookID, number, name) VALUES (?, ?, ?)",
			bookID,
			d.chapter.num,
			d.chapter.name,
		)
		if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				return fmt.Errorf("[submitToDB] 2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
			}
			return fmt.Errorf("[submitToDB] %w", err)
		}
		resChapterID, err := res.LastInsertId()
		if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				return fmt.Errorf("[submitToDB] 2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
			}
			return fmt.Errorf("[submitToDB] %w", err)
		}
		chapterID = int(resChapterID)
		Logger.Debug(
			"insert new chapter into db",
			"bookID", bookID,
			"chapterNum", d.chapter.num,
			"chapterName", d.chapter.name,
		)
	} else {
		// // TODO: Add to dynamicForm chapterID  <21-08-24, twin>
	}
	// stage exercise info to db
	stmtIdExer, err := tx.Prepare("INSERT INTO exerciseId (exNum, chapterID) VALUES (?, ?)")
	if err != nil {
		errTx := tx.Rollback()
		if errTx != nil {
			return fmt.Errorf("[submitToDB] 2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
		}
		return fmt.Errorf("[submitToDB] %w", err)
	}
	stmtExerData, err := tx.Prepare("INSERT INTO exerciseData (exID, imageName, imageOrder) VALUES (?, ?, ?)")
	if err != nil {
		errTx := tx.Rollback()
		if errTx != nil {
			return fmt.Errorf("[submitToDB] 2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
		}
		return fmt.Errorf("[submitToDB] %w", err)
	}

	for _, exNum := range d.exercisesNum {
		exNumInt, err := strconv.Atoi(exNum)
		if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				return fmt.Errorf("[submitToDB] 2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
			}
			return fmt.Errorf("[submitToDB] %w", err)
		}
		res, err := stmtIdExer.Exec(exNum, chapterID)
		if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				return fmt.Errorf("[submitToDB] 2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
			}
			return fmt.Errorf("[submitToDB] %w", err)
		}
		resExID, err := res.LastInsertId()
		if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				return fmt.Errorf("[submitToDB] 2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
			}
			return fmt.Errorf("[submitToDB] %w", err)
		}
		Logger.Debug(
			"insert new exercise into db",
			"exID", resExID,
			"exNumStr", exNum,
			"exNumInt", exNumInt,
			"chapterID", chapterID,
		)
		// stage images info to db
		for imageOrder, imagePath := range d.exerciseMap[exNumInt] {
			_, err := stmtExerData.Exec(int(resExID), imagePath, imageOrder)
			if err != nil {
				errTx := tx.Rollback()
				if errTx != nil {
					return fmt.Errorf("[submitToDB] 2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
				}
				return fmt.Errorf("[submitToDB] %w", err)
			}
			Logger.Debug(
				"inserting exercise image into db",
				"exID", resExID,
				"imageOrder", imageOrder,
				"imagePath", imagePath,
			)
		}
	}

	// commit if sucefull
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("[submitToDB] %w", err)
	}
	Logger.Debug("Commit db transaction")

	return nil
}
