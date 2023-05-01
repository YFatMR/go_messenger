package dialog

import (
	"context"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/dialog_service/apientity"
	"github.com/YFatMR/go_messenger/dialog_service/entity"
	"github.com/jackc/pgx/v5"
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
			CASE WHEN d.user_id_1 = $1 THEN d.dialog_name_1 ELSE d.dialog_name_2 END AS dialog_name,
			CASE WHEN d.user_id_1 = $1 THEN d.unread_messages_count_1 ELSE d.unread_messages_count_2 END AS unread_messages_count,
			m2.messages_count,
			m.id,
			m.created_at,
			m.sender_id,
			m.text
		FROM
			dialogs as d
		INNER JOIN
			messages as m
		ON
			m.dialog_id = d.id AND m.created_at = (SELECT MAX(created_at) FROM messages AS tmp_m WHERE tmp_m.dialog_id = d.id)
		INNER JOIN
			(SELECT tmp_m3.dialog_id, COUNT(*) as messages_count FROM messages as tmp_m3 GROUP BY tmp_m3.dialog_id) as m2
		ON
			m2.dialog_id = d.id
		WHERE
			d.user_id_1 = $1 OR d.user_id_2 = $1
		ORDER BY
			m.created_at
		LIMIT $2 OFFSET $3;`,
		userID.ID,
		limit,
		offset,
	)

	// raws, err := r.connPool.Query(
	// 	ctx, `
	// 	SELECT
	// 		d.id,
	// 		CASE WHEN d.user_id_1 = $1 THEN d.dialog_name_1 ELSE d.dialog_name_2 END AS dialog_name,
	// 		CASE WHEN d.user_id_1 = $1 THEN d.unread_messages_count_1 ELSE d.unread_messages_count_2 END AS unread_messages_count
	// 	FROM
	// 		dialogs as d
	// 	WHERE
	// 		d.user_id_1 = $1 OR d.user_id_2 = $1
	// 	LIMIT $2 OFFSET $3;`,
	// 	userID.ID,
	// 	limit,
	// 	offset,
	// )
	if err != nil {
		return nil, err
	}
	defer raws.Close()

	r.logger.DebugContext(ctx, "Dialog information:", zap.Uint64("uid", userID.ID), zap.Uint64("limit", limit), zap.Uint64("offset", offset))

	var dialogs []*entity.Dialog
	for raws.Next() {
		message := new(entity.DialogMessage)
		dialog := new(entity.Dialog)

		r.logger.DebugContext(ctx, "Dialog scabbing...")

		err := raws.Scan(
			&dialog.DialogID.ID, &dialog.Name, &dialog.UnreadMessagesCount, &dialog.MessagesCount,
			&message.MessageID.ID, &message.CreatedAt, &message.SenderID.ID, &message.Text,
		)

		// err := raws.Scan(
		// 	&dialog.DialogID.ID, &dialog.Name, &dialog.UnreadMessagesCount,
		// )
		if err != nil {
			return nil, err
		}
		dialog.LastMessage = *message
		dialogs = append(dialogs, dialog)
	}
	return dialogs, nil
}

func (r *dialogRepository) CreateDialogMessage(ctx context.Context, dialogID *entity.DialogID,
	message *entity.DialogMessage,
) (
	*entity.DialogMessage, error,
) {
	ctx, cancel := context.WithTimeout(ctx, r.settings.OperationTimeout)
	defer cancel()

	tx, err := r.connPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.logger.ErrorContext(ctx, "Unable to start transacrion to create message", zap.Error(err))
		return nil, ErrCreateMessage
	}

	responseMessage := entity.CopyDialogMessage(message)
	err = tx.QueryRow(
		ctx, `
		INSERT INTO
			messages (dialog_id, sender_id, text)
		VALUES
			($1, $2, $3)
		RETURNING
			created_at;`,
		dialogID.ID,
		message.SenderID.ID,
		message.Text,
	).Scan(&responseMessage.CreatedAt)
	if err != nil {
		r.logger.ErrorContext(ctx, "Unable to create message", zap.Error(err))
		return nil, ErrCreateMessage
	}

	_, err = tx.Exec(
		ctx, `
		UPDATE dialogs
		SET
			unread_messages_count_1 =
				CASE
					WHEN user_id_1 = $1 THEN unread_messages_count_1
					ELSE unread_messages_count_1 + 1
				END,
			unread_messages_count_2 =
				CASE
					WHEN user_id_2 = $1 THEN unread_messages_count_2
					ELSE unread_messages_count_2 + 1
				END
		WHERE id = $2 AND (user_id_1 = $1 OR user_id_2 = $1);`,
		message.SenderID.ID,
		dialogID.ID,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Can not update unread messages count", zap.Error(err))
		return nil, ErrCreateMessage
	}

	if err := tx.Commit(ctx); err != nil {
		r.logger.ErrorContext(ctx, "Unable to commit transaction", zap.Error(err))
	}

	return responseMessage, nil
}

func (r *dialogRepository) GetDialogMessages(ctx context.Context, dialogID *entity.DialogID,
	offset uint64, limit uint64,
) (
	[]*entity.DialogMessage, error,
) {
	ctx, cancel := context.WithTimeout(ctx, r.settings.OperationTimeout)
	defer cancel()

	raws, err := r.connPool.Query(
		ctx, `
		SELECT
			id, created_at, sender_id, text
		FROM
			messages
		WHERE
			dialog_id = $1
		ORDER BY
			created_at
		LIMIT $2 OFFSET $3`,
		dialogID.ID,
		limit,
		offset,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Unable to get messages", zap.Error(err))
		return nil, ErrCreateMessage
	}
	defer raws.Close()

	var messages []*entity.DialogMessage
	for raws.Next() {
		message := new(entity.DialogMessage)
		err := raws.Scan(
			&message.MessageID.ID, &message.CreatedAt, &message.SenderID.ID, &message.Text,
		)
		if err != nil {
			r.logger.ErrorContext(ctx, "Unable to scan messages", zap.Error(err))
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
