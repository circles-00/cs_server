-- name: GetPorts :many
SELECT * FROM ports ORDER BY port ASC;

-- name: GetAvailablePorts :many
SELECT * FROM ports WHERE server_id IS NULL ORDER BY port ASC;

-- name: InsertPort :exec
INSERT INTO ports(port) VALUES(?);

-- name: InsertServer :one
INSERT INTO servers(max_players, admin_nickname, admin_password) VALUES(?, ?, ?) RETURNING *;

-- name: UpdatePort :exec
UPDATE ports SET server_id=? WHERE id=?;

-- name: GetServers :many
SELECT s.*, p.port FROM servers s
JOIN ports p on s.id = p.server_id;

-- name: DeleteServer :exec
DELETE FROM servers WHERE id=?;

-- name: ResetPort :exec
UPDATE ports SET server_id=NULL WHERE id=?;

-- name: GetPortByValue :one
SELECT * FROM ports WHERE port=?;

-- name: UpdateServerExpiration :exec
UPDATE servers SET expires_at=? WHERE id=?;
