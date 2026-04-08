package postgres

import (
	"context"
	dto "eduplay-notification/internal/generated"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"

	_ "github.com/lib/pq"
)

type Storage struct {
	db      *pgxpool.Pool
	eventDb *pgxpool.Pool
	userDb  *pgxpool.Pool
}

func New(ctx context.Context, storagePath string, eventStoragePath string, userStoragePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	poolConfig, err := pgxpool.ParseConfig(storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s - %s", op, err)
	}

	eventPoolConfig, err := pgxpool.ParseConfig(eventStoragePath)
	if err != nil {
		return nil, fmt.Errorf("%s - %s", op, err)
	}

	userPoolConfig, err := pgxpool.ParseConfig(userStoragePath)
	if err != nil {
		return nil, fmt.Errorf("%s - %s", op, err)
	}

	poolConfig.MaxConns = 13
	poolConfig.MinConns = 5
	db, err := pgxpool.NewWithConfig(ctx, poolConfig)

	if err != nil {
		return nil, fmt.Errorf("%s - %s", op, err)
	}

	eventPoolConfig.MaxConns = 13
	eventPoolConfig.MinConns = 5
	eventDb, err := pgxpool.NewWithConfig(ctx, eventPoolConfig)

	if err != nil {
		return nil, fmt.Errorf("%s - %s", op, err)
	}

	userPoolConfig.MaxConns = 13
	userPoolConfig.MinConns = 5
	userDb, err := pgxpool.NewWithConfig(ctx, userPoolConfig)

	if err != nil {
		return nil, fmt.Errorf("%s - %s", op, err)
	}

	return &Storage{db: db, eventDb: eventDb, userDb: userDb}, nil
}

func (s *Storage) Stop(ctx context.Context) error {
	s.db.Close()
	return nil
}

func (s *Storage) GetNotifications(ctx context.Context, in *dto.Filters) (*dto.NotificationInfos, error) {
	const op = "storage.postgres.GetNotifications"

	notifications := make([]*dto.NotificationInfo, 0)

	state := `SELECT 
	notifId, 
	notifType, 
	notifDate, 
	userId, 
	eventId, 
	timeLeft, 
	eventName, 
	notStartedFavorite, 
	isRead
	FROM notifications WHERE userId = $1 
	ORDER BY notifDate DESC
	LIMIT $2 OFFSET $3;`

	res, err := s.db.Query(ctx, state, in.UserId, in.MaxOnPage, (in.Page-1)*in.MaxOnPage)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	for res.Next() {
		var notif dto.NotificationInfo
		var date time.Time
		err = res.Scan(&notif.NotificationId, &notif.Type, &date, &notif.UserId, &notif.EventId, &notif.TimeLeft, &notif.EventName, &notif.NotStartedFavorite, &notif.IsRead)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		notif.Date = timestamppb.New(date)
		notifications = append(notifications, &notif)
	}

	return &dto.NotificationInfos{Notifications: notifications}, nil
}

func (s *Storage) DeleteNotification(ctx context.Context, in *dto.Ids) error {
	const op = "storage.postgres.DeleteNotification"

	state := `DELETE FROM notifications WHERE userId = $1 AND notifId = $2;`

	_, err := s.db.Exec(ctx, state, in.UserId, in.NotificationId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil

}

func (s *Storage) AddNewNotification(ctx context.Context, in *dto.NotificationInfo) error {
	const op = "storage.postgres.AddNewNotification"

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM notifications WHERE userId = $1 AND eventId = $2 AND notifType = $3)`

	err := s.db.QueryRow(ctx, checkQuery, in.UserId, in.EventId, in.Type).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: check failed: %w", op, err)
	}

	if exists {
		return nil
	}

	insertQuery := `INSERT INTO notifications 
        (notifType, notifDate, userId, eventId, timeLeft, eventName, notStartedFavorite, isRead) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

	_, err = s.db.Exec(ctx, insertQuery,
		in.Type,
		in.Date.AsTime(),
		in.UserId,
		in.EventId,
		in.TimeLeft,
		in.EventName,
		in.NotStartedFavorite,
		in.IsRead)

	if err != nil {
		return fmt.Errorf("%s: insert failed: %w", op, err)
	}

	return nil
}

func (s *Storage) GetUserFavoriteStart(ctx context.Context, in *dto.Filters) error {
	const op = "storage.postgres.GetUserFavoriteStart"

	timeNow := time.Now().UTC().Add(3 * time.Hour)

	state := `SELECT e.eventId, e.title FROM events e JOIN userFavorites f ON e.eventId = f.eventId WHERE f.userId = $1 AND e.startDate < $2;`

	res, err := s.eventDb.Query(ctx, state, in.UserId, timeNow)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	for res.Next() {
		notif := &dto.NotificationInfo{UserId: in.UserId, Date: timestamppb.New(timeNow), Type: "favoriteEventStart"}
		err = res.Scan(&notif.EventId, &notif.EventName)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		err = s.AddNewNotification(ctx, notif)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *Storage) GetEndedEvents(ctx context.Context, in *dto.Filters) error {
	const op = "storage.postgres.GetEndedEvents"

	timeNowDay := time.Now().UTC().Add(3*time.Hour - time.Hour*24)
	timeNowHour := time.Now().UTC().Add(3*time.Hour - time.Hour)

	// TODO get not started favorite events hour

	state := `SELECT e.eventId, e.title FROM events e 
	JOIN userFavorites f ON e.eventId = f.eventId WHERE f.userId = $1 AND e.endDate < $2;`

	res, err := s.eventDb.Query(ctx, state, in.UserId, timeNowHour)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	for res.Next() {
		notif := &dto.NotificationInfo{UserId: in.UserId, Date: timestamppb.New(timeNowHour), Type: "eventEnd", NotStartedFavorite: true, TimeLeft: "hour"}
		err = res.Scan(&notif.EventId, &notif.EventName)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		state = `SELECT isParticipant FROM userLinks WHERE userId = $1 AND eventId = $2;`

		res2 := s.eventDb.QueryRow(ctx, state, in.UserId, notif.EventId)

		var isParticipant bool
		err = res2.Scan(&isParticipant)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				isParticipant = false
			}
		}

		if !isParticipant {
			err = s.AddNewNotification(ctx, notif)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}
	}

	// TODO get started events hour

	state = `SELECT e.eventId, e.title FROM events e JOIN 
	userLinks ul ON e.eventId = ul.eventId 
	WHERE ul.userId = $1 AND ul.finished = false AND ul.isParticipant = true AND e.endDate < $2;`

	res, err = s.eventDb.Query(ctx, state, in.UserId, timeNowHour)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	for res.Next() {
		notif := &dto.NotificationInfo{UserId: in.UserId, Date: timestamppb.New(timeNowHour), Type: "eventEnd", NotStartedFavorite: false, TimeLeft: "hour"}
		err = res.Scan(&notif.EventId, &notif.EventName)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		err = s.AddNewNotification(ctx, notif)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	// TODO get ended favorite events day

	state = `SELECT e.eventId, e.title FROM events e 
	JOIN userFavorites f ON e.eventId = f.eventId WHERE f.userId = $1 AND e.endDate < $2;`

	res, err = s.eventDb.Query(ctx, state, in.UserId, timeNowDay)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	for res.Next() {
		notif := &dto.NotificationInfo{UserId: in.UserId, Date: timestamppb.New(timeNowDay), Type: "eventEnd", NotStartedFavorite: true, TimeLeft: "day"}
		err = res.Scan(&notif.EventId, &notif.EventName)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		state = `SELECT isParticipant FROM userLinks WHERE userId = $1 AND eventId = $2;`

		res2 := s.eventDb.QueryRow(ctx, state, in.UserId, notif.EventId)

		var isParticipant bool
		err = res2.Scan(&isParticipant)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				isParticipant = false
			}
		}

		if !isParticipant {
			err = s.AddNewNotification(ctx, notif)
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}
	}

	// TODO get started events day

	state = `SELECT e.eventId, e.title FROM events e JOIN 
	userLinks ul ON e.eventId = ul.eventId 
	WHERE ul.userId = $1 AND ul.finished = false AND ul.isParticipant = true AND e.endDate < $2;`

	res, err = s.eventDb.Query(ctx, state, in.UserId, timeNowDay)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	for res.Next() {
		notif := &dto.NotificationInfo{UserId: in.UserId, Date: timestamppb.New(timeNowDay), Type: "eventEnd", NotStartedFavorite: false, TimeLeft: "day"}
		err = res.Scan(&notif.EventId, &notif.EventName)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		err = s.AddNewNotification(ctx, notif)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *Storage) GetUserNotifications(ctx context.Context) error {
	const op = "storage.postgres.GetUserNotifications"

	state := `SELECT userId FROM users;`

	res, err := s.userDb.Query(ctx, state)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer res.Close()

	for res.Next() {
		var userId string
		err = res.Scan(&userId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = s.GetUserFavoriteStart(ctx, &dto.Filters{UserId: userId})
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		err = s.GetEndedEvents(ctx, &dto.Filters{UserId: userId})
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
