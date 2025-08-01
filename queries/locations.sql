-- name: GetLocationByID :one
SELECT * FROM locations WHERE id = $1;

-- name: GetLocationByName :one
SELECT * FROM locations WHERE name = $1;

-- name: ListLocations :many
SELECT * FROM locations;

-- name: CreateLocation :one
INSERT INTO locations (name) 
VALUES ($1) 
RETURNING *;

-- name: UpdateLocation :one
UPDATE locations 
SET name = $2 
WHERE id = $1 
RETURNING *;

-- name: DeleteLocation :exec
DELETE FROM locations WHERE id = $1;
