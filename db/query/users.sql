-- name: FindUsersByEmailOrPhone :many
SELECT * FROM users
WHERE 
  (email = $1 AND email IS NOT NULL AND $1 <> '') 
  OR 
  (phone_number = $2 AND phone_number IS NOT NULL AND $2 <> '')
ORDER BY created_at ASC;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUsersByPrimaryID :many
SELECT * FROM users
WHERE id = $1 OR linked_id = $1
ORDER BY link_precedence ASC, created_at ASC;

-- name: CreatePrimaryUser :one
INSERT INTO users (
  email,
  phone_number,
  link_precedence
) VALUES (
  $1, $2, 'primary'
)
RETURNING *;

-- name: CreateSecondaryUser :one
INSERT INTO users (
  email,
  phone_number,
  linked_id,
  link_precedence
) VALUES (
  $1, $2, $3, 'secondary'
)
RETURNING *;

-- name: UpdateUserToSecondary :exec
UPDATE users
SET 
  linked_id = $1,
  link_precedence = 'secondary',
  updated_at = NOW()
WHERE id = $2;
