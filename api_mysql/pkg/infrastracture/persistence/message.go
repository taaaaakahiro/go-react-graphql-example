package persistence

import (
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/domain/entity"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/domain/repository"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mysql/pkg/io"
	errs "github.com/pkg/errors"
)

type MessageRepo struct {
	database *io.SQLDatabase
}

var _ repository.IMessageRepository = (*MessageRepo)(nil)

func NewMessageRepository(db *io.SQLDatabase) *MessageRepo {
	return &MessageRepo{
		database: db,
	}
}

func (r MessageRepo) ListMessages(userId int) ([]entity.Message, error) {
	messages := make([]entity.Message, 0)

	query := "SELECT id, user_id, message FROM message WHERE user_id = ? ORDER BY id DESC"
	stmtOut, err := r.database.Prepare(query)
	if err != nil {
		return nil, errs.WithStack(err)
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId)
	if err != nil {
		return nil, errs.WithStack(err)
	}

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		message := entity.Message{}

		err = rows.Scan(&message.Id, &message.UserId, &message.Message)
		if err != nil {
			return nil, errs.WithStack(err)
		}

		messages = append(messages, message)
	}
	if err = rows.Err(); err != nil {
		return nil, errs.WithStack(err)
	}
	return messages, nil
}

func (r MessageRepo) Messages() ([]entity.Message, error) {
	messages := make([]entity.Message, 0)

	query := "SELECT id, user_id, message FROM message ORDER BY id DESC"
	stmtOut, err := r.database.Prepare(query)
	if err != nil {
		return nil, errs.WithStack(err)
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query()
	if err != nil {
		return nil, errs.WithStack(err)
	}

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		message := entity.Message{}

		err = rows.Scan(&message.Id, &message.UserId, &message.Message)
		if err != nil {
			return nil, errs.WithStack(err)
		}

		messages = append(messages, message)
	}
	if err = rows.Err(); err != nil {
		return nil, errs.WithStack(err)
	}
	return messages, nil
}

func (r MessageRepo) CreateMessage(message *entity.Message) error {
	tx, cancel, err := r.database.Begin()
	if err != nil {
		return errs.WithStack(err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
		cancel()
	}()

	// insert
	query := "INSERT INTO message(user_id, message) VALUES (?, ?);"
	ret, err := tx.Exec(query, message.UserId, message.Message)
	if err != nil {
		return errs.WithStack(err)
	}
	err = tx.Commit()
	if err != nil {
		return errs.WithStack(err)
	}

	lastId, err := ret.LastInsertId()
	if err != nil {
		return errs.WithStack(err)
	}
	message.Id = int(lastId)

	return nil
}
