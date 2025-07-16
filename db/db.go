package db

import (
	"context"
	"sandbox/db/models"

	"github.com/jackc/pgx/v5"
)

func Connect(ctx context.Context, uri string) (*pgx.Conn, *models.Queries, error) {
	dbConn, err := pgx.Connect(ctx, uri)
	if err != nil {
		return nil, nil, err
	}
	if err := dbConn.Ping(ctx); err != nil {
		dbConn.Close(ctx)
		return nil, nil, err
	}
	db := models.New(dbConn)
	return dbConn, db, nil
}
