-- name: InsertBook :one
INSERT INTO books (
    title, author, volume, edition, publisher, year
) VALUES (
    ?, ?, ?, ?, ?, ?
)
RETURNING id;

-- name: InsertChapter :one
INSERT INTO chapters (
    book_id, number, name
) VALUES (
    ?, ?, ?
)
RETURNING id;

-- name: InsertExercise :one
INSERT INTO exercises (
    number, chapter_id
) VALUES (
    ?, ?
)
RETURNING id;

-- name: InsertImage :exec
INSERT INTO images (
    ex_id, file_name, sequence
) VALUES (
    ?, ?, ?
);

-- name: GetBooks :many
SELECT * FROM books;

-- name: GetBook :one
SELECT * FROM books WHERE id = ? LIMIT 1;

-- name: GetChapters :many
SELECT * FROM chapters
WHERE book_id = ?;

-- name: GetChapter :one
SELECT * FROM chapters WHERE id = ? LIMIT 1;

-- name: GetChapterIDs :many
SELECT id FROM chapters WHERE book_id = ?;

-- name: GetExercises :many
SELECT * FROM exercises
WHERE chapter_id = ?;

-- name: GetExeRange :many
SELECT number FROM exercises
WHERE chapter_id = ?;

-- name: GetImages :many
SELECT file_name, sequence FROM images
WHERE ex_id = ?;

