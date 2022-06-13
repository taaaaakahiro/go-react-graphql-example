package persistence

import (
	"context"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mongo/pkg/domain/entity"
	"github.com/cyberwo1f/go-and-react-graphql-example/api_mongo/pkg/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const (
	userCollection = "user"
)

type UserRepo struct {
	col *mongo.Collection
}

var _ repository.IUserRepository = (*UserRepo)(nil)

func NewUserRepository(db *mongo.Database) *UserRepo {
	return &UserRepo{
		col: db.Collection(userCollection),
	}
}

func (r UserRepo) ListUsers(ctx context.Context) ([]entity.User, error) {
	users := make([]entity.User, 0)
	srt := bson.D{
		primitive.E{Key: "id", Value: -1},
	}
	opt := options.Find().SetSort(srt)
	cur, err := r.col.Find(ctx, bson.D{}, opt)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var user entity.User
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	err = cur.Close(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r UserRepo) User(ctx context.Context, userId int) (entity.User, error) {
	user := entity.User{}
	flt := bson.D{
		primitive.E{Key: "id", Value: userId},
	}
	opt := options.FindOne()
	err := r.col.FindOne(ctx, flt, opt).Decode(&user)
	if err == mongo.ErrNoDocuments {
		log.Println("Documents not found")
	}
	return user, err
}
