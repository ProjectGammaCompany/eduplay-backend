package postgres

import (
	"context"
	"database/sql"
	dto "eduplay-event/internal/generated"
	"strings"
	"time"

	"errors"
	"fmt"

	"github.com/google/uuid"
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

func (s *Storage) SaveFile(ctx context.Context, fileName string, fileKey string) (string, error) {
	const op = "storage.postgres.SaveFile"

	state := `INSERT INTO files (fileKey, filename) VALUES ($1, $2) RETURNING fileId;`
	res := s.db.QueryRow(ctx, state, fileKey, fileName)

	var id string
	err := res.Scan(&id)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil

	// var id = strings.Split(fileUUID, ".")[0]

	// state := `SELECT count FROM files WHERE fileKey = $1`

	// res := s.db.QueryRow(ctx, state, fileKey)
	// var count int
	// err := res.Scan(&count)

	// if errors.Is(err, pgx.ErrNoRows) {
	// 	state := `INSERT INTO files (fileKey, filename, count) VALUES ($1, $2, $3) RETURNING fileId;`
	// 	res := s.db.QueryRow(ctx, state, fileKey, fileName, 1)

	// 	var id string
	// 	err = res.Scan(&id)

	// 	if err != nil {
	// 		return "", fmt.Errorf("%s: %w", op, err)
	// 	}
	// 	return id, nil
	// }

	// state = `UPDATE files SET count = $1 WHERE fileKey = $2;`
	// _, err = s.db.Exec(ctx, state, count+1, fileKey)

	// if err != nil {
	// 	return "", fmt.Errorf("%s: %w", op, err)
	// }

	// return "file saved", nil
}

func (s *Storage) PostEvent(ctx context.Context, in *dto.PostEventIn) (string, error) {
	const op = "storage.postgres.PostEvent"

	var (
		startDate *timestamppb.Timestamp
		endDate   *timestamppb.Timestamp
	)

	state := `INSERT INTO events (title, description, tags, cover, startDate, endDate, private, password, ownerId, lastEditionDate, allowDownloading, groupEvent) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING eventId;`

	if in.StartDate != nil {
		startDate = in.StartDate
	}

	if in.EndDate != nil {
		endDate = in.EndDate
	}

	fmt.Println(endDate)

	res := s.db.QueryRow(ctx, state, in.Title, in.Description, in.Tags, in.Cover, startDate.AsTime(), endDate.AsTime(), in.Private, in.Password, in.OwnerId, time.Now(), in.AllowDownloading, in.GroupEvent)

	var id string
	err := res.Scan(&id)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetEvent(ctx context.Context, id string) (*dto.PostEventIn, error) {
	const op = "storage.postgres.GetEvent"

	state := `SELECT title, description, tags, cover, startDate, endDate, private, password, ownerId, lastEditionDate, allowDownloading, groupEvent FROM events WHERE eventId = $1;`

	res := s.db.QueryRow(ctx, state, id)

	var (
		title            string
		description      string
		tags             []string
		cover            string
		startDate        time.Time
		endDate          time.Time
		private          bool
		password         string
		ownerId          string
		lastEditionDate  time.Time
		allowDownloading bool
		groupEvent       bool
	)

	err := res.Scan(&title, &description, &tags, &cover, &startDate, &endDate, &private, &password, &ownerId, &lastEditionDate, &allowDownloading, &groupEvent)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.PostEventIn{
		EventId:          id,
		Title:            title,
		Description:      description,
		Tags:             tags,
		Cover:            cover,
		StartDate:        timestamppb.New(startDate),
		EndDate:          timestamppb.New(endDate),
		Private:          private,
		Password:         password,
		OwnerId:          ownerId,
		LastEditionDate:  timestamppb.New(lastEditionDate),
		AllowDownloading: allowDownloading,
		GroupEvent:       groupEvent,
	}, nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventId string) (string, error) {
	const op = "storage.postgres.DeleteEvent"

	state := `DELETE FROM events WHERE eventId = $1;`

	_, err := s.db.Exec(ctx, state, eventId)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return "event" + eventId + "deleted", nil
}

func (s *Storage) GetRole(ctx context.Context, userId string, eventId string) (int64, error) {
	const op = "storage.postgres.GetRole"

	state := `SELECT isParticipant FROM userLinks WHERE userId = $1 AND eventId = $2;`

	res := s.db.QueryRow(ctx, state, userId, eventId)

	var role bool
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

	if role {
		return 0, nil
	}
	return 1, nil
}

func (s *Storage) GetGroups(ctx context.Context, eventId string) (*dto.GetGroupsOut, error) {
	const op = "storage.postgres.GetGroups"

	state := `SELECT groupId, login, password FROM groups WHERE eventId = $1;`

	res, err := s.db.Query(ctx, state, eventId)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	var groups []*dto.Group

	for res.Next() {
		var group dto.Group
		err = res.Scan(&group.Id, &group.Login, &group.Password)
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

	defer res.Close()

	// var collaborators []*dto.User
	collaborators := make([]*dto.User, 0)

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

func (s *Storage) DeleteEventBlock(ctx context.Context, blockId string) (string, error) {
	const op = "storage.postgres.DeleteEventBlock"

	state := `WITH deleted_block AS (
    DELETE FROM blocks
    WHERE blockId = $1
    RETURNING eventId, blockOrder
)
UPDATE blocks b
SET blockOrder = b.blockOrder - 1
FROM deleted_block d
WHERE b.eventId = d.eventId
  AND b.blockOrder > d.blockOrder;
`

	_, err := s.db.Exec(ctx, state, blockId)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return "block " + blockId + " deleted", nil
}

func (s *Storage) GetEventBlocks(ctx context.Context, eventId string) (*dto.GetEventBlocksOut, error) {
	const op = "storage.postgres.GetEventBlocks"

	state := `SELECT blockId, name, blockOrder, isParallel FROM blocks WHERE eventId = $1 ORDER BY blockOrder;`

	res, err := s.db.Query(ctx, state, eventId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

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

	state := `SELECT conditionId, prevBlockId, nextBlockId, groupName, min, max FROM conditions WHERE prevBlockId = $1 OR nextBlockId = $1;`

	res, err := s.db.Query(ctx, state, blockId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	var conditions []*dto.Condition

	for res.Next() {
		var condition dto.Condition
		var nextBlockId sql.NullString
		err = res.Scan(&condition.ConditionId, &condition.PreviousBlockId, &nextBlockId, &condition.GroupIds, &condition.Min, &condition.Max)
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

func (s *Storage) GetPublicEvent(ctx context.Context, ids *dto.UserEventIds) (*dto.GetPublicEvent, error) {
	const op = "storage.postgres.GetPublicEvent"

	state := `SELECT 
    e.eventId,
    e.title, 
    e.description,
    e.cover,
    e.lastEditionDate,
    e.tags,
    COALESCE(AVG(r.rating), 0) AS rate,
    EXISTS (
        SELECT 1 
        FROM userFavorites 
        WHERE userId = $2 
          AND eventId = e.eventId
    ) AS favorite
FROM events e
LEFT JOIN ratings r ON e.eventId = r.eventId  
WHERE e.private = false
  AND e.eventId = $1
GROUP BY e.eventId;`

	res := s.db.QueryRow(ctx, state, ids.EventId, ids.UserId)

	tags := make([]string, 0)
	lastEditionDate := time.Time{}
	// var event dto.GetPublicEvent
	event := dto.GetPublicEvent{}
	err := res.Scan(&event.EventId, &event.Title, &event.Description, &event.Cover, &lastEditionDate, &tags, &event.Rate, &event.Favorite)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	event.LastEditionDate = timestamppb.New(lastEditionDate)
	_, _, finished, _, err := s.GetEventProgress(ctx, ids.UserId, ids.EventId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		event.Status = "notStarted"
	}

	if finished {
		event.Status = "finished"
	} else {
		event.Status = "inProgress"
	}

	fullTags, err := s.GetTagsByIds(ctx, tags)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	event.Tags = fullTags.Tags

	return &event, nil
}

func (s *Storage) GetTagsByIds(ctx context.Context, ids []string) (*dto.Tags, error) {
	const op = "storage.postgres.GetTagsByIds"

	state := `SELECT tagId, name FROM tags WHERE tagId = ANY($1);`

	res, err := s.db.Query(ctx, state, ids)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	tags := make([]*dto.Tag, 0)
	for res.Next() {
		tag := &dto.Tag{}
		err = res.Scan(&tag.Id, &tag.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		tags = append(tags, tag)
	}

	return &dto.Tags{Tags: tags}, nil
}

func (s *Storage) GetEventProgress(ctx context.Context, userId string, eventId string) (currTaskId string, currBlockId string, finished bool, currTaskStartTime time.Time, err error) {
	const op = "storage.postgres.GetEventProgress"
	currTaskIdNull := sql.NullString{}
	currBlockIdNull := sql.NullString{}

	state := `SELECT currTaskId, currBlockId, finished, currTaskStartTime FROM userLinks WHERE userId = $1 AND eventId = $2 AND isParticipant = true;`

	err = s.db.QueryRow(ctx, state, userId, eventId).Scan(&currTaskIdNull, &currBlockIdNull, &finished, &currTaskStartTime)
	if err != nil {
		return "", "", false, time.Time{}, fmt.Errorf("%s: %w", op, err)
	}

	if currTaskIdNull.Valid {
		currTaskId = currTaskIdNull.String
	}

	if currBlockIdNull.Valid {
		currBlockId = currBlockIdNull.String
	}

	return
}

func (s *Storage) GetPublicEvents(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error) {
	const op = "storage.postgres.GetPublicEvents"

	var args = []interface{}{}

	limit := in.MaxOnPage
	offset := (in.Page - 1) * in.MaxOnPage

	args = append(args, limit, offset, in.UserId)
	userParamIdx := 3

	tagsParamIdx := 4
	if len(in.Tags) > 0 {
		args = append(args, in.Tags)
	}

	titleParamIdx := 5
	if in.Title != "" {
		args = append(args, in.Title)
	}

	var where []string
	where = append(where, "e.private = false")

	if in.Active {
		where = append(where, "e.startDate < now() AND e.endDate > now()")
	}

	if in.Favorites {
		where = append(where, fmt.Sprintf("EXISTS (SELECT 1 FROM userFavorites WHERE userId = $%d AND eventId = e.eventId)", userParamIdx))
	}

	if len(in.Tags) > 0 {
		where = append(where, fmt.Sprintf("e.tags && $%d", tagsParamIdx))
	}

	if in.Title != "" {
		where = append(where, fmt.Sprintf("e.title ILIKE $%d", titleParamIdx))
	}

	whereClause := strings.Join(where, " AND ")

	orderBy := "e.lastEditionDate DESC"
	if in.DecliningRating {
		orderBy = "rate DESC"
	}

	state := fmt.Sprintf(`
			SELECT
            e.eventId,
            e.title,
            e.description,
            e.cover,
            e.lastEditionDate,
            e.tags,
            COALESCE(AVG(r.rating), 0) AS rate,
            EXISTS (SELECT 1 FROM userFavorites WHERE userId = $%d AND eventId = e.eventId) AS favorite
        FROM events e
        LEFT JOIN ratings r ON e.eventId = r.eventId
        WHERE %s
        GROUP BY e.eventId
        ORDER BY %s
        LIMIT $1 OFFSET $2
	`, userParamIdx, whereClause, orderBy)

	// fmt.Println(state, args)

	res, err := s.db.Query(ctx, state, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	events := &dto.GetPublicEventsOut{}

	for res.Next() {
		tags := make([]string, 0)
		var event dto.GetPublicEvent
		var lastEditionDate time.Time
		err = res.Scan(&event.EventId, &event.Title, &event.Description, &event.Cover, &lastEditionDate, &tags, &event.Rate, &event.Favorite)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		fullTags, err := s.GetTagsByIds(ctx, tags)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		event.Tags = fullTags.Tags
		event.LastEditionDate = timestamppb.New(lastEditionDate)
		events.Events = append(events.Events, &event)
	}

	return events, nil

	// 	state := `SELECT
	//     e.eventId as id,
	//     e.title,
	//     e.description,
	//     e.cover,
	//     e.lastEditionDate,
	//     COALESCE(AVG(r.rating), 0) as rate,
	//     -- Проверяем, есть ли событие в избранном у пользователя
	//    	CASE
	//             WHEN $6::uuid IS NULL THEN false
	//             ELSE EXISTS(
	//                 SELECT 1 FROM userFavorites uf
	//                 WHERE uf.eventId = e.eventId AND uf.userId = $6::uuid
	//             )
	//         END as favorite,
	//     e.tags
	// FROM events e
	// LEFT JOIN ratings r ON e.eventId = r.eventId
	// WHERE
	//     e.private = false
	//     AND (
	//         $4 = false
	//         OR (
	//             e.startDate < NOW()
	//             AND (e.endDate IS NULL OR e.endDate > NOW())
	//         )
	//     )
	//     AND (
	//         $5::text[] IS NULL
	//         OR cardinality($5::text[]) = 0
	//         OR e.tags && $5::text[]
	//     )
	// GROUP BY e.eventId
	// ORDER BY
	//     CASE
	//         WHEN $3 = true THEN -COALESCE(AVG(r.rating), 0)
	//         ELSE 1
	//     END,
	//     e.lastEditionDate DESC
	// LIMIT $2
	// OFFSET (($1 - 1) * $2);`

	// 	emptyTags := []string{}

	// 	res, err := s.db.Query(ctx, state, in.Page, in.MaxOnPage, in.DecliningRating, in.Active, emptyTags, in.UserId)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("%s: %w", op, err)
	// 	}
	// 	fmt.Println(state, in.Page, in.MaxOnPage, in.DecliningRating, in.Active, emptyTags, in.UserId)
	// 	events := &dto.GetPublicEventsOut{}

	// 	for res.Next() {
	// 		fmt.Println("trying to get row")
	// 		var event *dto.GetPublicEvent
	// 		var lastEd time.Time
	// 		err = res.Scan(&event.EventId, &event.Title, &event.Description, &event.Cover, &lastEd, &event.Rate, &event.Favorite, &event.Tags)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("%s: %w", op, err)
	// 		}
	// 		fmt.Println("got one row")
	// 		event.LastEditionDate = timestamppb.New(lastEd)
	// 		events.Events = append(events.Events, event)
	// 	}

	// 	return events, nil
}

func (s *Storage) GetUserFavorites(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error) {
	const op = "storage.postgres.GetUserFavorites"

	state := `SELECT 
    e.eventId,
    e.title, 
    e.description,
    e.cover,
    e.lastEditionDate,
    e.tags,
    COALESCE(AVG(r.rating), 0) as rate,
    true as favorite
FROM events AS e
LEFT JOIN ratings r ON e.eventId = r.eventId
WHERE EXISTS (
    SELECT 1 
    FROM userFavorites 
    WHERE userId = $3 AND eventId = e.eventId
)
GROUP BY e.eventId
ORDER BY e.lastEditionDate DESC
LIMIT $1 
OFFSET $2;`

	res, err := s.db.Query(ctx, state, in.MaxOnPage, (in.Page-1)*in.MaxOnPage, in.UserId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	events := &dto.GetPublicEventsOut{}

	for res.Next() {
		tags := make([]string, 0)
		var event dto.GetPublicEvent
		var lastEd time.Time
		err = res.Scan(&event.EventId, &event.Title, &event.Description, &event.Cover, &lastEd, &tags, &event.Rate, &event.Favorite)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		fullTags, err := s.GetTagsByIds(ctx, tags)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		event.Tags = fullTags.Tags
		event.LastEditionDate = timestamppb.New(lastEd)
		events.Events = append(events.Events, &event)
	}

	return events, nil
}

func (s *Storage) GetOwnedEvents(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error) {
	const op = "storage.postgres.GetOwnedEvents"

	state := `SELECT 
	e.eventId,
    e.title, 
    e.description,
    e.cover,
    e.lastEditionDate,
    e.tags,
    COALESCE(AVG(r.rating), 0) as rate,
    EXISTS (
        SELECT 1 
        FROM userFavorites 
        WHERE userId = $3 AND eventId = e.eventId
    ) as favorite
FROM events AS e
LEFT JOIN ratings r ON e.eventId = r.eventId
WHERE e.ownerId = $3
GROUP BY e.eventId
ORDER BY e.lastEditionDate DESC
LIMIT $1 
OFFSET $2;`

	res, err := s.db.Query(ctx, state, in.MaxOnPage, (in.Page-1)*in.MaxOnPage, in.UserId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	events := &dto.GetPublicEventsOut{}

	for res.Next() {
		tags := make([]string, 0)
		var event dto.GetPublicEvent
		var lastEd time.Time
		err = res.Scan(&event.EventId, &event.Title, &event.Description, &event.Cover, &lastEd, &tags, &event.Rate, &event.Favorite)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		fullTags, err := s.GetTagsByIds(ctx, tags)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		event.Tags = fullTags.Tags
		event.LastEditionDate = timestamppb.New(lastEd)
		events.Events = append(events.Events, &event)
	}

	return events, nil

}

func (s *Storage) GetHistory(ctx context.Context, in *dto.EventBaseFilters) (*dto.GetPublicEventsOut, error) {
	const op = "storage.postgres.GetHistory"

	state := `SELECT 
    e.eventId,
    e.title, 
    e.description,
    e.cover,
    e.lastEditionDate,
    e.tags,
    COALESCE(AVG(r.rating), 0) as rate,
    EXISTS (
        SELECT 1 
        FROM userFavorites 
        WHERE userId = $3 AND eventId = e.eventId
    ) as favorite,
    ul.finished as is_finished
FROM events AS e
INNER JOIN userLinks ul ON e.eventId = ul.eventId 
    AND ul.userId = $3 
    AND ul.finished = true  
LEFT JOIN ratings r ON e.eventId = r.eventId
GROUP BY e.eventId, ul.finished
ORDER BY e.lastEditionDate DESC
LIMIT $1 
OFFSET $2;`

	// С учетом статуса прохождения
	// 	SELECT
	//     e.eventId,
	//     e.title,
	//     e.description,
	//     e.cover,
	//     e.lastEditionDate,
	//     e.tags,
	//     COALESCE(AVG(r.rating), 0) as rate,
	//     EXISTS (
	//         SELECT 1
	//         FROM userFavorites
	//         WHERE userId = $3 AND eventId = e.eventId
	//     ) as favorite,
	//     COALESCE(ul.finished, false) as is_finished,
	//     -- Определяем статус участия
	//     CASE
	//         WHEN ul.finished = true THEN 'completed'
	//         WHEN ul.currBlockId IS NOT NULL THEN 'in_progress'
	//         WHEN ul.eventId IS NOT NULL THEN 'joined'
	//         ELSE 'not_joined'
	//     END as user_status
	// FROM events AS e
	// -- LEFT JOIN чтобы получить ВСЕ квесты, но со статусом пользователя
	// LEFT JOIN userLinks ul ON e.eventId = ul.eventId AND ul.userId = $3
	// LEFT JOIN ratings r ON e.eventId = r.eventId
	// -- Фильтр по статусу (если нужно)
	// -- WHERE ul.finished = true  -- для завершенных
	// -- WHERE ul.eventId IS NOT NULL  -- для всех, в которых пользователь участвовал
	// GROUP BY e.eventId, ul.finished, ul.currBlockId
	// ORDER BY
	//     -- Сортировка по статусу: завершенные -> в процессе -> остальные
	//     CASE
	//         WHEN ul.finished = true THEN 1
	//         WHEN ul.currBlockId IS NOT NULL THEN 2
	//         ELSE 3
	//     END,
	//     e.lastEditionDate DESC
	// LIMIT $1
	// OFFSET $2;

	res, err := s.db.Query(ctx, state, in.MaxOnPage, (in.Page-1)*in.MaxOnPage, in.UserId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	events := &dto.GetPublicEventsOut{}

	for res.Next() {
		tags := make([]string, 0)
		var event dto.GetPublicEvent
		var lastEd time.Time
		err = res.Scan(&event.EventId, &event.Title, &event.Description, &event.Cover, &lastEd, &tags, &event.Rate, &event.Favorite)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		fullTags, err := s.GetTagsByIds(ctx, tags)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		event.Tags = fullTags.Tags
		event.LastEditionDate = timestamppb.New(lastEd)
		events.Events = append(events.Events, &event)
	}

	return events, nil
}

func (s *Storage) PutFavorite(ctx context.Context, in *dto.PutFavoriteIn) (string, error) {
	const op = "storage.postgres.PutFavorite"

	if !in.Favorite {
		state := `DELETE FROM userFavorites WHERE userId = $1 AND eventId = $2;`
		_, err := s.db.Exec(ctx, state, in.UserId, in.EventId)
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}
		return "removed", nil
	}
	state := `SELECT * FROM userFavorites WHERE userId = $1 AND eventId = $2;`
	_, err := s.db.Exec(ctx, state, in.UserId, in.EventId)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	state = `INSERT INTO userFavorites (userId, eventId) VALUES ($1, $2);`
	_, err = s.db.Exec(ctx, state, in.UserId, in.EventId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return "added", nil
}

func (s *Storage) GetAllTags(ctx context.Context) (*dto.Tags, error) {
	const op = "storage.postgres.GetAllTags"

	state := `SELECT tagId, name FROM tags;`

	res, err := s.db.Query(ctx, state)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	tags := &dto.Tags{}

	for res.Next() {
		var tag dto.Tag
		err = res.Scan(&tag.Id, &tag.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		tags.Tags = append(tags.Tags, &tag)
	}

	return tags, nil
}

func (s *Storage) PostTask(ctx context.Context, in *dto.Task) (string, error) {
	const op = "storage.postgres.PostTask"

	var id string
	var order int64

	state := `SELECT COALESCE(MAX(taskOrder), 0) FROM tasks WHERE blockId = $1;`

	res := s.db.QueryRow(ctx, state, in.BlockId)

	err := res.Scan(&order)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	fmt.Println(order)
	order++

	state = `INSERT INTO tasks (blockId, name, description, type, files, time, points, partialPoint, taskOrder) VALUES ($1, $2, $3, $4, COALESCE($5, '{}'::text[]), $6, $7, $8, $9) RETURNING taskId;`

	res = s.db.QueryRow(ctx, state, in.BlockId, in.Name, in.Description, in.Type, in.Files, in.Time, in.Points, in.PartialPoints, order)

	err = res.Scan(&id)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	for _, option := range in.Options {
		state := `INSERT INTO options (taskId, value, isCorrect) VALUES ($1, $2, $3);`
		_, err := s.db.Exec(ctx, state, id, option.Value, option.IsCorrect)
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}
	}

	return id, nil
}

func (s *Storage) PostBlockCondition(ctx context.Context, in *dto.Condition) (*dto.PostConditionOut, error) {
	const op = "storage.postgres.PostBlockCondition"

	state := `INSERT INTO conditions (prevBlockId, nextBlockId, groupName, min, max) VALUES ($1, $2, COALESCE($3, '{}'::text[]), $4, $5) RETURNING conditionId;`

	res := s.db.QueryRow(ctx, state, in.PreviousBlockId, in.NextBlockId, in.GroupIds, in.Min, in.Max)

	var id string
	err := res.Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	state = `SELECT blockOrder FROM blocks WHERE blockId = $1;`
	res = s.db.QueryRow(ctx, state, in.NextBlockId)

	var order int64
	err = res.Scan(&order)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.PostConditionOut{ConditionId: id, BlockOrder: order}, nil
}

func (s *Storage) DeleteBlockCondition(ctx context.Context, conditionId string) (string, error) {
	const op = "storage.postgres.DeleteBlockCondition"

	state := `DELETE FROM conditions WHERE conditionId = $1;`

	_, err := s.db.Exec(ctx, state, conditionId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return "block condition " + conditionId + " deleted", nil
}

func (s *Storage) GetBlockInfo(ctx context.Context, blockId string) (*dto.PostEventBlockIn, error) {
	const op = "storage.postgres.GetBlockInfo"

	state := `SELECT eventId, name, blockOrder, isParallel, showPoints, showAnswers FROM blocks WHERE blockId = $1;`

	res := s.db.QueryRow(ctx, state, blockId)

	var info dto.PostEventBlockIn
	err := res.Scan(&info.EventId, &info.Name, &info.Order, &info.IsParallel, &info.ShowPoints, &info.ShowAnswers)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &info, nil
}

func (s *Storage) GetBlockConditionsFull(ctx context.Context, blockId string) (*dto.BlockInfo, error) {
	const op = "storage.postgres.GetBlockConditions"

	state := `SELECT bl.blockOrder, c.prevBlockId, c.nextBlockId, c.groupName, c.min, c.max, c.conditionId FROM blocks bl INNER JOIN conditions c ON bl.blockId = c.prevBlockId WHERE bl.blockId = $1;`

	res, err := s.db.Query(ctx, state, blockId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	conditions := make([]*dto.Condition, 0)
	for res.Next() {
		condition := &dto.Condition{}
		// var emptyGroup pgtype.Array[string]
		err = res.Scan(&condition.NextBlockOrder, &condition.PreviousBlockId, &condition.NextBlockId, &condition.GroupIds, &condition.Min, &condition.Max, &condition.ConditionId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		// if emptyGroup.Valid {
		// 	condition.GroupIds = emptyGroup.Elements
		// }
		conditions = append(conditions, condition)
	}

	return &dto.BlockInfo{BlockId: blockId, Conditions: conditions}, nil
}

func (s *Storage) GetBlockTasks(ctx context.Context, blockId string) (*dto.Tasks, error) {
	const op = "storage.postgres.GetBlockTasks"

	state := `SELECT taskId, name, description, type, files, time, points, partialPoint, taskOrder FROM tasks WHERE blockId = $1 ORDER BY taskOrder;`

	res, err := s.db.Query(ctx, state, blockId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	tasks := make([]*dto.Task, 0)
	for res.Next() {
		task := &dto.Task{}
		err = res.Scan(&task.TaskId, &task.Name, &task.Description, &task.Type, &task.Files, &task.Time, &task.Points, &task.PartialPoints, &task.Order)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		options, err := s.GetTaskOptions(ctx, task.TaskId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		task.Options = options.Options
		tasks = append(tasks, task)
	}

	return &dto.Tasks{Tasks: tasks}, nil
}

func (s *Storage) GetTaskOptions(ctx context.Context, taskId string) (*dto.TaskOptions, error) {
	const op = "storage.postgres.GetTaskOptions"

	state := `SELECT value, isCorrect, optionId FROM options WHERE taskId = $1;`

	res, err := s.db.Query(ctx, state, taskId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	options := make([]*dto.TaskOption, 0)
	for res.Next() {
		option := &dto.TaskOption{}
		err = res.Scan(&option.Value, &option.IsCorrect, &option.OptionId)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		options = append(options, option)
	}

	return &dto.TaskOptions{Options: options}, nil
}

func (s *Storage) GetTaskById(ctx context.Context, taskId string) (*dto.Task, error) {
	const op = "storage.postgres.GetTaskById"

	state := `SELECT taskId, name, description, type, files, time, points, partialPoint, taskOrder, blockId FROM tasks WHERE taskId = $1;`

	res := s.db.QueryRow(ctx, state, taskId)

	task := &dto.Task{}
	err := res.Scan(&task.TaskId, &task.Name, &task.Description, &task.Type, &task.Files, &task.Time, &task.Points, &task.PartialPoints, &task.Order, &task.BlockId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	options, err := s.GetTaskOptions(ctx, task.TaskId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	task.Options = options.Options

	return task, nil
}

func (s *Storage) DeleteTaskById(ctx context.Context, taskId string) (string, error) {
	const op = "storage.postgres.DeleteTaskById"

	state := `WITH deleted_task AS (
        DELETE FROM tasks 
        WHERE taskId = $1
        RETURNING blockId, taskOrder
    )
    UPDATE tasks t
    SET t.taskOrder = t.taskOrder - 1
    FROM deleted_task d
    WHERE t.blockId = d.blockId 
      AND t.taskOrder > d.taskOrder;`

	_, err := s.db.Exec(ctx, state, taskId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return "task " + taskId + " removed", nil
}

func (s *Storage) PostAnswer(ctx context.Context, answer *dto.Answer) (string, error) {
	const op = "storage.postgres.PostAnswer"

	state := `INSERT INTO answers (userId, taskId, values, points) VALUES ($1, $2, $3, $4) RETURNING answerId;`

	var answerId string
	err := s.db.QueryRow(ctx, state, answer.UserId, answer.TaskId, answer.Answer, answer.Points).Scan(&answerId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return answerId, nil
}

func (s *Storage) PutNextStage(ctx context.Context, stage *dto.EventBlockTaskUserIds) (string, error) {
	const op = "storage.postgres.PostNextStage"

	var (
		currTaskId  *uuid.UUID
		currBlockId *uuid.UUID
	)

	if stage.TaskId == "" {
		currTaskId = nil
	} else {
		taskId, err := uuid.Parse(stage.TaskId)
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}
		currTaskId = &taskId
	}

	if stage.BlockId == "" {
		currBlockId = nil
	} else {
		blockId, err := uuid.Parse(stage.BlockId)
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}
		currBlockId = &blockId
	}

	state := `UPDATE userLinks SET currTaskId = $1, currBlockId = $2, finished = $5 WHERE userId = $3 AND eventId = $4;`

	_, err := s.db.Exec(ctx, state, currTaskId, currBlockId, stage.UserId, stage.EventId, stage.Finished)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return "updated", nil
}

func (s *Storage) PutTimestamp(ctx context.Context, userId string, eventId string, timestamp *timestamppb.Timestamp) (string, error) {
	const op = "storage.postgres.PutTimestamp"

	state := `UPDATE userLinks SET currTaskStartTime = $3 WHERE userId = $1 AND eventId = $2;`

	_, err := s.db.Exec(ctx, state, userId, eventId, timestamp.AsTime())
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return "updated", nil
}

func (s *Storage) GetNextStage(ctx context.Context, stage *dto.UserEventIds) (linkId string, currTaskId string, currBlockId string, finished bool, startTime *timestamppb.Timestamp, err error) {
	const op = "storage.postgres.GetNextStage"

	var (
		currTaskStartTime time.Time
	)

	state := `SELECT linkId, COALESCE(currTaskId::text, ''), COALESCE(currBlockId::text, ''), finished, currTaskStartTime FROM userLinks WHERE userId = $1 AND eventId = $2 AND isParticipant = true;`

	fmt.Println("state", stage.EventId)
	err = s.db.QueryRow(ctx, state, stage.UserId, stage.EventId).Scan(&linkId, &currTaskId, &currBlockId, &finished, &currTaskStartTime)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			state = `INSERT INTO userLinks (userId, eventId, isParticipant, currTaskStartTime)
VALUES ($1, $2, true, '1970-01-01 00:00:00'::timestamp) 
RETURNING linkId, COALESCE(currTaskId::text, ''), COALESCE(currBlockId::text, ''), finished, currTaskStartTime;`

			err = s.db.QueryRow(ctx, state, stage.UserId, stage.EventId).Scan(&linkId, &currTaskId, &currBlockId, &finished, &currTaskStartTime)
			if err != nil {
				return "", "", "", false, nil, fmt.Errorf("%s: %w", op, err)
			}
		}
		return "", "", "", false, nil, fmt.Errorf("%s: %w", op, err)
	}

	startTime = timestamppb.New(currTaskStartTime)

	return
}

func (s *Storage) EndMe(ctx context.Context, userId string, eventId string) (string, error) {
	const op = "storage.postgres.EndMe"

	state := `UPDATE userLinks SET finished = true, currTaskStartTime = '1970-01-01 00:00:00'::timestamp, currTaskId = NULL, currBlockId = NULL WHERE userId = $1 AND eventId = $2;`

	_, err := s.db.Exec(ctx, state, userId, eventId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return "updated to finished", nil
}

func (s *Storage) GetUserBlockPointsSum(ctx context.Context, userId string, blockId string) (int64, error) {
	const op = "storage.postgres.GetUserBlockPointsSum"

	state := `
SELECT COALESCE(SUM(a.points), 0) AS total_points
FROM answers a
JOIN tasks t ON t.taskId = a.taskId
WHERE a.userId = $1
  AND t.blockId = $2;
`

	var total int64
	err := s.db.QueryRow(ctx, state, userId, blockId).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return total, nil
}

func (s *Storage) GetUserBlockTasksShort(ctx context.Context, blockId string, userId string) ([]*dto.NextStageTaskShort, error) {
	const op = "storage.postgres.GetUserBlockTasksShort"

	state := `SELECT 
    t.taskId,
    t.name,
    t.time,
    t.taskOrder,
    EXISTS (
        SELECT 1
        FROM answers a
        WHERE a.taskId = t.taskId
          AND a.userId = $2
    ) AS isCompleted,
	t.type, 
	t.description

FROM tasks t
WHERE t.blockId = $1
ORDER BY t.taskOrder;
`

	res, err := s.db.Query(ctx, state, blockId, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	tasks := make([]*dto.NextStageTaskShort, 0)
	for res.Next() {
		task := &dto.NextStageTaskShort{}
		err = res.Scan(&task.TaskId, &task.Name, &task.Time, &task.Order, &task.IsCompleted, &task.Type, &task.Description)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
