package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
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
	Logger.Info("retriving books from db")
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

			title:     infoMap["titulo"],
			author:    infoMap["autor"],
			volume:    infoMap["volume"],
			edition:   infoMap["edição"],
			publisher: infoMap["editora"],
			year:      infoMap["ano"],
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
	if bookID == 0 {
		return nil, nil
	}
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

func addChapters(bookID, chapterNum int, name string) (int, error) {
	db, err := openDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO chapters(bookID, number, name) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(bookID, chapterNum, name)
	if err != nil {
		return 0, err
	}

	chapterID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(chapterID), nil
}

func addExerciseID(chapterID, exerNum int) (int, error) {
	db, err := openDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO exerciseId(chapterID, exNum) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(chapterID, exerNum)
	if err != nil {
		return 0, err
	}
	exerID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(exerID), nil
}

func addExerciseImage(exerID int, imgName string) error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO exerciseData (exID, imageName) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(exerID, imgName)
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

func (d *dynamicForm) submitToDB() error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
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
				return fmt.Errorf("2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
			}
			return err
		}

		resBookID, err := res.LastInsertId()
		if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				return fmt.Errorf("2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
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
				return fmt.Errorf("2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
			}
			return err
		}
		defer stmtBookInfo.Close()
		for typeField, content := range infos {
			_, err = stmtBookInfo.Exec(bookID, typeField, content)
			if err != nil {
				errTx := tx.Rollback()
				if errTx != nil {
					return fmt.Errorf("2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
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
			d.chapterNum,
			d.chapterName,
		)
		if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				return fmt.Errorf("2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
			}
			return err
		}
		resChapterID, err := res.LastInsertId()
		if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				return fmt.Errorf("2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
			}
			return err
		}
		chapterID = int(resChapterID)
		Logger.Debug(
			"insert new chapter into db",
			"bookID", bookID,
			"chapterNum", d.chapterNum,
			"chapterName", d.chapterName,
		)
	} else {
		// // TODO: Add to dynamicForm chapterID  <21-08-24, twin>
	}
	// stage exercise info to db
	stmtIdExer, err := tx.Prepare("INSERT INTO exerciseId (exNum, chapterID) VALUES (?, ?)")
	if err != nil {
		errTx := tx.Rollback()
		if errTx != nil {
			return fmt.Errorf("2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
		}
		return err
	}
	stmtExerData, err := tx.Prepare("INSERT INTO exerciseData (exID, imageName) VALUES (?, ?)")
	if err != nil {
		errTx := tx.Rollback()
		if errTx != nil {
			return fmt.Errorf("2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
		}
		return err
	}

	for _, exNum := range d.exercisesNum {
		exNumInt, err := strconv.Atoi(exNum)
		if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				return fmt.Errorf("2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
			}
			return err
		}
		res, err := stmtIdExer.Exec(exNum, chapterID)
		if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				return fmt.Errorf("2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
			}
			return err
		}
		resExID, err := res.LastInsertId()
		if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				return fmt.Errorf("2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
			}
			return err
		}
		Logger.Debug(
			"insert new exercise into db",
			"exID", resExID,
			"exNumStr", exNum,
			"exNumInt", exNumInt,
			"chapterID", chapterID,
		)
		// stage images info to db
		for _, imagePath := range d.exerciseMap[exNumInt] {
			_, err := stmtExerData.Exec(int(resExID), imagePath)
			if err != nil {
				errTx := tx.Rollback()
				if errTx != nil {
					return fmt.Errorf("2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
				}
				return err
			}
			Logger.Debug(
				"inserting exercise image into db",
				"exID", resExID,
				"imagePath", imagePath,
			)
		}
	}

	// commit if sucefull
	err = tx.Commit()
	if err != nil {
		return err
	}
	Logger.Debug("Commit db transaction")

	return nil
}
