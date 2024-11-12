-- name: GetPorts :many
SELECT * FROM ports ORDER BY port ASC;

-- name: GetAvailablePorts :many
SELECT * FROM ports WHERE server_id IS NULL ORDER BY port ASC;

-- name: InsertPort :exec
INSERT INTO ports(port) VALUES(?);

-- name: InsertServer :one
INSERT INTO servers(max_players, start_map) VALUES(?, ?) RETURNING *;

-- name: UpdatePort :exec
UPDATE ports SET server_id=? WHERE id=?;
