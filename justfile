
db := "exercises.db"

run:
    go run ./cmd/app/main.go

debug:
    go run -tags debug .
