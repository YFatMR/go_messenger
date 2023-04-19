package dialog

import (
	"context"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/dialog_service/apientity"
	"github.com/YFatMR/go_messenger/dialog_service/entity"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DialogRepositorySettings struct {
	OperationTimeout time.Duration
}

type dialogRepository struct {
	settings DialogRepositorySettings
	connPool *pgxpool.Pool
	logger   *czap.Logger
}

func NewPosgreRepository(settings DialogRepositorySettings, connPool *pgxpool.Pool, logger *czap.Logger,
) apientity.DialogRepository {
	return &dialogRepository{
		settings: settings,
		connPool: connPool,
		logger:   logger,
	}
}

// name userID unreadMessagesCount

func (r *dialogRepository) CreateDialog(ctx context.Context, userID1 *entity.UserID, userData1 *entity.UserData,
	userID2 *entity.UserID, userData2 *entity.UserData,
) (*entity.Dialog, error) {
	createDialogName := func(userData *entity.UserData) string {
		return userData.Name + " " + userData.Surname
	}

	ctx, cancel := context.WithTimeout(ctx, r.settings.OperationTimeout)
	defer cancel()

	dialog := new(entity.Dialog)
	err := r.connPool.QueryRow(
		ctx, `
		INSERT INTO
			dialogs (user_id_1, dialog_name_1, user_id_2, dialog_name_2)
		VALUES
			($1, $2, $3, $4)
		RETURNING
			id, dialog_name_1, unread_messages_count_1;`,
		userID1.ID, createDialogName(userData2),
		userID2.ID, createDialogName(userData1),
	).Scan(&dialog.DialogID.ID, &dialog.Name, &dialog.UnreadMessagesCount)
	if err != nil {
		r.logger.ErrorContext(ctx, "Unable to create dialog", zap.Error(err))
		return nil, ErrCreateDialog
	}
	return dialog, nil
}

func (r *dialogRepository) GetDialogs(ctx context.Context, userID *entity.UserID, offset uint64, limit uint64) (
	[]*entity.Dialog, error,
) {
	ctx, cancel := context.WithTimeout(ctx, r.settings.OperationTimeout)
	defer cancel()

	raws, err := r.connPool.Query(
		ctx, `
		SELECT
			d.id,
			IF(d.user_id_1 = $1, d.dialog_name_1, d.dialog_name_2) AS dialog_name,
			IF(d.user_id_1 = $1, d.unread_messages_count_1, d.unread_messages_count_2) AS unread_messages_count,
			m.created_at,
			m.sender_id,
			m.text
		FROM
			dialogs as d
		INNER JOIN
			messages as m
		ON
			m.dialog_id = d.id AND m.created_at = (SELECT MAX(created_at) FROM messages AS tmp_m WHERE tmp_m.dialog_id = d.id)
		WHERE
			userID1 = $1 OR userID2 = $1
		ORDER BY
			lastMessageDate
		LIMIT $2 OFFSET $3
		;`,
		userID.ID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer raws.Close()

	var dialogs []*entity.Dialog
	for raws.Next() {
		message := new(entity.DialogMessage)
		dialog := new(entity.Dialog)
		var timestamp time.Time

		err := raws.Scan(
			&dialog.DialogID.ID, &dialog.Name, &dialog.UnreadMessagesCount,
			&timestamp, &message.SenderID, &message.Text,
		)
		message.CreatedAt = uint64(timestamp.UnixNano())
		if err != nil {
			return nil, err
		}
		dialogs = append(dialogs, dialog)
	}
	return dialogs, nil
}
