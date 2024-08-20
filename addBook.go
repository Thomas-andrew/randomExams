package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func makeAddBook(g *GUI) fyne.CanvasObject {
	titleEntry := widget.NewEntry()
	authorEntry := widget.NewEntry()
	volumeEntry := widget.NewEntry()
	editionEntry := widget.NewEntry()
	publisherEntry := widget.NewEntry()
	yearEntry := widget.NewEntry()

	results := widget.NewLabel("")

	updateResults := func() {
		var str string = ""
		str += "titulo:\t\t" + titleEntry.Text + "\n"
		str += "autor:\t\t" + authorEntry.Text + "\n"
		str += "volume:\t" + volumeEntry.Text + "\n"
		str += "edição:\t" + editionEntry.Text + "\n"
		str += "editora:\t" + publisherEntry.Text + "\n"
		str += "ano:\t\t" + yearEntry.Text + "\n"

		results.SetText(str)
	}
	titleEntry.OnChanged = func(s string) { updateResults() }
	authorEntry.OnChanged = func(s string) { updateResults() }
	volumeEntry.OnChanged = func(s string) { updateResults() }
	editionEntry.OnChanged = func(s string) { updateResults() }
	publisherEntry.OnChanged = func(s string) { updateResults() }
	yearEntry.OnChanged = func(s string) { updateResults() }

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "titulo:", Widget: titleEntry},
			{Text: "autor:", Widget: authorEntry},
			{Text: "volume:", Widget: volumeEntry},
			{Text: "edição:", Widget: editionEntry},
			{Text: "editora:", Widget: publisherEntry},
			{Text: "ano:", Widget: yearEntry},
			{Text: "resultado:", Widget: results},
		},
		OnSubmit: func() {
			db, err := sql.Open("sqlite3", "./randomEx.db")
			if err != nil {
				dialog.ShowError(err, g.window)
				return
			}
			id, err := insertbookIdDB(db)
			if err != nil {
				dialog.ShowError(err, g.window)
				return
			}
			values := map[string]string{
				"titulo":  titleEntry.Text,
				"autor":   authorEntry.Text,
				"volume":  volumeEntry.Text,
				"edição":  editionEntry.Text,
				"editora": publisherEntry.Text,
				"ano":     yearEntry.Text,
			}

			for key, val := range values {
				err := insertBookInfoDB(db, id, key, val)
				if err != nil {
					dialog.ShowError(err, g.window)
					return
				}
			}
			g.startPage()
		},
	}

	return form
}

func insertBookInfoDB(db *sql.DB, bookID int64, typeField, content string) error {
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

func insertbookIdDB(db *sql.DB) (int64, error) {
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
