package main

import (
	"fyne.io/fyne/v2/widget"
)

func (d *dynamicForm) addNewBook(g *GUI) {
	Logger.Info("adding a new book for ingest")

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
			d.isNewBook = true
			d.book = &bookInfo{
				title:     titleEntry.Text,
				author:    authorEntry.Text,
				volume:    volumeEntry.Text,
				edition:   editionEntry.Text,
				publisher: publisherEntry.Text,
				year:      yearEntry.Text,
			}
			d.book.generateInfo()
			Logger.Debug(
				"book enter into ingest form",
				"title", d.book.title,
				"author", d.book.author,
				"volume", d.book.volume,
				"edition", d.book.edition,
				"publisher", d.book.publisher,
				"year", d.book.year,
			)
			d.isNewBook = true
			d.choseChapterOption(g)
		},
	}

	g.window.SetContent(form)
}
