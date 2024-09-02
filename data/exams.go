package data

type Exam struct {
	Gui  *GUI
	Set  SetTable
	Pull Exercises
	Num  int
}

func NewExam(g *GUI) *Exam {
	return &Exam{
		Gui:  g,
		Set:  NewSetTable(),
		Pull: NewExercises(),
	}
}

type SetTable struct {
	BookIDs    []int
	ChapterIDs []int
}

func NewSetTable() SetTable {
	return SetTable{
		BookIDs:    []int{},
		ChapterIDs: []int{},
	}
}

func (e *SetTable) AddBookID(id int) {
	e.BookIDs = append(e.BookIDs, id)
}

func (e *SetTable) AddChapterID(id int) {
	e.ChapterIDs = append(e.ChapterIDs, id)
}
