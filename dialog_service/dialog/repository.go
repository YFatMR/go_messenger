package dialog

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/YFatMR/go_messenger/core/pkg/czap"
	"github.com/YFatMR/go_messenger/dialog_service/apientity"
	"github.com/YFatMR/go_messenger/dialog_service/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type nullableDialogMessage struct {
	MessageID sql.NullInt64
	SenderID  sql.NullInt64
	Text      sql.NullString
	CreatedAt sql.NullTime
}

func nullableDialogMessageToEntity(msg *nullableDialogMessage) *entity.DialogMessage {
	return &entity.DialogMessage{
		MessageID: entity.MessageID{
			ID: uint64(msg.MessageID.Int64),
		},
		SenderID: entity.UserID{
			ID: uint64(msg.SenderID.Int64),
		},
		Text:      msg.Text.String,
		CreatedAt: msg.CreatedAt.Time,
	}
}

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

func (r *dialogRepository) CreateDialog(ctx context.Context, userID1 *entity.UserID, userData1 *entity.UserData,
	userID2 *entity.UserID, userData2 *entity.UserData,
) (*entity.Dialog, error) {
	createDialogName := func(userData *entity.UserData) string {
		return userData.Name + " " + userData.Surname
	}

	ctx, cancel := context.WithTimeout(ctx, r.settings.OperationTimeout)
	defer cancel()

	tx, err := r.connPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		tx.Rollback(ctx)
		r.logger.ErrorContext(ctx, "Can not start tranzaction", zap.Error(err))
		return nil, err
	}

	dialogID := new(entity.DialogID)
	err = tx.QueryRow(
		ctx, `
		INSERT INTO
			dialogs
		DEFAULT VALUES
		RETURNING
			id;`,
	).Scan(&dialogID.ID)
	if err != nil {
		tx.Rollback(ctx)
		r.logger.ErrorContext(ctx, "Can not create dialog", zap.Error(err))
		return nil, err
	}

	lastMessage := &entity.DialogMessage{
		SenderID: entity.UserID{
			ID: 0,
		},
		Text: "Welcome!",
	}
	err = tx.QueryRow(
		ctx, `
		INSERT INTO
			messages (dialog_id, sender_id, text, viewed)
		VALUES
			($1, $2, $3, TRUE)
		RETURNING
			id, created_at;`,
		dialogID.ID, lastMessage.SenderID.ID, lastMessage.Text,
	).Scan(&lastMessage.MessageID.ID, &lastMessage.CreatedAt)

	if err != nil {
		tx.Rollback(ctx)
		r.logger.ErrorContext(ctx, "Can not create first dialog message", zap.Error(err))
		return nil, err
	}

	_, err = tx.Exec(
		ctx, `
		INSERT INTO
			dialog_members (dialog_id, user_id, dialog_name)
		VALUES
			($1, $2, $3),
			($1, $4, $5);`,
		dialogID.ID,
		userID1.ID, createDialogName(userData2),
		userID2.ID, createDialogName(userData1),
	)
	if err != nil {
		tx.Rollback(ctx)
		r.logger.ErrorContext(ctx, "Can not create info about dialog_members", zap.Error(err))
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		r.logger.ErrorContext(ctx, "Can not commit tranzaction for dialog creation", zap.Error(err))
		return nil, err
	}

	dialog := &entity.Dialog{
		DialogID:            *dialogID,
		Name:                createDialogName(userData2),
		UnreadMessagesCount: 0,
		LastMessage:         *lastMessage,
		LastReadMessage:     *lastMessage,
	}
	return dialog, nil
}

func (r *dialogRepository) GetDialog(ctx context.Context, userID *entity.UserID, dialogID *entity.DialogID) (
	*entity.Dialog, error,
) {
	ctx, cancel := context.WithTimeout(ctx, r.settings.OperationTimeout)
	defer cancel()

	lastMessage := new(entity.DialogMessage)
	lastReadMessage := new(entity.DialogMessage)
	dialog := new(entity.Dialog)

	err := r.connPool.QueryRow(
		ctx, `
		SELECT
			dm.dialog_id,
			dm.dialog_name,
			(
				SELECT
					COUNT(*)
				FROM
					messages AS tmp_msg1
				WHERE
					dm.user_id != tmp_msg1.sender_id AND
					dm.dialog_id = tmp_msg1.dialog_id AND
					viewed = FALSE
			) AS unread_messages_count,
			last_message.id,
			last_message.created_at,
			last_message.sender_id,
			last_message.text,
			last_read_message.id,
			last_read_message.created_at,
			last_read_message.sender_id,
			last_read_message.text
		FROM
			dialog_members AS dm
		INNER JOIN
			messages as last_message
		ON
			last_message.dialog_id = dm.dialog_id AND
			last_message.dialog_id = $1 AND last_message.id = (SELECT tmp_msg3.id FROM messages as tmp_msg3 WHERE tmp_msg3.dialog_id = dm.dialog_id ORDER BY created_at DESC LIMIT 1)
		LEFT JOIN
			messages as last_read_message
		ON
			last_read_message.dialog_id = dm.dialog_id AND
			last_read_message.dialog_id = $1 AND
			last_read_message.id = (SELECT id FROM messages AS tmp_msg4 WHERE tmp_msg4.sender_id != $2 AND tmp_msg4.dialog_id = dm.dialog_id AND tmp_msg4.viewed = TRUE ORDER BY tmp_msg4.created_at DESC LIMIT 1)
		WHERE
			dm.user_id = $2 AND dm.dialog_id = $1;`,
		dialogID.ID,
		userID.ID,
	).Scan(&dialog.DialogID.ID, &dialog.Name, &dialog.UnreadMessagesCount,
		&lastMessage.MessageID.ID, &lastMessage.CreatedAt, &lastMessage.SenderID.ID, &lastMessage.Text,
		&lastReadMessage.MessageID.ID, &lastReadMessage.CreatedAt, &lastReadMessage.SenderID.ID, &lastReadMessage.Text)
	if err != nil {
		r.logger.ErrorContext(ctx, "Can not get dialog", zap.Error(err))
		return nil, err
	}
	dialog.LastMessage = *lastMessage
	dialog.LastReadMessage = *lastReadMessage
	return dialog, nil
}

func (r *dialogRepository) GetDialogs(ctx context.Context, userID *entity.UserID, offset uint64, limit uint64) (
	[]*entity.Dialog, error,
) {
	ctx, cancel := context.WithTimeout(ctx, r.settings.OperationTimeout)
	defer cancel()

	r.logger.DebugContext(
		ctx, "Dialog information:",
		zap.Uint64("uid", userID.ID),
		zap.Uint64("limit", limit),
		zap.Uint64("offset", offset),
	)

	raws, err := r.connPool.Query(
		ctx, `
		SELECT
			dm.dialog_id,
			dm.dialog_name,
			(
				SELECT
					COUNT(*)
				FROM
					messages AS tmp_msg1
				WHERE
					dm.user_id != tmp_msg1.sender_id AND
					dm.dialog_id = tmp_msg1.dialog_id AND
					viewed = FALSE
			) AS unread_messages_count,
			last_message.id,
			last_message.created_at,
			last_message.sender_id,
			last_message.text,
			last_read_message.id,
			last_read_message.created_at,
			last_read_message.sender_id,
			last_read_message.text
		FROM
			dialog_members AS dm
		INNER JOIN
			messages as last_message
		ON
			last_message.dialog_id = dm.dialog_id AND
			last_message.id = (SELECT tmp_msg3.id FROM messages as tmp_msg3 WHERE tmp_msg3.dialog_id = dm.dialog_id ORDER BY created_at DESC LIMIT 1)
		LEFT JOIN
			messages as last_read_message
		ON
			last_read_message.dialog_id = dm.dialog_id AND
			last_read_message.id = (SELECT id FROM messages AS tmp_msg4 WHERE tmp_msg4.sender_id != $1 AND tmp_msg4.dialog_id = dm.dialog_id AND tmp_msg4.viewed = TRUE ORDER BY tmp_msg4.created_at DESC LIMIT 1)
		WHERE
			dm.user_id = $1
		ORDER BY
			last_message.created_at DESC
		LIMIT $2 OFFSET $3;`,
		userID.ID,
		limit,
		offset,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Can not get dialogs", zap.Error(err))
		return nil, err
	}
	defer raws.Close()

	var dialogs []*entity.Dialog
	for raws.Next() {
		lastMessage := new(entity.DialogMessage)
		nullableLastReadMessage := new(nullableDialogMessage)
		dialog := new(entity.Dialog)
		err := raws.Scan(
			&dialog.DialogID.ID, &dialog.Name, &dialog.UnreadMessagesCount,
			&lastMessage.MessageID.ID, &lastMessage.CreatedAt, &lastMessage.SenderID.ID, &lastMessage.Text,
			&nullableLastReadMessage.MessageID, &nullableLastReadMessage.CreatedAt, &nullableLastReadMessage.SenderID,
			&nullableLastReadMessage.Text,
		)
		if err != nil {
			r.logger.ErrorContext(ctx, "Can not scan raws", zap.Error(err))
			return nil, err
		}

		dialog.LastMessage = *lastMessage
		dialog.LastReadMessage = *nullableDialogMessageToEntity(nullableLastReadMessage)

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

	responseMessage := entity.CopyDialogMessage(message)
	err := r.connPool.QueryRow(
		ctx, `
		INSERT INTO
			messages (dialog_id, sender_id, text)
		VALUES
			($1, $2, $3)
		RETURNING
			id, created_at;`,
		dialogID.ID,
		message.SenderID.ID,
		message.Text,
	).Scan(&responseMessage.MessageID.ID, &responseMessage.CreatedAt)
	if err != nil {
		r.logger.ErrorContext(ctx, "Unable to create message", zap.Error(err))
		return nil, ErrCreateMessage
	}
	return responseMessage, nil
}

func (r *dialogRepository) GetDialogMessageByID(ctx context.Context, dialogID *entity.DialogID,
	messageID *entity.MessageID,
) (
	*entity.DialogMessage, error,
) {
	ctx, cancel := context.WithTimeout(ctx, r.settings.OperationTimeout)
	defer cancel()

	message := new(entity.DialogMessage)
	err := r.connPool.QueryRow(
		ctx, `
		SELECT
			id, created_at, sender_id, text, viewed
		FROM
			messages
		WHERE
			dialog_id = $1 AND id = $2;`,
		dialogID.ID,
		messageID.ID,
	).Scan(&message.MessageID.ID, &message.CreatedAt, &message.SenderID.ID, &message.Text, &message.Viewed)
	if err != nil {
		r.logger.ErrorContext(ctx, "Can not get message by ID", zap.Error(err))
		return nil, err
	}
	return message, nil
}

func (r *dialogRepository) getDialogMessagesBefore(ctx context.Context, dialogID *entity.DialogID,
	messageID *entity.MessageID, limit uint64, include bool,
) (
	[]*entity.DialogMessage, error,
) {
	ctx, cancel := context.WithTimeout(ctx, r.settings.OperationTimeout)
	defer cancel()

	sign := "<"
	if include {
		sign = "<="
	}

	raws, err := r.connPool.Query(
		ctx, fmt.Sprintf(`
		SELECT *
		FROM (
			SELECT
				id, created_at, sender_id, text, viewed
			FROM
				messages
			WHERE
				dialog_id = $1 AND id %s $2
			ORDER BY
				created_at DESC
			LIMIT $3) as tmp
		ORDER BY tmp.id`, sign),
		dialogID.ID,
		messageID.ID,
		limit,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Unable to get messages before", zap.Error(err))
		return nil, ErrCreateMessage
	}
	defer raws.Close()

	var messages []*entity.DialogMessage
	for raws.Next() {
		message := new(entity.DialogMessage)
		err := raws.Scan(
			&message.MessageID.ID, &message.CreatedAt, &message.SenderID.ID, &message.Text, &message.Viewed,
		)
		if err != nil {
			r.logger.ErrorContext(ctx, "Unable to scan messages", zap.Error(err))
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (r *dialogRepository) GetDialogMessagesBefore(ctx context.Context, dialogID *entity.DialogID,
	messageID *entity.MessageID, limit uint64,
) (
	[]*entity.DialogMessage, error,
) {
	return r.getDialogMessagesBefore(ctx, dialogID, messageID, limit, false)
}

func (r *dialogRepository) GetDialogMessagesBeforeAndInclude(ctx context.Context, dialogID *entity.DialogID,
	messageID *entity.MessageID, limit uint64,
) (
	[]*entity.DialogMessage, error,
) {
	return r.getDialogMessagesBefore(ctx, dialogID, messageID, limit, true)
}

func (r *dialogRepository) getDialogMessagesAfter(ctx context.Context, dialogID *entity.DialogID,
	messageID *entity.MessageID, limit uint64, include bool,
) (
	[]*entity.DialogMessage, error,
) {
	ctx, cancel := context.WithTimeout(ctx, r.settings.OperationTimeout)
	defer cancel()

	sign := ">"
	if include {
		sign = ">="
	}

	raws, err := r.connPool.Query(
		ctx, fmt.Sprintf(`
		SELECT
			id, created_at, sender_id, text, viewed
		FROM
			messages
		WHERE
			dialog_id = $1 AND id %s $2
		ORDER BY
			created_at
		LIMIT $3`, sign),
		dialogID.ID,
		messageID.ID,
		limit,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Unable to get messages after", zap.Error(err))
		return nil, ErrCreateMessage
	}
	defer raws.Close()

	var messages []*entity.DialogMessage
	for raws.Next() {
		message := new(entity.DialogMessage)
		err := raws.Scan(
			&message.MessageID.ID, &message.CreatedAt, &message.SenderID.ID, &message.Text, &message.Viewed,
		)
		if err != nil {
			r.logger.ErrorContext(ctx, "Unable to scan messages after", zap.Error(err))
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (r *dialogRepository) GetDialogMessagesAfter(ctx context.Context, dialogID *entity.DialogID,
	messageID *entity.MessageID, limit uint64,
) (
	[]*entity.DialogMessage, error,
) {
	return r.getDialogMessagesAfter(ctx, dialogID, messageID, limit, false)
}

func (r *dialogRepository) GetDialogMessagesAfterAndInclude(ctx context.Context, dialogID *entity.DialogID,
	messageID *entity.MessageID, limit uint64,
) (
	[]*entity.DialogMessage, error,
) {
	return r.getDialogMessagesAfter(ctx, dialogID, messageID, limit, true)
}

// func (r *dialogRepository) GetDialogMessages(ctx context.Context, dialogID *entity.DialogID,
// 	messageID *entity.MessageID, limit uint64, offsetType entity.OffserType,
// ) (
// 	[]*entity.DialogMessage, error,
// ) {
// 	ctx, cancel := context.WithTimeout(ctx, r.settings.OperationTimeout)
// 	defer cancel()

// 	getMessagesBefore := func() (pgx.Rows, error) {
// 		return r.connPool.Query(
// 			ctx, `
// 				SELECT *
// 				FROM (
// 					SELECT
// 						id, created_at, sender_id, text, viewed
// 					FROM
// 						messages
// 					WHERE
// 						dialog_id = $1 AND id < $2
// 					ORDER BY
// 						created_at DESC
// 					LIMIT $3) as tmp
// 				ORDER BY tmp.id`,
// 			dialogID.ID,
// 			messageID.ID,
// 			limit,
// 		)
// 	}

// 	getMessagesAfter := func() (pgx.Rows, error) {
// 		return r.connPool.Query(
// 			ctx, `
// 				SELECT
// 					id, created_at, sender_id, text, viewed
// 				FROM
// 					messages
// 				WHERE
// 					dialog_id = $1 AND id > $2
// 				ORDER BY
// 					created_at
// 				LIMIT $3`,
// 			dialogID.ID,
// 			messageID.ID,
// 			limit,
// 		)
// 	}

// 	var raws pgx.Rows
// 	var err error
// 	if offsetType == entity.BEFORE {
// 		r.logger.DebugContext(ctx, "Process BEFORE messages selection")
// 		raws, err = getMessagesBefore()
// 	} else {
// 		r.logger.DebugContext(ctx, "Process AFTER messages selection")
// 		raws, err = getMessagesAfter()
// 	}

// 	if err != nil {
// 		r.logger.ErrorContext(ctx, "Unable to get messages", zap.Error(err))
// 		return nil, ErrCreateMessage
// 	}
// 	defer raws.Close()

// 	var messages []*entity.DialogMessage
// 	for raws.Next() {
// 		message := new(entity.DialogMessage)
// 		err := raws.Scan(
// 			&message.MessageID.ID, &message.CreatedAt, &message.SenderID.ID, &message.Text, &message.Viewed,
// 		)
// 		if err != nil {
// 			r.logger.ErrorContext(ctx, "Unable to scan messages", zap.Error(err))
// 			return nil, err
// 		}
// 		messages = append(messages, message)
// 	}
// 	return messages, nil
// }

func (r *dialogRepository) GetDialogMembers(ctx context.Context, dialogID *entity.DialogID) (
	[]*entity.UserID, error,
) {
	ctx, cancel := context.WithTimeout(ctx, r.settings.OperationTimeout)
	defer cancel()

	raws, err := r.connPool.Query(
		ctx, `
		SELECT
			user_id
		FROM
			dialog_members
		WHERE
			dialog_id = $1`,
		dialogID.ID,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Unable to get dialog members", zap.Error(err))
		return nil, ErrCreateDialog
	}

	userIDs := make([]*entity.UserID, 0, 2)
	for raws.Next() {
		userID := new(entity.UserID)
		err := raws.Scan(&userID.ID)
		if err != nil {
			r.logger.ErrorContext(ctx, "Unable to scan dialog members", zap.Error(err))
			return nil, ErrCreateDialog
		}
		userIDs = append(userIDs, userID)
	}
	return userIDs, nil
}

// Include message
func (r *dialogRepository) ReadAllMessagesBeforeAndIncl(ctx context.Context, userID *entity.UserID, dialogID *entity.DialogID,
	messageID *entity.MessageID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.settings.OperationTimeout)
	defer cancel()

	_, err := r.connPool.Exec(
		ctx, `
		UPDATE messages
		SET
			viewed = TRUE
		WHERE
			sender_id != $1 AND
			dialog_id = $3 AND
			created_at <= (SELECT tmp_msg1.created_at FROM messages as tmp_msg1 WHERE tmp_msg1.id = $2);`,
		userID.ID,
		messageID.ID,
		dialogID.ID,
	)
	if err != nil {
		r.logger.ErrorContext(ctx, "Unable to update last read message", zap.Error(err))
		return err
	}
	return nil
}
