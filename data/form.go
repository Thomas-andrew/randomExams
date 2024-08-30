package data

type IngestForm struct {
	IsNewBook bool
	Book      *BookInfo

	IsNewChapter bool
	Chapter      Chapter

	ExercisesNum []string
	ExerciseMap  map[int][]string
}
