// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
)

type Port struct {
	ID        int64
	Port      int64
	CreatedAt sql.NullTime
	ServerID  sql.NullInt64
}

type Server struct {
	ID         int64
	MaxPlayers sql.NullInt64
	StartMap   sql.NullString
	CreatedAt  sql.NullTime
}