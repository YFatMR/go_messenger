package user

import (
	"context"
	"errors"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/user_service/apientity"
	"github.com/YFatMR/go_messenger/user_service/entity"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type userPosgreSQLRepository struct {
	connPool         *pgxpool.Pool
	operationTimeout time.Duration
	logger           *czap.Logger
}

func NewPosgreSQLRepository(connPool *pgxpool.Pool, operationTimeout time.Duration,
	logger *czap.Logger,
) apientity.UserRepository {
	return &userPosgreSQLRepository{
		connPool:         connPool,
		operationTimeout: operationTimeout,
		logger:           logger,
	}
}

func (r *userPosgreSQLRepository) Create(ctx context.Context, user *entity.User, credential *entity.Credential) (
	*entity.UserID, error,
) {
	dbCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	userID := new(entity.UserID)
	err := r.connPool.QueryRow(
		dbCtx, `
		INSERT INTO
			users (email, hashed_password, role, nickname, name, surname)
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING
			id;`,
		credential.Email, credential.HashedPassword, credential.Role.Name,
		user.Nickname, user.Name, user.Surname,
	).Scan(&userID.ID)

	if err != nil {
		r.logger.ErrorContext(ctx, "Can't insert new user", zap.Error(err))
		return nil, ErrUserCreation
	}
	return userID, nil
}

func (r *userPosgreSQLRepository) GetByID(ctx context.Context, userID *entity.UserID) (
	*entity.User, error,
) {
	dbCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	userData := new(entity.User)
	err := r.connPool.QueryRow(
		dbCtx, `
		SELECT
			nickname, name, surname, github, linkedin, public_email, status
		FROM
			users
		WHERE
			id = $1;`,
		userID.ID,
	).Scan(
		&userData.Nickname, &userData.Name, &userData.Surname, &userData.Github,
		&userData.Linkedin, &userData.PublicEmail, &userData.Status,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.ErrorContext(ctx, "Can't get user", zap.Uint64("user_id", userID.ID), zap.Error(err))
		return nil, ErrUserNotFound
	} else if err != nil {
		r.logger.ErrorContext(ctx, "Database connection error", zap.Error(err))
		return nil, ErrUserNotFound
	}
	return userData, nil
}

func (r *userPosgreSQLRepository) DeleteByID(ctx context.Context, userID *entity.UserID) error {
	dbCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	_, err := r.connPool.Exec(
		dbCtx, `
		DELETE FROM
			users
		WHERE
			id = $1;`,
		userID.ID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.ErrorContext(ctx, "User not found by id", zap.Error(err))
		return ErrUserNotFound
	} else if err != nil {
		r.logger.ErrorContext(ctx, "Database connection error", zap.Error(err))
		return ErrUserNotFound
	}
	return nil
}

func (r *userPosgreSQLRepository) GetAccountByEmail(ctx context.Context, email string) (
	*entity.Account, error,
) {
	dbCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	account := new(entity.Account)
	err := r.connPool.QueryRow(
		dbCtx, `
		SELECT
			id, email, hashed_password, role, nickname, name, surname, github, linkedin, public_email, status
		FROM
			users
		WHERE
			email = $1;`,
		email,
	).Scan(
		&account.UserID.ID, &account.Email, &account.HashedPassword,
		&account.Role.Name, &account.Nickname, &account.Name, &account.Surname,
		&account.Github, &account.Linkedin, &account.PublicEmail, &account.Status,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.ErrorContext(ctx, "User credential not found", zap.String("email", email), zap.Error(err))
		return nil, ErrUserNotFound
	} else if err != nil {
		r.logger.ErrorContext(ctx, "Database connection error", zap.Error(err))
		return nil, ErrGetToken
	}
	return account, nil
}

func (r *userPosgreSQLRepository) UpdateUserData(ctx context.Context, userID *entity.UserID,
	request *entity.UpdateUserRequest,
) error {
	dbCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	_, err := r.connPool.Exec(
		dbCtx, `
		UPDATE
			users
		SET
			name = $1,
			surname = $2,
			github = $3,
			linkedin = $4,
			public_email = $5,
			status = $6
		WHERE
			id = $7;`,
		request.Name, request.Status, request.Github,
		request.Linkedin, request.PublicEmail,
		userID.ID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.ErrorContext(ctx, "User not found")
		return ErrUserNotFound
	} else if err != nil {
		r.logger.ErrorContext(ctx, "Database connection error", zap.Error(err))
		return ErrUserNotFound
	}
	return nil
}

func (r *userPosgreSQLRepository) GetUsersByPrefix(ctx context.Context, selfID *entity.UserID,
	prefix string, limit uint64,
) (
	[]*entity.UserWithID, error,
) {
	dbCtx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()

	rows, err := r.connPool.Query(
		dbCtx, `
		SELECT
			id, nickname, name, surname, github, linkedin, public_email, status
		FROM
			users
		WHERE
			nickname LIKE $1 AND id != $3
		ORDER BY
			LENGTH(nickname)
		LIMIT
			$2;`,
		prefix+"%",
		limit,
		selfID.ID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.DebugContext(ctx, "User not found")
		return make([]*entity.UserWithID, 0), nil
	} else if err != nil {
		r.logger.ErrorContext(ctx, "Database connection error", zap.Error(err))
		return nil, err
	}

	result := make([]*entity.UserWithID, 0, 16)
	for rows.Next() {
		var userWithID entity.UserWithID
		err = rows.Scan(&userWithID.UserID.ID, &userWithID.Nickname, &userWithID.Name,
			&userWithID.Surname, &userWithID.Github, &userWithID.Linkedin,
			&userWithID.PublicEmail, &userWithID.Status)
		if err != nil {
			r.logger.ErrorContext(ctx, "Failed to scan row", zap.Error(err))
			return nil, err
		}
		result = append(result, &userWithID)
	}
	return result, nil
}
