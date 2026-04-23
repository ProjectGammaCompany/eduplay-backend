package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) DeleteUserData(ctx context.Context, userId string) (string, error) {
	const op = "storage.postgres.DeleteUserData"

	err := s.DeleteUserAnswers(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = s.DeleteUserComplaints(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = s.DeleteUserRatings(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = s.DeleteUserGroups(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = s.DeleteUserLinks(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = s.DeleteUserFavorites(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = s.DeleteUserEvents(ctx, userId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return "user " + userId + " data deleted", nil
}

func (s *Storage) DeleteUserAnswers(ctx context.Context, userId string) error {
	const op = "storage.postgres.DeleteUserAnswers"

	state := `DELETE FROM answers WHERE userId=$1;`

	_, err := s.db.Exec(ctx, state, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteUserEvents(ctx context.Context, userId string) error {
	const op = "storage.postgres.DeleteUserEvents"

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	state := `DELETE FROM options WHERE taskId IN 
	(SELECT taskId FROM tasks WHERE blockId IN 
	(SELECT blockId FROM blocks WHERE eventId IN (SELECT eventId FROM events WHERE ownerId=$1)));`
	_, err = tx.Exec(ctx, state, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	state = `DELETE FROM tasks WHERE blockId IN 
	(SELECT blockId FROM blocks WHERE eventId IN (SELECT eventId FROM events WHERE ownerId=$1));`
	_, err = tx.Exec(ctx, state, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	state = `DELETE FROM conditions WHERE prevBlockId IN 
	(SELECT blockId FROM blocks WHERE eventId IN (SELECT eventId FROM events WHERE ownerId=$1));`
	_, err = tx.Exec(ctx, state, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	state = `DELETE FROM blocks WHERE eventId IN (SELECT eventId FROM events WHERE ownerId=$1);`
	_, err = tx.Exec(ctx, state, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	state = `DELETE FROM groups WHERE eventId IN (SELECT eventId FROM events WHERE ownerId=$1);`
	_, err = tx.Exec(ctx, state, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	state = `DELETE FROM ratings WHERE eventId IN (SELECT eventId FROM events WHERE ownerId=$1);`
	_, err = tx.Exec(ctx, state, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	state = `DELETE FROM userFavorites WHERE eventId IN (SELECT eventId FROM events WHERE ownerId=$1);`
	_, err = tx.Exec(ctx, state, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	state = `DELETE FROM complaints WHERE eventId IN (SELECT eventId FROM events WHERE ownerId=$1);`
	_, err = tx.Exec(ctx, state, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	state = `DELETE FROM events WHERE ownerId=$1;`
	_, err = tx.Exec(ctx, state, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return tx.Commit(ctx)
}

func (s *Storage) DeleteUserLinks(ctx context.Context, userId string) error {
	const op = "storage.postgres.DeleteUserLinks"

	state := `DELETE FROM userLinks WHERE userId=$1;`

	_, err := s.db.Exec(ctx, state, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteUserFavorites(ctx context.Context, userId string) error {
	const op = "storage.postgres.DeleteUserFavorites"

	state := `DELETE FROM userFavorites WHERE userId=$1;`

	_, err := s.db.Exec(ctx, state, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteUserComplaints(ctx context.Context, userId string) error {
	const op = "storage.postgres.DeleteUserComplaints"

	state := `DELETE FROM complaints WHERE userId=$1;`

	_, err := s.db.Exec(ctx, state, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteUserRatings(ctx context.Context, userId string) error {
	const op = "storage.postgres.DeleteUserRatings"

	state := `DELETE FROM ratings WHERE userId=$1;`

	_, err := s.db.Exec(ctx, state, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteUserGroups(ctx context.Context, userId string) error {
	const op = "storage.postgres.DeleteUserGroups"

	state := `DELETE FROM userGroups WHERE userId=$1;`

	_, err := s.db.Exec(ctx, state, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
