package main

// type ingestData struct {
// 	rows  []*ingestDataRow
// 	mapEx map[int][]string
//
// 	num int
// }

// func makeExIngest(g *GUI) fyne.CanvasObject {
// 	entry := widget.NewEntry()
// 	entry.SetPlaceHolder("Diga quantos exercicios:")
//
// 	book := widget.NewEntry()
// 	bookLabel := widget.NewLabel("0")
//
// 	db, err := sql.Open("sqlite3", "./randomEx.db")
// 	if err != nil {
// 		dialog.ShowError(err, g.window)
// 		return nil
// 	}
// 	bookList, err := getBooks(db)
// 	if err != nil {
// 		dialog.ShowError(err, g.window)
// 		return nil
// 	}
//
// 	// var completion string
//
// 	book.OnChanged = func(text string) {
// 		subs := powerSet(text)
//
// 		sort.Slice(bookList, func(i, j int) bool {
// 			return bookScore(subs, bookList[i]) > bookScore(subs, bookList[j])
// 		})
//
// 		var str string = ""
//
// 		for _, book := range bookList {
// 			str += fmt.Sprintf(
// 				"score: %v,\t %v\n",
// 				bookScore(subs, book),
// 				book,
// 			)
// 		}
//
// 		bookLabel.SetText(str)
// 	}
//
// 	bookSearch := container.NewVBox(
// 		book,
// 		bookLabel,
// 	)
//
// 	form := &widget.Form{
// 		Items: []*widget.FormItem{
// 			{Text: "exercicios", Widget: entry},
// 			{Text: "Livro", Widget: bookSearch},
// 		},
// 		OnSubmit: func() {
// 			numEntries, err := strconv.Atoi(entry.Text)
// 			if err != nil {
// 				dialog.ShowError(err, g.window)
// 			}
//
// 			if numEntries <= 0 {
// 				return
// 			}
//
// 			ingestData := &ingestData{
// 				num:   numEntries,
// 				mapEx: make(map[int][]string),
// 			}
//
// 			buttonTop := widget.NewButton(
// 				"Save",
// 				func() {
// 					confirmDialog := dialog.NewConfirm(
// 						"Confirmação para salvar",
// 						"Você quer salvar esses exercicios?",
// 						func(response bool) {
// 							if response {
// 								for i, row := range ingestData.rows {
// 									log.Printf("{ex: %v, path: %v}\n", i+1, row.CurrentImages())
// 								}
// 							} else {
// 								log.Println("noppers...")
// 							}
// 						},
// 						g.window,
// 					)
// 					confirmDialog.Show()
// 				},
// 			)
//
// 			vertList := container.New(layout.NewVBoxLayout())
//
// 			scrollVertList := container.NewScroll(vertList)
// 			border := container.NewBorder(buttonTop, nil, nil, nil, scrollVertList)
// 			g.window.SetContent(border)
//
// 			for i := 1; i <= numEntries; i++ {
// 				ingestDataRow := NewIngestData(i, g.window)
// 				ingestDataRow.AddImage(g.window)
//
// 				ingestData.rows = append(ingestData.rows, ingestDataRow)
//
// 				ingestRow := container.New(
// 					NewIngestRowLayout(),
// 					ingestDataRow.buttons,
// 					ingestDataRow.images,
// 				)
// 				vertList.Add(ingestRow)
// 				g.window.SetContent(border)
// 			}
// 		},
// 	}
// 	return form
// 	// fyne.Layout
// }
//
// type ingestDataRow struct {
// 	images  *fyne.Container
// 	buttons *fyne.Container
//
// 	path     string
// 	imgPaths []string
// 	num      int
// 	id       int
// }
//
// func NewIngestData(id int, w fyne.Window) *ingestDataRow {
// 	images := container.NewVBox()
// 	buttons := container.NewVBox()
//
// 	ingest := &ingestDataRow{
// 		images:  images,
// 		buttons: buttons,
//
// 		path: "./imgs/img_test-pog",
// 		num:  0,
// 		id:   id,
// 	}
//
// 	addButton := widget.NewButton(
// 		"Add image",
// 		func() {
// 			ingest.AddImage(w)
// 		},
// 	)
//
// 	ingest.buttons.Add(addButton)
//
// 	return ingest
// }
//
// func (g *ingestDataRow) CurrentImages() []string {
// 	return g.imgPaths
// }
//
// func (g *ingestDataRow) AddImage(w fyne.Window) {
// 	g.num += 1
// 	path := g.path + "-" + strconv.Itoa(g.id) + strconv.Itoa(g.num) + ".png"
// 	err := screenshoot(path)
// 	if err != nil {
// 		dialog.ShowError(err, w)
// 	}
//
// 	img := canvas.NewImageFromFile(path)
// 	img.SetMinSize(fyne.NewSize(700, 500))
// 	img.FillMode = canvas.ImageFillContain
// 	g.images.Add(img)
// 	g.imgPaths = append(g.imgPaths, img.File)
//
// 	retakeButton := widget.NewButton(
// 		fmt.Sprintf("retake %v", g.num),
// 		func() {
// 			err := screenshoot(path)
// 			if err != nil {
// 				dialog.ShowError(err, w)
// 			}
//
// 			img.Refresh()
// 		},
// 	)
//
// 	g.buttons.Add(retakeButton)
// }
