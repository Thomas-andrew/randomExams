package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/Twintat/randomExams/data"
	_ "github.com/mattn/go-sqlite3"
)

func openDB() (*sql.DB, error) {
	return sql.Open("sqlite3", "./exercises.db")
}

func SubmitToDB(d *data.IngestForm) error {
	db, err := openDB()
	if err != nil {
		return fmt.Errorf("[submitToDB] %w", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("[submitToDB] %w", err)
	}
	slog.Debug("open db transaction")
	// stage book info to db
	var bookID int
	if d.IsNewBook {
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
		slog.Debug("insert new book into db", "bookID", bookID)

		// add book infos
		infos := d.Book.GetInfos()
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
			slog.Debug(
				"insert book info into db",
				"bookID", bookID,
				"typeField", typeField,
				"content", content,
			)
		}
	} else {
		bookID = d.Book.Id
	}
	// stage chapter info to db
	var chapterID int
	if d.IsNewChapter {
		res, err := tx.Exec(
			"INSERT INTO chapters (bookID, number, name) VALUES (?, ?, ?)",
			bookID,
			d.Chapter.Num,
			d.Chapter.Name,
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
		slog.Debug(
			"insert new chapter into db",
			"bookID", bookID,
			"chapterNum", d.Chapter.Num,
			"chapterName", d.Chapter.Name,
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

	for _, exNum := range d.ExercisesNum {
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
		slog.Debug(
			"insert new exercise into db",
			"exID", resExID,
			"exNumStr", exNum,
			"exNumInt", exNumInt,
			"chapterID", chapterID,
		)
		// stage images info to db
		for imageOrder, imagePath := range d.ExerciseMap[exNumInt] {
			_, err := stmtExerData.Exec(int(resExID), imagePath, imageOrder)
			if err != nil {
				errTx := tx.Rollback()
				if errTx != nil {
					return fmt.Errorf("[submitToDB] 2 ERRORS:\n\tERROR1: %w\n\tERROR2: %w\n", err, errTx)
				}
				return fmt.Errorf("[submitToDB] %w", err)
			}
			slog.Debug(
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
	slog.Debug("Commit db transaction")

	return nil
}

func GetExercises(chapID int) (data.Exercises, error) {
	if chapID == 0 {
		return nil, NoID{objectName: "chapter"}
	}
	db, err := openDB()
	if err != nil {
		return nil, fmt.Errorf("[getExercises] %w", err)
	}
	defer db.Close()

	// get exercise IDs from this chapter
	query := "SELECT id, exNum FROM exerciseId WHERE chapterID = ?"
	rowsExerIDs, err := db.Query(query, chapID)
	if err != nil {
		return nil, fmt.Errorf("[getExercises] %w", err)
	}
	defer rowsExerIDs.Close()

	exers := data.NewExercises()
	// from exercise IDs get images
	for rowsExerIDs.Next() {
		// get exercise id and number
		var exID, exNum int
		err := rowsExerIDs.Scan(&exID, &exNum)
		if err != nil {
			return nil, fmt.Errorf("[getExercises] %w", err)
		}
		// from id and number get exercise order and image name
		query = "SELECT imageName, imageOrder FROM exerciseData WHERE exID = ?"
		rowsExerData, err := db.Query(query, exID)
		if err != nil {
			return nil, fmt.Errorf("[getExercises] %w", err)
		}
		ex := data.NewExercise(exID, exNum)
		for rowsExerData.Next() {
			var imageName string
			var imageOrder int
			err := rowsExerData.Scan(&imageName, &imageOrder)
			if err != nil {
				return nil, fmt.Errorf("[getExercises] %w", err)
			}
			ex.Images[imageOrder] = imageName
		}
		exers = append(exers, ex)
	}

	return exers, nil
}
