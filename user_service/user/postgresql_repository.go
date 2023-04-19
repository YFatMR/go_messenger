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

// type userDocument struct {
// 	ID             primitive.ObjectID `bson:"_id,omitempty"`
// 	Email          string             `bson:"login,omitempty"`
// 	HashedPassword string             `bson:"hashed_password,omitempty"`
// 	UserRole       string             `bson:"user_role,omitempty"`
// 	Nickname       string             `bson:"nickname,omitempty"`
// 	Name           string             `bson:"name,omitempty"`
// 	Surname        string             `bson:"surname,omitempty"`
// }

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
			nickname, name, surname
		FROM
			users
		WHERE
			id = $1;`,
		userID.ID,
	).Scan(&userData.Nickname, &userData.Name, &userData.Surname)
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
			id, email, hashed_password, role, nickname, name, surname
		FROM
			users
		WHERE
			email = $1;`,
		email,
	).Scan(
		&account.UserID.ID, &account.Email, &account.HashedPassword,
		&account.Role.Name, &account.Nickname, &account.Name, &account.Surname,
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
