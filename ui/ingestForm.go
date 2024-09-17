package ui

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Twintat/randomExams/config"
	"github.com/Twintat/randomExams/data"
	"github.com/Twintat/randomExams/db"
)

func makeIngestForm(g *data.GUI) {
	slog.Info("making new ingest Form")
	form := &data.IngestForm{
		Gui: g,
	}
	selectBook(form)
}

func selectBook(form *data.IngestForm) {
	slog.Debug("selecting a book")
	bookEntry := widget.NewEntry()
	bookLabel := widget.NewLabel("")
	// open dbSource
	dbSource, err := db.OpenDB()
	if err != nil {
		dialog.ShowError(err, form.Gui.Window)
		return
	}
	defer dbSource.Close()
	ctx := context.Background()
	qdb := db.New(dbSource)

	bookList, err := qdb.GetBooks(ctx)
	if err != nil {
		dialog.ShowError(err, form.Gui.Window)
		return
	}

	var bookListText string = ""
	for i, book := range bookList {
		bookListText += fmt.Sprintf(
			"score 0,\t %v\n",
			book.Info(),
		)
		if i > 5 {
			break
		}
	}
	bookLabel.SetText(bookListText)

	// var completion string

	bookEntry.OnChanged = func(text string) {
		subs := data.PowerSet(text)

		sort.Slice(bookList, func(i, j int) bool {
			return data.BookScore(subs, bookList[i].Info()) > data.BookScore(subs, bookList[j].Info())
		})

		var str string = ""

		for i, book := range bookList {
			str += fmt.Sprintf(
				"score %v,\t %v\n",
				data.BookScore(subs, book.Info()),
				book.Info(),
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
			addNewBook(form)
		},
	)

	buttonGoBack := widget.NewButton(
		"voltar",
		func() {
			slog.Debug("reset selecting a book")
			selectBook(form)
		},
	)
	buttonSubimit := widget.NewButton(
		"avançar",
		func() {
			subs := data.PowerSet(bookEntry.Text)
			sort.Slice(bookList, func(i, j int) bool {
				return data.BookScore(subs, bookList[i].Info()) > data.BookScore(subs, bookList[j].Info())
			})
			if len(bookList) == 0 {
				dialog.ShowError(
					fmt.Errorf("bookList has zero items"),
					form.Gui.Window,
				)
				return
			}
			bestMatchBook := bookList[0]
			slog.Info(
				"book chosen",
				"bookID", bestMatchBook.ID,
				"BookInfo", bestMatchBook.Info(),
			)

			form.IsNewBook = false
			form.Book = bestMatchBook
			choseChapterOption(form)
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

	form.Gui.Window.SetContent(bookSearch)
}

func choseChapterOption(form *data.IngestForm) {
	bookChosen := widget.NewLabel("livro escolhido: " + form.Book.Info())

	// new Books don't have chapters saved
	if form.IsNewBook {
		form.IsNewChapter = true
		slog.Info("new book doesn't have chapters saved. Automaticly adding new chapter.")
		newChapter(form)
		return
	}

	// open db
	dbSource, err := db.OpenDB()
	if err != nil {
		dialog.ShowError(err, form.Gui.Window)
		return
	}
	defer dbSource.Close()

	qdb := db.New(dbSource)
	ctx := context.Background()

	// list available chapters
	chapters, err := qdb.GetChapters(ctx, form.Book.ID)
	if err != nil {
		dialog.ShowError(err, form.Gui.Window)
		return
	}

	var chList string = "Capitulos no banco de dados\n\n"

	if len(chapters) == 0 {
		chList = "nenhum capitulo no banco de dados"
	} else {
		for _, chapter := range chapters {
			chList += chapter.Info() + "\n"
		}
	}

	chapterListLabel := widget.NewLabel(chList)

	newChapterButton := widget.NewButton(
		"adicionar capitulo",
		func() {
			form.IsNewChapter = true
			newChapter(form)
		},
	)

	oldChapter := widget.NewButton(
		"capitulo antigo",
		func() {
			form.IsNewChapter = false
			oldChapter(form)
		},
	)

	content := container.NewVBox(
		bookChosen,
		chapterListLabel,
		container.NewCenter(
			container.NewHBox(
				newChapterButton,
				oldChapter,
			),
		),
	)
	form.Gui.Window.SetContent(content)
}

func oldChapter(form *data.IngestForm) {
	slog.Debug("chosen old chapter")
	bookChosen := widget.NewLabel("livro escolhido: " + form.Book.Info())

	chapterEntry := widget.NewEntry()
	oldChapterListLabel := widget.NewLabel("")

	// open db
	dbSource, err := db.OpenDB()
	if err != nil {
		dialog.ShowError(err, form.Gui.Window)
		return
	}
	defer dbSource.Close()

	qdb := db.New(dbSource)
	ctx := context.Background()

	oldChapters, err := qdb.GetChapters(ctx, form.Book.ID)
	if err != nil {
		dialog.ShowError(err, form.Gui.Window)
		return
	}

	if len(oldChapters) == 0 {
		err := fmt.Errorf("[oldChapter] book chosen doesn't have old chapters")
		slog.Warn("len of oldChapters equal to 0", "ERROR", err)
		dialog.ShowError(
			err,
			form.Gui.Window,
		)
		newChapter(form)
		return
	}

	var msg string
	for _, oldChapter := range oldChapters {
		msg += oldChapter.Info() + "\n"
	}
	oldChapterListLabel.SetText(msg)

	// sort

	chapterEntry.OnChanged = func(s string) {
		subs := data.PowerSet(s)
		sort.Slice(
			oldChapters,
			func(i, j int) bool {
				return data.BookScore(subs, oldChapters[i].Info()) > data.BookScore(subs, oldChapters[j].Info())
			},
		)

		var str string = ""
		for _, chapter := range oldChapters {
			str += fmt.Sprintf(
				"s %v,\t %v\n",
				data.BookScore(subs, chapter.Info()),
				chapter.Info(),
			)
		}
		oldChapterListLabel.SetText(str)
	}

	submitButton := widget.NewButton(
		"escolher",
		func() {
			if len(oldChapters) == 0 {
				msg := "[oldChapter] oldChapters has zero lenght"
				slog.Error(msg)
				dialog.ShowError(
					fmt.Errorf(msg),
					form.Gui.Window,
				)
				return
			}
			form.Chapter = oldChapters[0]
			slog.Info(
				"chosen old chapter",
				"chapterInfo", form.Chapter.Info(),
			)
			howManyExercises(form)
			return
		},
	)

	content := container.NewVBox(
		bookChosen,
		chapterEntry,
		oldChapterListLabel,
		container.NewCenter(
			submitButton,
		),
	)
	form.Gui.Window.SetContent(content)
}

func newChapter(form *data.IngestForm) {
	slog.Info("adding new chapter")
	bookChosen := widget.NewLabel("livro escolhido: " + form.Book.Info())

	// open db
	dbSource, err := db.OpenDB()
	if err != nil {
		dialog.ShowError(err, form.Gui.Window)
		return
	}
	defer dbSource.Close()

	qdb := db.New(dbSource)
	ctx := context.Background()

	oldChapters := widget.NewLabel("")
	chapters, err := qdb.GetChapters(ctx, form.Book.ID)
	if err != nil {
		dialog.ShowError(err, form.Gui.Window)
		return
	}
	var oldChaptersStr string = "Capitulos existentes\n"
	for _, chapter := range chapters {
		oldChaptersStr += chapter.Info() + "\n"
	}

	oldChapters.SetText(oldChaptersStr)

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
				dialog.ShowError(err, form.Gui.Window)
				return
			}
			for num := range chapters {
				if chapterNum == int(chapters[num].Number) {
					dialog.ShowError(
						fmt.Errorf("chapter %v already exists! cannot have to of it.", chapterNum),
						form.Gui.Window,
					)
					slog.Warn("this chapter number already exists for this book",
						"func",
						"newChapter",
						"chapterNum",
						chapterNum,
					)
					newChapter(form)
					return
				}
			}
			chapter := db.Chapter{
				BookID: form.Book.ID,
				Name:   chapterNameEntry.Text,
				Number: int64(chapterNum),
			}
			form.Chapter = chapter
			// go to the ask for how many screenshoots
			howManyExercises(form)
		},
	)

	content := container.NewVBox(
		bookChosen,
		oldChapters,
		chapterNumEntry,
		chapterNameEntry,
		chapterEntry,
		container.NewCenter(submitButton),
	)
	form.Gui.Window.SetContent(content)
}

func howManyExercises(form *data.IngestForm) {
	slog.Debug("entering howManyExercises")

	// display chosen book and chapter
	bookChosen := widget.NewLabel("livro escolhido: " + form.Book.Info())
	chapterChosen := widget.NewLabel("capitulo escolhido: " + form.Chapter.Info())

	// entry for how many exercises as ranges
	numOfExercises := widget.NewEntry()
	numOfExercises.SetPlaceHolder("Quantos exercicios? Ex: '1-3,5-8'")

	submitButton := widget.NewButton(
		"salvar",
		func() {
			slog.Debug("running the saving function")
			// decode range
			exerRanges, err := expandRanges(numOfExercises.Text)
			if err != nil {
				dialog.ShowError(err, form.Gui.Window)
				howManyExercises(form)
			}
			// check for colisions between range chosen and the exercises in db
			err = checkRanges(form.Chapter, exerRanges)
			if err != nil {
				if _, ok := err.(ExerciseColisions); ok {
					slog.Warn("conflicting range of of exercises passed")
					dialog.ShowError(
						fmt.Errorf("[howManyExercises] %w", err),
						form.Gui.Window,
					)
					howManyExercises(form)
					return
				} else {
					slog.Error("error in checkRanges", "error", err)
					dialog.ShowError(
						fmt.Errorf("[howManyExercises] %w", err),
						form.Gui.Window,
					)
					return
				}
			}
			// range valid and no colisions
			// create exercises for form
			for tmpID, num := range exerRanges {
				form.Exercises = append(form.Exercises, db.Exercise{
					ID:     int64(tmpID), // temporary id for exercise
					Number: int64(num),
				})
			}
			slog.Info("chose how many exercises", "answer", len(exerRanges))
			takeScreenshoots(form)
		},
	)

	content := container.NewVBox(
		bookChosen,
		chapterChosen,
		numOfExercises,
		container.NewCenter(
			submitButton,
		),
	)

	form.Gui.Window.SetContent(content)
}

func takeScreenshoots(form *data.IngestForm) {
	slog.Info("enter func: takeScreenshoots")

	saveButton := widget.NewButton(
		"salvar",
		func() {
			info := "livro: " + form.Book.Info() + "\n"
			info += "capitulo: " + form.Chapter.Info() + "\n"
			info += "exercicios: "
			for _, ex := range form.Exercises {
				info += strconv.Itoa(int(ex.Number)) + ", "
			}
			info += "\n"
			saveImgsDialog := dialog.NewConfirm(
				"Confirmação para salvar imagens",
				info,
				func(response bool) {
					if response {
						err := form.SubmitToDB()
						if err != nil {
							dialog.ShowError(err, form.Gui.Window)
						}
						StartPage(form.Gui)
					}
				},
				form.Gui.Window,
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
	form.Gui.Window.SetContent(border)

	for _, ex := range form.Exercises {
		// creating a exercise image container
		// display the image + button for adding more images
		exer := newExerciseRow(ex, form)
		exer.AddImage(form)

		// adding the exercise/image container to the pull
		ingestRow := container.New(
			NewIngestRowLayout(),
			exer.buttons,
			exer.images,
		)

		vertList.Add(ingestRow)
		vertListScroll.ScrollToBottom()
		form.Gui.Window.SetContent(border)
	}
}

type exerciseRow struct {
	images  *fyne.Container
	buttons *fyne.Container

	exercise db.Exercise
	num      int64 // number of images
}

func newExerciseRow(ex db.Exercise, form *data.IngestForm) *exerciseRow {
	images := container.NewVBox()
	buttons := container.NewVBox()

	ingest := &exerciseRow{
		images:  images,
		buttons: buttons,

		exercise: ex,
	}

	addButton := widget.NewButton(
		"Adicionar imagem",
		func() {
			ingest.AddImage(form)
		},
	)

	ingest.buttons.Add(addButton)

	return ingest
}

func (e *exerciseRow) AddImage(form *data.IngestForm) {
	// e.numOfPhotos += 1
	//      ./imgs/01012024-0101-000000.png
	imgName := imageName()
	path := config.ImagesDirectory() + "/" + imgName
	err := screenshoot(path)
	if err != nil {
		dialog.ShowError(err, form.Gui.Window)
	}

	img := canvas.NewImageFromFile(path)
	img.SetMinSize(fyne.NewSize(700, 500))
	img.FillMode = canvas.ImageFillContain
	e.images.Add(img)

	e.num += 1
	form.Images = append(form.Images,
		db.Image{
			ExID:     e.exercise.ID, // tmpID
			FileName: imgName,
			Sequence: e.num,
		},
	)

	retakeButton := widget.NewButton(
		fmt.Sprintf("retake %v", e.num),
		func() {
			err := screenshoot(path)
			if err != nil {
				dialog.ShowError(err, form.Gui.Window)
			}

			img.Refresh()
		},
	)

	e.buttons.Add(retakeButton)
}
