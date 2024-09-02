
db := "exercises.db"

run:
    go run ./cmd/app/main.go

debug:
    go run -tags debug .

clean:
    rm ./imgs/*
    rm {{db}}
    sqlite3 {{db}} ".read schema.sql"
