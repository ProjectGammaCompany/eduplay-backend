package postgres

import (
	"context"
	"strings"

	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/lib/pq"

	pgx "github.com/jackc/pgx/v5"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	poolConfig, err := pgxpool.ParseConfig(storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s - %s", op, err)
	}

	poolConfig.MaxConns = 13
	poolConfig.MinConns = 5
	db, err := pgxpool.NewWithConfig(ctx, poolConfig)

	if err != nil {
		return nil, fmt.Errorf("%s - %s", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop(ctx context.Context) error {
	s.db.Close()
	return nil
}

func (s *Storage) SaveFile(ctx context.Context, fileName string, fileUUID string) (string, error) {
	const op = "storage.postgres.SaveFile"

	var id = strings.Split(fileUUID, ".")[0]

	state := `SELECT count FROM files WHERE fileId = $1`

	res := s.db.QueryRow(ctx, state, id)
	var count int
	err := res.Scan(&count)

	if errors.Is(err, pgx.ErrNoRows) {
		state := `INSERT INTO files (fileId, filename, count) VALUES ($1, $2, $3) RETURNING fileId;`
		res := s.db.QueryRow(ctx, state, id, fileName, 1)

		var id string
		err = res.Scan(&id)

		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}
		return id, nil
	}

	state = `UPDATE files SET count = $1 WHERE fileId = $2;`
	_, err = s.db.Exec(ctx, state, count+1, id)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return "file saved", nil
}
