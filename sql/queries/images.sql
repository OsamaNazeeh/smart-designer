-- name: CreateNewImage :one 
INSERT INTO images(id, objectKey, ext, created_at, updated_at, owner_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;


-- name: GetImage :one 
SELECT * FROM images
WHERE id = $1; 


-- name: UpdateImage :exec
UPDATE images
SET updated_at = $1, ext = $2

WHERE id = $3;