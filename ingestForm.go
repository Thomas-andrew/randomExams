package main

import (
	"fmt"
	"sort"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type dynamicForm struct {
	isNewBook bool
	book      *bookInfo

	isNewChapter bool
	chapterNum   int
	chapterName  string

	exercisesNum []string
	exerciseMap  map[int][]string
}

func makeIngestForm(g *GUI) {
	Logger.Info("making new ingest Form")
	form := &dynamicForm{
		exerciseMap: make(map[int][]string),
	}
	form.selectBook(g)
}

func (d *dynamicForm) selectBook(g *GUI) {
	Logger.Info("[ingestForm:selectBook] selecting a book")
	bookEntry := widget.NewEntry()
	bookLabel := widget.NewLabel("")

	bookList, err := getBooks()
	if err != nil {
		dialog.ShowError(err, g.window)
		return
	}

	var bookListText string = ""
	for i, book := range bookList {
		bookListText += fmt.Sprintf(
			"score: 0,\t %v\n",
			book.info,
		)
		if i > 5 {
			break
		}
	}
	bookLabel.SetText(bookListText)

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

	buttonNewBook := widget.NewButton(
		"novo livro",
		func() {
			d.addNewBook(g)
		},
	)

	buttonGoBack := widget.NewButton(
		"voltar",
		func() {
			Logger.Info("[ingestForm:selectBook] reset selecting a book")
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
			Logger.Info("[ingestForm:selectBook] book chosen", "bookID", bestMatchBook.id, "BookInfo", bestMatchBook.info)

			d.isNewBook = false
			d.book = bestMatchBook
			d.choseChapterOption(g)
		},
	)

	formSelector := container.NewHBox(
		buttonGoBack,
		buttonNewBook,
		buttonSubimit,
	)

	bookSearch := container.NewVBox(
		bookEntry,
		bookLabel,
		container.NewCenter(formSelector),
	)

	g.window.SetContent(bookSearch)
}

func (d *dynamicForm) choseChapterOption(g *GUI) {
	bookChosen := widget.NewLabel("livro escolhido: " + d.book.info)

	// new Books don't have chapters saved
	if d.isNewBook {
		d.isNewChapter = true
		Logger.Info("[ingestForm:choseChapterOption] new book doesn't have chapters saved. Automaticly adding new chapter.")
		d.newChapter(g)
	}

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
			d.isNewChapter = true
			d.newChapter(g)
		},
	)

	oldChapter := widget.NewButton(
		"capitulo antigo",
		func() {
			d.isNewChapter = false
			// d.choseOldChapter(g, d.book)
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
	Logger.Info("[ingestForm:newChapter] adding new chapter")
	bookChosen := widget.NewLabel("livro escolhido: " + d.book.info)

	oldChapters := widget.NewLabel("")
	chapters, err := listChapters(d.book.id)
	if err != nil {
		dialog.ShowError(err, g.window)
		return
	}
	if !d.isNewChapter {
		var oldChaptersStr string = "Capitulos existentes\n"
		for num, name := range chapters {
			oldChaptersStr += strconv.Itoa(num) + " - " + name + "\n"
		}

		oldChapters.SetText(oldChaptersStr)
	}

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
			for num := range chapters {
				if chapterNum == num {
					dialog.ShowError(
						fmt.Errorf("chapter %v already exists! cannot have to of it.", chapterNum),
						g.window,
					)
					Logger.Warn("this chapter number already exists for this book",
						"func",
						"newChapter",
						"chapterNum",
						chapterNum,
					)
					d.newChapter(g)
					return
				}
			}
			d.chapterNum = chapterNum
			d.chapterName = chapterNameEntry.Text
			// go to the ask for how many screenshoots
			d.howManyExercises(g)
		},
	)

	form := container.NewVBox(
		bookChosen,
		oldChapters,
		chapterNumEntry,
		chapterNameEntry,
		chapterEntry,
		container.NewCenter(submitButton),
	)
	g.window.SetContent(form)
}

func (d *dynamicForm) howManyExercises(g *GUI) {
	bookChosen := widget.NewLabel("livro escolhido: " + d.book.info)
	chapterChosen := widget.NewLabel(
		"capitulo escolhido: " + strconv.Itoa(d.chapterNum) + " - " + d.chapterName,
	)

	numOfExercises := widget.NewEntry()
	numOfExercises.SetPlaceHolder("Quantos exercicios? Ex: '1-3,5-8'")

	submitButton := widget.NewButton(
		"salvar",
		func() {
			exerRanges, err := expandRanges(numOfExercises.Text)
			if err != nil {
				dialog.ShowError(err, g.window)
				d.howManyExercises(g)
			}
			d.exercisesNum = append(d.exercisesNum, exerRanges...)
			Logger.Info("chose how many exercises", "answer", len(exerRanges))
			d.takeScreenshoots(g)
		},
	)

	form := container.NewVBox(
		bookChosen,
		chapterChosen,
		numOfExercises,
		container.NewCenter(
			submitButton,
		),
	)

	g.window.SetContent(form)
}

func (d *dynamicForm) takeScreenshoots(g *GUI) {
	Logger.Info("enter func: takeScreenshoots")
	ingestExercises := &ingestData{
		num:   len(d.exercisesNum),
		mapEx: make(map[int][]string),
	}

	saveButton := widget.NewButton(
		"salvar",
		func() {
			d.exerciseMap = ingestExercises.retrivePaths()
			info := "livro: " + d.book.info + "\n"
			info += "capitulo: " + strconv.Itoa(d.chapterNum) + " - " + d.chapterName
			saveImgsDialog := dialog.NewConfirm(
				"Confirmação para salvar imagens",
				info,
				func(response bool) {
					if response {
						err := d.submitToDB()
						if err != nil {
							dialog.ShowError(err, g.window)
						}
					}
				},
				g.window,
			)
			saveImgsDialog.Show()
		},
	)

	vertList := container.NewVBox()
	vertListScroll := container.NewVScroll(vertList)
	border := container.NewBorder(
		saveButton,     // top
		nil,            // left
		nil,            // right
		nil,            // botton
		vertListScroll, // center
	)
	g.window.SetContent(border)

	for i := 1; i <= len(d.exercisesNum); i++ {
		exer := newExerciseRow(i, g)
		exer.AddImage(g)

		ingestExercises.rows = append(ingestExercises.rows, exer)
		ingestRow := container.New(
			NewIngestRowLayout(),
			exer.buttons,
			exer.images,
		)

		vertList.Add(ingestRow)
		g.window.SetContent(border)
	}
}

type ingestData struct {
	rows  []*exerciseRow
	mapEx map[int][]string

	num int
}

func (i *ingestData) retrivePaths() map[int][]string {
	result := make(map[int][]string)
	for _, row := range i.rows {
		result[row.exerciseNum] = append(result[row.exerciseNum], row.imgPaths...)
	}
	i.mapEx = result
	return result
}

type exerciseRow struct {
	images  *fyne.Container
	buttons *fyne.Container

	path        string
	imgPaths    []string
	numOfPhotos int
	exerciseNum int
}

func newExerciseRow(num int, g *GUI) *exerciseRow {
	images := container.NewVBox()
	buttons := container.NewVBox()

	ingest := &exerciseRow{
		images:  images,
		buttons: buttons,

		numOfPhotos: 0,
		exerciseNum: num,
	}

	addButton := widget.NewButton(
		"Adicionar imagem",
		func() {
			ingest.AddImage(g)
		},
	)

	ingest.buttons.Add(addButton)

	return ingest
}

func (g *exerciseRow) CurrentImages() []string {
	return g.imgPaths
}

func (e *exerciseRow) AddImage(g *GUI) {
	e.numOfPhotos += 1
	//      ./imgs/01012024-0101-000000.png
	imgName := imageName()
	path := imagesDirectory() + "/" + imgName
	err := screenshoot(path)
	if err != nil {
		dialog.ShowError(err, g.window)
	}

	img := canvas.NewImageFromFile(path)
	img.SetMinSize(fyne.NewSize(700, 500))
	img.FillMode = canvas.ImageFillContain
	e.images.Add(img)
	e.imgPaths = append(e.imgPaths, imgName)

	retakeButton := widget.NewButton(
		fmt.Sprintf("retake %v", e.numOfPhotos),
		func() {
			err := screenshoot(path)
			if err != nil {
				dialog.ShowError(err, g.window)
			}

			img.Refresh()
		},
	)

	e.buttons.Add(retakeButton)
}
