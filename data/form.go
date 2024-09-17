package data

import (
	"context"
	"log/slog"

	"github.com/Twintat/randomExams/db"
)

type IngestForm struct {
	Gui *GUI

	IsNewBook bool
	Book      db.Book

	IsNewChapter bool
	Chapter      db.Chapter

	Exercises []db.Exercise

	Images []db.Image
}

func (form *IngestForm) SubmitToDB() error {
	// open dbSrc
	dbSrc, err := db.OpenDB()
	if err != nil {
		return err
	}
	defer dbSrc.Close()

	// create transaction
	ctx := context.Background()
	tx, err := dbSrc.Begin()
	if err != nil {
		return nil
	}
	defer tx.Rollback()

	qtx := db.New(tx)

	// submit book
	var bookID int64
	if form.IsNewBook {
		bookID, err = qtx.InsertBook(ctx, db.InsertBookParams{
			Title:     form.Book.Title,
			Author:    form.Book.Author,
			Volume:    form.Book.Volume,
			Edition:   form.Book.Edition,
			Publisher: form.Book.Publisher,
			Year:      form.Book.Year,
		})
		if err != nil {
			return err
		}
		slog.Debug(
			"ingest form, transaction add book",
			"bookID", bookID,
			"bookInfo", form.Book.Info(),
		)
	} else {
		slog.Debug(
			"[SubmitToDB] old book",
			"bookID", form.Book.ID,
		)
		bookID = form.Book.ID
	}

	// submit chapter
	var chapterID int64
	if form.IsNewChapter {
		chapterID, err = qtx.InsertChapter(ctx, db.InsertChapterParams{
			BookID: bookID,
			Number: form.Chapter.Number,
			Name:   form.Chapter.Name,
		})
		if err != nil {
			return err
		}
		slog.Debug(
			"ingest form, transaction, add chapter",
			"chapterID", chapterID,
			"name", form.Chapter.Name,
			"bookID", bookID,
			"number", form.Chapter.Number,
		)
	} else {
		slog.Debug("[SubmitToDB] old chapter")
		chapterID = form.Chapter.ID
	}

	// submit exercise and images
	for _, ex := range form.Exercises {
		// tmpID was set as ExID value for each image when taken the scrshoot
		tmpID := ex.ID
		exID, err := qtx.InsertExercise(ctx, db.InsertExerciseParams{
			ChapterID: chapterID,
			Number:    ex.Number,
		})
		if err != nil {
			return err
		}
		slog.Debug(
			"Ingest form, transaction, add exercise",
			"number", ex.Number,
			"chapterID", chapterID,
		)
		for _, img := range form.Images {
			if img.ExID == tmpID {
				err := qtx.InsertImage(ctx, db.InsertImageParams{
					ExID:     exID,
					FileName: img.FileName,
					Sequence: img.Sequence,
				})
				if err != nil {
					return err
				}
				slog.Debug(
					"Ingest form, transaction, add image",
					"exID", exID,
					"file", img.FileName,
					"squence", img.Sequence,
				)
			}
		}
	}

	return tx.Commit()
}
