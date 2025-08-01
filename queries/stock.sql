-- name: GetStockByProductAndLocation :one
SELECT * FROM stock WHERE product_id = $1 AND location_id = $2;

-- name: GetStockByProduct :many
SELECT * FROM stock WHERE product_id = $1;

-- name: GetStockByLocation :many
SELECT * FROM stock WHERE location_id = $1;

-- name: CreateStock :one
INSERT INTO stock (product_id, location_id, quantity) 
VALUES ($1, $2, $3) 
RETURNING *;

-- name: UpdateStock :one
UPDATE stock 
SET quantity = $3, updated_at = NOW() 
WHERE product_id = $1 AND location_id = $2 
RETURNING *;

-- name: DeleteStock :exec
DELETE FROM stock WHERE product_id = $1 AND location_id = $2;

-- name: AddStock :one
UPDATE stock 
SET quantity = quantity + $3, updated_at = NOW() 
WHERE product_id = $1 AND location_id = $2 
RETURNING *;

-- name: RemoveStock :one
UPDATE stock 
SET quantity = quantity - $3, updated_at = NOW() 
WHERE product_id = $1 AND location_id = $2 
RETURNING *;
