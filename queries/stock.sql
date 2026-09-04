-- name: GetStockByProductAndLocation :one
SELECT * FROM stock WHERE product_id = $1 AND location_id = $2;

-- name: GetStockByProduct :many
SELECT * FROM stock WHERE product_id = $1;

-- name: GetStockByLocation :many
SELECT * FROM stock WHERE location_id = $1;

-- name: GetLowStock :many
SELECT * FROM stock WHERE quantity < $1;

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
INSERT INTO stock (product_id, location_id, quantity)
VALUES ($1, $2, $3)
ON CONFLICT (product_id, location_id)
DO UPDATE SET quantity = stock.quantity + EXCLUDED.quantity, updated_at = NOW()
RETURNING *;

-- name: RemoveStock :one
UPDATE stock
SET quantity = quantity - $3, updated_at = NOW()
WHERE product_id = $1 AND location_id = $2 AND quantity >= $3
RETURNING *;
