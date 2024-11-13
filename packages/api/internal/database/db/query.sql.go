// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package db

import (
	"context"
	"database/sql"
)

const deleteServer = `-- name: DeleteServer :exec
DELETE FROM servers WHERE id=?
`

func (q *Queries) DeleteServer(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteServer, id)
	return err
}

const getAvailablePorts = `-- name: GetAvailablePorts :many
SELECT id, port, created_at, server_id FROM ports WHERE server_id IS NULL ORDER BY port ASC
`

func (q *Queries) GetAvailablePorts(ctx context.Context) ([]Port, error) {
	rows, err := q.db.QueryContext(ctx, getAvailablePorts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Port
	for rows.Next() {
		var i Port
		if err := rows.Scan(
			&i.ID,
			&i.Port,
			&i.CreatedAt,
			&i.ServerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPortByValue = `-- name: GetPortByValue :one
SELECT id, port, created_at, server_id FROM ports WHERE port=?
`

func (q *Queries) GetPortByValue(ctx context.Context, port int64) (Port, error) {
	row := q.db.QueryRowContext(ctx, getPortByValue, port)
	var i Port
	err := row.Scan(
		&i.ID,
		&i.Port,
		&i.CreatedAt,
		&i.ServerID,
	)
	return i, err
}

const getPorts = `-- name: GetPorts :many
SELECT id, port, created_at, server_id FROM ports ORDER BY port ASC
`

func (q *Queries) GetPorts(ctx context.Context) ([]Port, error) {
	rows, err := q.db.QueryContext(ctx, getPorts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Port
	for rows.Next() {
		var i Port
		if err := rows.Scan(
			&i.ID,
			&i.Port,
			&i.CreatedAt,
			&i.ServerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getServers = `-- name: GetServers :many
SELECT s.id, s.max_players, s.admin_nickname, s.admin_password, s.expires_at, s.created_at, p.port FROM servers s
JOIN ports p on s.id = p.server_id
`

type GetServersRow struct {
	ID            int64
	MaxPlayers    sql.NullInt64
	AdminNickname string
	AdminPassword string
	ExpiresAt     sql.NullTime
	CreatedAt     sql.NullTime
	Port          int64
}

func (q *Queries) GetServers(ctx context.Context) ([]GetServersRow, error) {
	rows, err := q.db.QueryContext(ctx, getServers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetServersRow
	for rows.Next() {
		var i GetServersRow
		if err := rows.Scan(
			&i.ID,
			&i.MaxPlayers,
			&i.AdminNickname,
			&i.AdminPassword,
			&i.ExpiresAt,
			&i.CreatedAt,
			&i.Port,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertPort = `-- name: InsertPort :exec
INSERT INTO ports(port) VALUES(?)
`

func (q *Queries) InsertPort(ctx context.Context, port int64) error {
	_, err := q.db.ExecContext(ctx, insertPort, port)
	return err
}

const insertServer = `-- name: InsertServer :one
INSERT INTO servers(max_players, admin_nickname, admin_password) VALUES(?, ?, ?) RETURNING id, max_players, admin_nickname, admin_password, expires_at, created_at
`

type InsertServerParams struct {
	MaxPlayers    sql.NullInt64
	AdminNickname string
	AdminPassword string
}

func (q *Queries) InsertServer(ctx context.Context, arg InsertServerParams) (Server, error) {
	row := q.db.QueryRowContext(ctx, insertServer, arg.MaxPlayers, arg.AdminNickname, arg.AdminPassword)
	var i Server
	err := row.Scan(
		&i.ID,
		&i.MaxPlayers,
		&i.AdminNickname,
		&i.AdminPassword,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const resetPort = `-- name: ResetPort :exec
UPDATE ports SET server_id=NULL WHERE id=?
`

func (q *Queries) ResetPort(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, resetPort, id)
	return err
}

const updatePort = `-- name: UpdatePort :exec
UPDATE ports SET server_id=? WHERE id=?
`

type UpdatePortParams struct {
	ServerID sql.NullInt64
	ID       int64
}

func (q *Queries) UpdatePort(ctx context.Context, arg UpdatePortParams) error {
	_, err := q.db.ExecContext(ctx, updatePort, arg.ServerID, arg.ID)
	return err
}

const updateServerExpiration = `-- name: UpdateServerExpiration :exec
UPDATE servers SET expires_at=? WHERE id=?
`

type UpdateServerExpirationParams struct {
	ExpiresAt sql.NullTime
	ID        int64
}

func (q *Queries) UpdateServerExpiration(ctx context.Context, arg UpdateServerExpirationParams) error {
	_, err := q.db.ExecContext(ctx, updateServerExpiration, arg.ExpiresAt, arg.ID)
	return err
}
