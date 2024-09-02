package data

type IngestForm struct {
	Gui *GUI

	IsNewBook bool
	Book      *BookInfo

	IsNewChapter bool
	Chapter      Chapter

	ExercisesNum []string
	ExerciseMap  map[int][]string
}
