-- name: CreateTodo :one
INSERT INTO todos (title) VALUES ($1)
RETURNING id, title, is_done, created_at, updated_at;

-- name: ListTodos :many
SELECT * FROM todos
ORDER BY id DESC
LIMIT $1 OFFSET $2;

-- name: GetTodo :one
SELECT * FROM todos
WHERE id = $1;

-- name: UpdateTodo :one
UPDATE todos
SET title = $2,
    is_done = $3,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = $1;
