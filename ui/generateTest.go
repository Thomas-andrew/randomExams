package ui

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
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
	"github.com/Twintat/randomExams/ui/layouts"
)

func startRandomExam(g *data.GUI) {
	form := data.NewExam(g)
	newSet(form)
}

func isBook(id widget.TreeNodeID) bool {
	if len(id) < 1 {
		return false
	} else if id[0] == 'b' {
		return true
	}
	return false
}

func chaptersIDs(id widget.TreeNodeID) ([]widget.TreeNodeID, error) {
	// get ridof of prefix type
	intID, err := strconv.Atoi(id[1:])
	if err != nil {
		return nil, err
	}

	// open db
	dbSrc, err := db.OpenDB()
	if err != nil {
		return nil, err
	}
	defer dbSrc.Close()

	qdb := db.New(dbSrc)
	ctx := context.Background()

	// get chapters
	chapters, err := qdb.GetChapters(ctx, int64(intID))
	if err != nil {
		return nil, err
	}

	// add prefix type
	result := []string{}
	for _, chapter := range chapters {
		result = append(result, "c"+strconv.Itoa(int(chapter.ID)))
	}
	return result, nil
}

func isChapter(id widget.TreeNodeID) bool {
	if len(id) == 0 {
		return false
	} else if id[0] == 'c' {
		return true
	}
	return false
}

// make set of exercise
func newSet(form *data.Exam) {
	// help func to fail
	fail := func(err error) {
		slog.Error("[newSet]", "error", err)
		dialog.ShowError(err, form.Gui.Window)
	}

	dbSrc, err := db.OpenDB()
	if err != nil {
		fail(err)
	}
	defer dbSrc.Close()

	qdb := db.New(dbSrc)
	ctx := context.Background()

	books, err := qdb.GetBooks(ctx)
	if err != nil {
		fail(err)
	}

	set := data.SetTable{}

	// add prefix 'b' in the nodeID for indicating the type
	var bookIDs []string
	for _, book := range books {
		bookIDs = append(bookIDs, "b"+strconv.Itoa(int(book.ID)))
	}

	// creating the tree
	tree := widget.NewTree(
		// given an nodeID get the sub nodeIDs
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			switch {
			case id == "":
				return bookIDs
			case isBook(id):
				chapters, err := chaptersIDs(id)
				if err != nil {
					fail(err)
				}
				return chapters
			}
			return []string{}
		},
		// bool func to determine if is branch
		func(id widget.TreeNodeID) bool {
			if id == "" {
				return true
			}
			f := id[0]
			if f == 'b' {
				return true
			}
			return false
		},
		// base widget for the tree node
		func(_ bool) fyne.CanvasObject {
			return widget.NewCheck("", nil)
		},
		// set the content for the tree node
		func(id widget.TreeNodeID, isBranch bool, obj fyne.CanvasObject) {
			if len(id) == 0 {
				return
			}

			// nodeID prefix type
			treeType := id[0]
			realID, err := strconv.Atoi(id[1:])
			if err != nil {
				fail(err)
			}

			var checkChange func(b bool)
			var text string
			switch treeType {
			case 'b':
				checkChange = func(b bool) {
					set.AddBookID(realID)
				}
				// pulling book from db and geting the info
				book, err := qdb.GetBook(ctx, int64(realID))
				if err != nil {
					fail(err)
				}
				text = book.Info()
			case 'c':
				checkChange = func(b bool) {
					set.AddBookID(realID)
				}
				// pulling chapter from db and geting the info
				chapter, err := qdb.GetChapter(ctx, int64(realID))
				if err != nil {
					fail(err)
				}
				text = chapter.Info()
			default:
				checkChange = nil
			}
			if val, ok := obj.(*widget.Check); ok {
				val.SetText(text)
				val.OnChanged = checkChange
			}
		},
	)

	// saveButton := widget.NewButton(
	// 	"salvar",
	// 	func() {
	// 		// save set and go to next form
	// 		// form.Set = set
	// 		// defineExam(form)
	// 	},
	// )
	contButton := widget.NewButton(
		"continuar",
		func() {
			// generate pull of exercises
			pull, err := set.GenExerPull()
			if err != nil {
				fail(err)
			}
			form.Pull = append(form.Pull, pull...)
			defineExam(form)
		},
	)
	content := container.New(
		&newSetLayout{},
		tree,
		contButton,
	)
	form.Gui.Window.SetContent(content)
}

// chose existen sets
// func latestSets(form *data.Exam) data.ExerciseSets {
// 	return data.ExerciseSets{}
// }

func defineExam(form *data.Exam) {
	entryExNum := widget.NewEntry()
	entryExNum.SetPlaceHolder("quantos exercicios na prova?")
	entryDuration := widget.NewEntry()
	entryDuration.SetPlaceHolder("quanto tempo de prova?")

	contButton := widget.NewButton(
		"continuar",
		func() {
			numEx, err := strconv.Atoi(entryExNum.Text)
			if err != nil {
				slog.Error("[defineExam]", "error", err)
				dialog.ShowError(err, form.Gui.Window)
			}
			form.Num = numEx
			form.Duration = entryDuration.Text
			runExam(form)
		},
	)

	content := container.NewVBox(
		entryExNum,
		entryDuration,
		contButton,
	)
	form.Gui.Window.SetContent(content)
}

func runExam(form *data.Exam) {
	// build a exercise canvas

	list := container.NewVBox()

	chosen := []int{}
beforeLoop:
	for i := 0; i < form.Num; i++ {
		loteryID := rand.Intn(len(form.Pull))
		slog.Info("pulled random exercise", "id", loteryID)

		// test if loteryID already chosen
		for _, ch := range chosen {
			if loteryID == ch {
				// already matched, unchose!!!
				slog.Info("pulled id was already chosen")
				i--
				continue beforeLoop
			}
		}
		chosen = append(chosen, loteryID)

		exCanvas, err := exerciseCont(form.Pull[loteryID])
		if err != nil {
			slog.Error("[runExam]", "error", err)
			dialog.ShowError(err, form.Gui.Window)
		}
		list.Add(exCanvas)
	}
	topBarWidgets := newTimer(form.Gui, form.Duration)
	cont := container.New(
		&layouts.ExamLayout{},
		topBarWidgets.timer,
		topBarWidgets.bar, topBarWidgets.button,
		container.NewVScroll(list),
	)
	form.Gui.Window.SetContent(cont)
}

// return the exercise container with all images
func exerciseCont(ex db.Exercise) (fyne.CanvasObject, error) {
	fail := func(err error) (fyne.CanvasObject, error) {
		return nil, fmt.Errorf("[ExerciseCont] error creating container: %v", err)
	}
	cont := container.New(
		&ExerciseExamLayout{},
	)
	// open db
	dbSrc, err := db.OpenDB()
	if err != nil {
		return fail(err)
	}
	defer dbSrc.Close()

	qdb := db.New(dbSrc)
	ctx := context.Background()

	// get images
	images, err := qdb.GetImages(ctx, ex.ID)
	if err != nil {
		return fail(err)
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].Sequence < images[j].Sequence
	})

	for _, img := range images {
		slog.Info("images loaded", "path", img.FileName)
		img := canvas.NewImageFromFile(config.ImagesDirectory() + "/" + img.FileName)
		img.FillMode = canvas.ImageFillContain
		cont.Add(img)
	}
	return cont, nil
}
