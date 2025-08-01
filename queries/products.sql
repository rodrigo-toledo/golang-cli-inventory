-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: GetProductBySKU :one
SELECT * FROM products WHERE sku = $1;

-- name: ListProducts :many
SELECT * FROM products;

-- name: CreateProduct :one
INSERT INTO products (sku, name, description, price) 
VALUES ($1, $2, $3, $4) 
RETURNING *;

-- name: UpdateProduct :one
UPDATE products 
SET name = $2, description = $3, price = $4 
WHERE id = $1 
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;
