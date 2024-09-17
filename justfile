
db := "exercises.db"
main := "./cmd/app/main.go"

run:
    go run {{main}}

debug:
    go run -tags debug {{main}}


reset_db:
    [[ "$(ls -A imgs)" ]] && rm -r imgs/* || true
    [[ -f "{{db}}" ]] && rm "{{db}}" || true
    sqlite3 "{{db}}" < schema.sql
