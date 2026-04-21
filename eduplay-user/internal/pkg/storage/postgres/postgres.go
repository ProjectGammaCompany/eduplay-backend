package postgres

import (
	"context"
	dto "eduplay-user/internal/generated"
	"eduplay-user/internal/model"
	"eduplay-user/internal/storage"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/lib/pq"

	pgx "github.com/jackc/pgx/v5"
)

type Storage struct {
	db *pgxpool.Pool
	// loc *time.Location
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

	// loc, err := time.LoadLocation("Europe/Moscow")
	// if err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	return &Storage{db: db}, nil
}

func (s *Storage) Stop(ctx context.Context) error {
	s.db.Close()
	return nil
}

func (s *Storage) SignUpUser(ctx context.Context, email string, password string) (string, error) {
	const op = "storage.postgres.SignUpUser"

	state := `SELECT userId FROM users WHERE email = $1`

	res := s.db.QueryRow(ctx, state, email)
	var id string
	err := res.Scan(&id)

	if errors.Is(err, pgx.ErrNoRows) {
		state := `INSERT INTO users (email, passwordHash, name) VALUES ($1, $2, $3) RETURNING userId;`

		userName := strings.Split(email, "@")[0]

		res := s.db.QueryRow(ctx, state, email, password, userName)

		var id string
		err = res.Scan(&id)

		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}
		return id, nil
	}

	return "", storage.ErrUserAlreadyExists
}

func (s *Storage) SaveSession(ctx context.Context, userId string, refreshToken string) error {
	const op = "storage.postgres.SaveSession"

	// authExpiry := time.Now().UTC().Add(3 * time.Hour).Add(15 * time.Minute)
	refreshExpiry := time.Now().UTC().Add(3 * time.Hour).Add(time.Hour)
	state := `INSERT INTO sessions (userId, refresh_token, refresh_expires, active) VALUES ($1, $2, $3, $4) RETURNING sessionId;`
	_, err := s.db.Exec(ctx, state, userId, refreshToken, refreshExpiry, true)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// state = `UPDATE users SET role = $1 WHERE userId = $2;`
	// _, err = s.db.Exec(ctx, state, "user", userId)

	// if err != nil {
	// 	return fmt.Errorf("%s: %w", op, err)
	// }

	// state = `INSERT INTO userSubscriptions (userId, subscriptionLevel, date, expiresAt, sessions) VALUES ($1, $2, $3, $4, $5);`
	// _, err = s.db.Exec(ctx, state, userId, 0, time.Now().UTC().Add(3*time.Hour), time.Now().UTC().Add(3*time.Hour), 0)
	// if err != nil {
	// 	return fmt.Errorf("%s: %w", op, err)
	// }

	return nil
}

func (s *Storage) UpdateSession(ctx context.Context, userId string, refreshToken string) error {
	const op = "storage.postgres.UpdateSession"

	// authExpiry := time.Now().UTC().Add(3 * time.Hour).Add(15 * time.Minute)
	refreshExpiry := time.Now().UTC().Add(3 * time.Hour).Add(180 * 24 * time.Hour)
	state := `UPDATE sessions SET refresh_token = $1, refresh_expires = $2, active = true WHERE userId = $3;`
	_, err := s.db.Exec(ctx, state, refreshToken, refreshExpiry, userId)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) SaveSubscription(ctx context.Context, userId string, subscriptionLevel int64, acquiringDate time.Time, subscriptionExpiry time.Time, sessionNums int64) error {
	const op = "storage.postgres.SaveSubscription"

	state := `INSERT INTO userSubscriptions (userId, subscriptionLevel, date, expiresAt, sessions) VALUES ($1, $2, $3, $4, $5);`
	_, err := s.db.Exec(ctx, state, userId, subscriptionLevel, acquiringDate, subscriptionExpiry, sessionNums)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	const op = "storage.postgres.GetUserByEmail"

	state := `SELECT userId, email, passwordHash FROM users WHERE email = $1`
	res := s.db.QueryRow(ctx, state, email)
	user := &model.User{}
	err := res.Scan(&user.Id, &user.Email, &user.Password)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) GetSessionByUserId(ctx context.Context, userId string) (*model.Session, error) {
	const op = "storage.postgres.GetSessionByUserId"

	state := `SELECT refresh_token, active FROM sessions WHERE userId = $1`
	res := s.db.QueryRow(ctx, state, userId)
	session := &model.Session{}
	err := res.Scan(&session.RefreshToken, &session.IsActive)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// state = `SELECT role FROM users WHERE userId = $1`
	// res = s.db.QueryRow(ctx, state, userId)
	// err = res.Scan(&session.Role)
	// if err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	// state = `SELECT subscriptionLevel FROM userSubscriptions WHERE userId = $1`
	// res = s.db.QueryRow(ctx, state, userId)
	// err = res.Scan(&session.AccessLevel)
	// if err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	return session, nil
}

func (s *Storage) GetUserByRefreshToken(ctx context.Context, refreshToken string) (*model.User, error) {
	const op = "storage.postgres.GetUserByAuthToken"

	state := `SELECT userId FROM sessions WHERE refresh_token = $1`
	res := s.db.QueryRow(ctx, state, refreshToken)
	var userId string
	err := res.Scan(&userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	state = `SELECT userId, email, passwordHash FROM users WHERE userId = $1`
	res = s.db.QueryRow(ctx, state, userId)
	user := &model.User{}
	err = res.Scan(&user.Id, &user.Email, &user.Password)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) CheckRefreshTokenExists(ctx context.Context, refreshToken string) (bool, time.Time, error) {
	const op = "storage.postgres.CheckRefreshTokenExists"
	state := `SELECT refresh_token, refresh_expires FROM sessions WHERE refresh_token = $1`
	res := s.db.QueryRow(ctx, state, refreshToken)
	// if err != nil {
	// 	if errors.Is(err, pgx.ErrNoRows) {
	// 		return false, time.Time{}, storage.ErrInvalidRefresh
	// 	}
	// 	return false, time.Time{}, fmt.Errorf("%s: %w", op, err)
	// }
	var token string
	var expiry time.Time
	err := res.Scan(&token, &expiry)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, time.Time{}, storage.ErrInvalidRefresh
		}
		return false, time.Time{}, fmt.Errorf("%s: %w", op, err)
	}
	// for res.Next() {
	// 	err := res.Scan(&token, &expiry)
	// 	if err != nil {
	// 		return false, time.Time{}, fmt.Errorf("%s: %w", op, err)
	// 	}
	// }
	// if res.Err() != nil {
	// 	return false, time.Time{}, fmt.Errorf("%s: %w", op, res.Err())
	// }

	return token == refreshToken, expiry, nil
}

func (s *Storage) PutAvatar(ctx context.Context, in *dto.Profile) (string, error) {
	const op = "storage.postgres.PutAvatar"

	state := `UPDATE users SET avatar = $1 WHERE userId = $2`

	_, err := s.db.Exec(ctx, state, in.Avatar, in.Email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return "success", nil
}

func (s *Storage) PutUsername(ctx context.Context, in *dto.Profile) (string, error) {
	const op = "storage.postgres.PutUsername"

	state := `UPDATE users SET name = $1 WHERE userId = $2`

	_, err := s.db.Exec(ctx, state, in.UserName, in.Email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return "success", nil
}

func (s *Storage) PutVerificationCode(ctx context.Context, email string, code string) error {
	const op = "storage.postgres.PutVerificationCode"

	expiryTime := time.Now().UTC().Add(3 * time.Hour).Add(time.Minute * 10)

	state := `UPDATE users SET pwdChangeCode = $1, pwdChangeCodeExpiry = $2 WHERE email = $3`

	_, err := s.db.Exec(ctx, state, code, expiryTime, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrUserNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) GetVerificationCode(ctx context.Context, code string) (string, string, error) {
	const op = "storage.postgres.GetVerificationCode"

	state := `SELECT userId, email, pwdChangeCodeExpiry FROM users WHERE pwdChangeCode = $1`

	var email string
	var userId string
	var expiry time.Time
	err := s.db.QueryRow(ctx, state, code).Scan(&userId, &email, &expiry)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", storage.ErrUserNotFound
		}
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if expiry.Before(time.Now().UTC()) {
		return userId, email, storage.ErrCodeExpired
	}

	return userId, email, nil
}

func (s *Storage) ChangeUserPassword(ctx context.Context, newHash string, userID string) error {
	const op = "storage.postgres.ChangeUserPassword"

	state := `UPDATE users SET passwordHash = $1 WHERE userId = $2`

	result, err := s.db.Exec(ctx, state, newHash, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("%s: password change failed", op)
	}

	return nil
}

// func (s *Storage) GetUserInfoByAuthToken(ctx context.Context, userID string) (*model.UserInfo, error) {
// 	const op = "storage.postgres.GetUserInfoByAuthToken"

// 	state := `SELECT name, surname, jobTitle, organisation, phone, email, city, shortOrganisationTitle, INN, organisationType, currentTarrif FROM users WHERE userId = $1`
// 	res := s.db.QueryRow(ctx, state, userID)
// 	user := &model.UserInfo{}
// 	err := res.Scan(&user.Name, &user.Surname, &user.JobTitle, &user.Organisation, &user.Phone, &user.Email, &user.City,
// 		&user.ShortOrgTitle, &user.INN, &user.OrganisationType, &user.CurrentTarrif)

// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return user, nil
// }

// func (s *Storage) ChangeUserInfo(ctx context.Context, userInfo *dto.ChangeUserInfoIn, userID string) (*model.UserInfo, error) {
// 	const op = "storage.postgres.ChangeUserInfo"

// 	state := `UPDATE users SET (name, surname, jobTitle, organisation, phone, email, city, shortOrganisationTitle, INN, organisationType, currentTarrif) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) WHERE userId = $12 RETURNING name, surname, jobTitle, organisation, phone, email, city, shortOrganisationTitle, INN, organisationType, currentTarrif`
// 	res := s.db.QueryRow(ctx, state, userInfo.Name, userInfo.Surname, userInfo.JobTitle, userInfo.Organisation, userInfo.Phone, userInfo.Email, userInfo.City, userInfo.ShortOrganisationTitle, userInfo.INN, userInfo.OrganisationType, userInfo.CurrentTarrif, userID)
// 	user := &model.UserInfo{}
// 	err := res.Scan(&user.Name, &user.Surname, &user.JobTitle, &user.Organisation, &user.Phone, &user.Email, &user.City,
// 		&user.ShortOrgTitle, &user.INN, &user.OrganisationType, &user.CurrentTarrif)

// 	fmt.Println(err)

// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return user, nil
// }

func (s *Storage) DeleteAccount(ctx context.Context, userID string) error {
	const op = "storage.postgres.DeleteAccount"

	state := `DELETE FROM users WHERE userId = $1`

	row, err := s.db.Exec(ctx, state, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected := row.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("%s: account deletion failed: user not found", op)
	}

	return nil
}

func (s *Storage) GetUserPasswordById(ctx context.Context, userId string) (string, error) {
	const op = "storage.postgres.GetUserPasswordById"

	state := `SELECT passwordHash FROM users WHERE userId = $1`

	res := s.db.QueryRow(ctx, state, userId)

	var hash string

	err := res.Scan(&hash)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return hash, nil
}

func (s *Storage) SignOutUser(ctx context.Context, userId string) error {
	const op = "storage.postgres.LogOutUser"

	state := `UPDATE sessions SET active = false WHERE userId = $1`

	row, err := s.db.Exec(ctx, state, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected := row.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("%s: logout failed: user not found", op)
	}

	return nil
}

func (s *Storage) GetProfile(ctx context.Context, userId string) (*dto.Profile, error) {
	const op = "storage.postgres.GetUserById"

	state := `SELECT avatar, email, name FROM users WHERE userId = $1`
	res := s.db.QueryRow(ctx, state, userId)
	user := &dto.Profile{}
	err := res.Scan(&user.Avatar, &user.Email, &user.UserName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user.UserId = userId

	return user, nil
}

func (s *Storage) GetProfileByLogin(ctx context.Context, login string) (*dto.Profile, error) {
	const op = "storage.postgres.GetUserByLogin"

	state := `SELECT userId, avatar FROM users WHERE email = $1`
	res := s.db.QueryRow(ctx, state, login)
	user := &dto.Profile{}
	err := res.Scan(&user.UserId, &user.Avatar)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user.Email = login

	return user, nil
}
