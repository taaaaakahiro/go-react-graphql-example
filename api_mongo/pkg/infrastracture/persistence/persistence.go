package persistence

import (
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mongo/pkg/domain/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	User    repository.IUserRepository
	Message repository.IMessageRepository
}

func NewRepositories(db *mongo.Database) (*Repositories, error) {
	return &Repositories{
		User:    NewUserRepository(db),
		Message: NewMessageRepository(db),
	}, nil
}
