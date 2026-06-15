-- name: CreateUser :one
INSERT INTO users (name, dob)
VALUES ($1, $2)
    RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users 
ORDER BY id;

-- name: UpdateUser :one
UPDATE users
SET name = $2,
    dob = $3
WHERE id = $1
    RETURNING *;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
    RETURNING id;

-- name: UserExists :one
SELECT EXISTS(
    SELECT 1
    FROM users
    WHERE id = $1
);

-- name: ListUsersPaginated :many
SELECT *
FROM users
ORDER BY id
    LIMIT $1
OFFSET $2;