-- name: GetOneArticleById :one
SELECT *
FROM articles
WHERE id = $1 LIMIT 1;

-- name: GetAllArticles :many
SELECT *
FROM articles
ORDER BY title;

-- name: CreateArticle :one
INSERT INTO articles (id, title, url, savedOn, readOn)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: DeleteArticle :exec
DELETE
FROM articles
WHERE id = $1;

-- name: UpdateArticle :one
UPDATE articles
SET title = $1,
    url = $2,
    savedOn = $3,
    readOn = $4
WHERE id = $1
    RETURNING *;