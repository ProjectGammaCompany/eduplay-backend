package postgres

import (
	"context"
	"database/sql"
	dto "eduplay-event/internal/generated"
	"strings"
	"time"

	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"

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

func (s *Storage) PostEvent(ctx context.Context, in *dto.PostEventIn) (string, error) {
	const op = "storage.postgres.PostEvent"

	var (
		startDate *timestamppb.Timestamp
		endDate   *timestamppb.Timestamp
	)

	state := `INSERT INTO events (title, description, tags, cover, startDate, endDate, private, password, ownerId, lastEditionDate) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING eventId;`

	if in.StartDate != nil {
		startDate = in.StartDate
	}

	if in.EndDate != nil {
		endDate = in.EndDate
	}

	fmt.Println(endDate)

	res := s.db.QueryRow(ctx, state, in.Title, in.Description, in.Tags, in.Cover, startDate.AsTime(), endDate.AsTime(), in.Private, in.Password, in.OwnerId, time.Now())

	var id string
	err := res.Scan(&id)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetEvent(ctx context.Context, id string) (*dto.PostEventIn, error) {
	const op = "storage.postgres.GetEvent"

	state := `SELECT title, description, tags, cover, startDate, endDate, private, password, ownerId, lastEditionDate FROM events WHERE eventId = $1;`

	res := s.db.QueryRow(ctx, state, id)

	var (
		title           string
		description     string
		tags            []string
		cover           string
		startDate       time.Time
		endDate         time.Time
		private         bool
		password        string
		ownerId         string
		lastEditionDate time.Time
	)

	err := res.Scan(&title, &description, &tags, &cover, &startDate, &endDate, &private, &password, &ownerId, &lastEditionDate)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.PostEventIn{
		Title:           title,
		Description:     description,
		Tags:            tags,
		Cover:           cover,
		StartDate:       timestamppb.New(startDate),
		EndDate:         timestamppb.New(endDate),
		Private:         private,
		Password:        password,
		OwnerId:         ownerId,
		LastEditionDate: timestamppb.New(lastEditionDate),
	}, nil
}

func (s *Storage) GetRole(ctx context.Context, userId string, eventId string) (int64, error) {
	const op = "storage.postgres.GetRole"

	state := `SELECT isParticipant FROM userLinks WHERE userId = $1 AND eventId = $2;`

	res := s.db.QueryRow(ctx, state, userId, eventId)

	var role int64
	err := res.Scan(&role)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			state = `SELECT ownerId FROM events WHERE eventId = $1;`
			res = s.db.QueryRow(ctx, state, eventId)

			var ownerId string
			err = res.Scan(&ownerId)

			if err != nil {
				return 0, fmt.Errorf("%s: %w", op, err)
			}

			if userId == ownerId {
				return 1, nil
			}
			return 0, nil
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}

func (s *Storage) GetGroups(ctx context.Context, eventId string) (*dto.GetGroupsOut, error) {
	const op = "storage.postgres.GetGroups"

	state := `SELECT login, password FROM groups WHERE eventId = $1;`

	res, err := s.db.Query(ctx, state, eventId)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var groups []*dto.Group

	for res.Next() {
		var group dto.Group
		err = res.Scan(&group.Login, &group.Password)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		groups = append(groups, &group)
	}

	return &dto.GetGroupsOut{Groups: groups}, nil
}

func (s *Storage) GetCollaborators(ctx context.Context, eventId string) (*dto.GetCollaboratorsOut, error) {
	const op = "storage.postgres.GetCollaborators"

	state := `SELECT u.userId, u.email, u.avatar FROM users AS u LEFT JOIN userLinks AS ul ON u.userId = ul.userId WHERE ul.eventId = $1 AND ul.isParticipant = false;`

	res, err := s.db.Query(ctx, state, eventId)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var collaborators []*dto.User

	for res.Next() {
		var collaborator dto.User
		err = res.Scan(&collaborator.Id, &collaborator.Email, &collaborator.Avatar)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		collaborators = append(collaborators, &collaborator)
	}

	return &dto.GetCollaboratorsOut{Users: collaborators}, nil
}

func (s *Storage) PostEventBlock(ctx context.Context, in *dto.PostEventBlockIn) (string, error) {
	const op = "storage.postgres.PostEventBlock"

	state := `INSERT INTO blocks (eventId, name, blockOrder, isParallel) VALUES ($1, $2, $3, $4) RETURNING blockId;`

	res := s.db.QueryRow(ctx, state, in.EventId, in.Name, in.Order, in.IsParallel)

	var id string
	err := res.Scan(&id)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetEventBlocks(ctx context.Context, eventId string) (*dto.GetEventBlocksOut, error) {
	const op = "storage.postgres.GetEventBlocks"

	state := `SELECT blockId, name, blockOrder, isParallel FROM blocks WHERE eventId = $1;`

	res, err := s.db.Query(ctx, state, eventId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var blocks []*dto.BlockInfo

	for res.Next() {
		var block dto.BlockInfo
		err = res.Scan(&block.BlockId, &block.Name, &block.Order, &block.IsParallel)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		blocks = append(blocks, &block)

		conditions, err := s.GetBlockConditions(ctx, block.BlockId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		block.Conditions = conditions
	}

	return &dto.GetEventBlocksOut{Blocks: blocks}, nil
}

func (s *Storage) GetBlockConditions(ctx context.Context, blockId string) ([]*dto.Condition, error) {
	const op = "storage.postgres.GetBlockConditions"

	state := `SELECT prevBlockId, nextBlockId, group, min, max FROM conditions WHERE prevBlockId = $1 OR nextBlockId = $1;`

	res, err := s.db.Query(ctx, state, blockId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var conditions []*dto.Condition

	for res.Next() {
		var condition dto.Condition
		var nextBlockId sql.NullString
		err = res.Scan(&condition.PreviousBlockId, &nextBlockId, &condition.GroupIds, &condition.Min, &condition.Max)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if nextBlockId.Valid {
			condition.NextBlockId = nextBlockId.String
		}

		conditions = append(conditions, &condition)
	}

	return conditions, nil
}
