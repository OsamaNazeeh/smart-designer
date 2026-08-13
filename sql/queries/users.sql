-- name: CreateUser :one
INSERT INTO users(id, username, password, created_at, updated_at)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING id, username, created_at, updated_at;


-- name: GetUserByID :one 
SELECT id, username, created_at, updated_at
FROM users 
WHERE id = $1; 

-- name: GetUserByName :one
SELECT *
FROM users 
WHERE username = $1; 
