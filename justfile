
db := "exercises.db"
main := "./cmd/app/main.go"

run:
    go run {{main}}

debug:
    go run -tags debug {{main}}
