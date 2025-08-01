-- name: CreateStockMovement :one
INSERT INTO stock_movements (product_id, from_location_id, to_location_id, quantity, movement_type) 
VALUES ($1, $2, $3, $4, $5) 
RETURNING *;

-- name: ListStockMovements :many
SELECT * FROM stock_movements ORDER BY created_at DESC;

-- name: GetStockMovementsByProduct :many
SELECT * FROM stock_movements WHERE product_id = $1 ORDER BY created_at DESC;

-- name: GetStockMovementsByLocation :many
SELECT * FROM stock_movements WHERE from_location_id = $1 OR to_location_id = $1 ORDER BY created_at DESC;
