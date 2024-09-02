package ui

import (
	"log/slog"

	"fyne.io/fyne/v2/widget"
	"github.com/Twintat/randomExams/data"
)

func addNewBook(form *data.IngestForm) {
	slog.Info("adding a new book for ingest")

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

	content := &widget.Form{
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
			form.IsNewBook = true
			form.Book = &data.BookInfo{
				Title:     titleEntry.Text,
				Author:    authorEntry.Text,
				Volume:    volumeEntry.Text,
				Edition:   editionEntry.Text,
				Publisher: publisherEntry.Text,
				Year:      yearEntry.Text,
			}
			form.Book.GenerateInfo()
			slog.Debug(
				"book enter into ingest form",
				"title", form.Book.Title,
				"author", form.Book.Author,
				"volume", form.Book.Volume,
				"edition", form.Book.Edition,
				"publisher", form.Book.Publisher,
				"year", form.Book.Year,
			)
			form.IsNewBook = true
			choseChapterOption(form)
		},
	}

	form.Gui.Window.SetContent(content)
}
