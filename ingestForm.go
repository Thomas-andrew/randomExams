package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type dynamicForm struct {
	*fyne.Container

	book bookInfo
}

func makeIngestForm(g *GUI) {
	form := &dynamicForm{
		Container: container.NewVBox(),
	}
	form.selectBook(g)
}

func (d *dynamicForm) selectBook(g *GUI) {
	// clear form
	d.Objects = nil

	bookEntry := widget.NewEntry()
	bookLabel := widget.NewLabel("0")

	bookList, err := getBooks()
	if err != nil {
		dialog.ShowError(err, g.window)
		return
	}

	// var completion string

	bookEntry.OnChanged = func(text string) {
		subs := powerSet(text)

		sort.Slice(bookList, func(i, j int) bool {
			return bookScore(subs, bookList[i].info) > bookScore(subs, bookList[j].info)
		})

		var str string = ""

		for i, book := range bookList {
			str += fmt.Sprintf(
				"score: %v,\t %v\n",
				bookScore(subs, book.info),
				book.info,
			)
			// only the first 5 best match books
			if i > 5 {
				break
			}
		}

		bookLabel.SetText(str)
	}

	buttonGoBack := widget.NewButton(
		"voltar",
		func() {
			log.Println("go back to book selecting")
			d.selectBook(g)
		},
	)
	buttonSubimit := widget.NewButton(
		"avançar",
		func() {
			subs := powerSet(bookEntry.Text)
			sort.Slice(bookList, func(i, j int) bool {
				return bookScore(subs, bookList[i].info) > bookScore(subs, bookList[j].info)
			})
			bestMatchBook := bookList.bestMatch()
			log.Printf(
				"chosen: {bookID: %v, info: '%v'}\n",
				bestMatchBook.id,
				bestMatchBook.info,
			)
			d.book = bestMatchBook
			d.choseChapterOption(g)
		},
	)

	formSelector := container.NewHBox(
		buttonGoBack,
		buttonSubimit,
	)

	bookSearch := container.NewVBox(
		bookEntry,
		bookLabel,
		container.NewCenter(formSelector),
	)

	d.Add(bookSearch)
	d.Refresh()
	g.window.SetContent(d.Container)
}

func (d *dynamicForm) choseChapterOption(g *GUI) {
	// clear objects
	d.Objects = nil
	bookChosen := widget.NewLabel("livro escolhido: " + d.book.info)

	// list available chapters
	chapters, err := listChapters(d.book.id)
	if err != nil {
		dialog.ShowError(err, g.window)
		return
	}

	var chList string = "Capitulos no banco de dados\n\n"

	if len(chapters) == 0 {
		chList = "nenhum capitulo no banco de dados"
	} else {
		for k, v := range chapters {
			chList += strconv.Itoa(k) + " - " + v + "\n"
		}
	}

	chapterListLabel := widget.NewLabel(chList)

	newChapterButton := widget.NewButton(
		"adicionar capitulo",
		func() {
			d.newChapter(g)
		},
	)

	oldChapter := widget.NewButton(
		"capitulo antigo",
		func() {
			// d.choseOldChapter(g, d.book)
			log.Println("capitulo antigo")
		},
	)

	form := container.NewVBox(
		bookChosen,
		chapterListLabel,
		container.NewCenter(
			container.NewHBox(
				newChapterButton,
				oldChapter,
			),
		),
	)
	g.window.SetContent(form)
}

func (d *dynamicForm) newChapter(g *GUI) {
	d.Objects = nil
	bookChosen := widget.NewLabel("livro escolhido: " + d.book.info)

	chapterNumEntry := widget.NewEntry()
	chapterNumEntry.SetPlaceHolder("numero do capitulo")

	chapterNameEntry := widget.NewEntry()
	chapterNameEntry.SetPlaceHolder("nome do capitulo")

	chapterEntry := widget.NewLabel("")

	updateChapterEntry := func() {
		var str string = ""
		str += "numero do capitulo: " + chapterNumEntry.Text + "\n"
		str += "nome do capitulo:   " + chapterNameEntry.Text
		chapterEntry.SetText(str)
	}

	chapterNumEntry.OnChanged = func(s string) { updateChapterEntry() }
	chapterNameEntry.OnChanged = func(s string) { updateChapterEntry() }

	submitButton := widget.NewButton(
		"adicionar capitulo",
		func() {
			chapterNum, err := strconv.Atoi(chapterNumEntry.Text)
			if err != nil {
				dialog.ShowError(err, g.window)
				return
			}
			err = addChapters(d.book.id, chapterNum, chapterNameEntry.Text)
			if err != nil {
				dialog.ShowError(err, g.window)
				return
			}
			// go to the ask for how many screenshoots
		},
	)

	form := container.NewVBox(
		bookChosen,
		chapterNumEntry,
		chapterNameEntry,
		chapterEntry,
		container.NewCenter(submitButton),
	)
	g.window.SetContent(form)
}
