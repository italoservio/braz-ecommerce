package storage

import (
	"context"
	"errors"
	"time"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/packages/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryInterface interface {
	GetByEmail(
		ctx context.Context,
		collection string,
		email string,
		structure any) error
}

type UserRepositoryImpl struct {
	logger   logger.LoggerInterface
	database *database.Database
}

func NewUserRepositoryImpl(lg logger.LoggerInterface, db *database.Database) *UserRepositoryImpl {
	return &UserRepositoryImpl{logger: lg, database: db}
}

func (cr *UserRepositoryImpl) GetByEmail(
	ctx context.Context,
	collection string,
	email string,
	structure any,
) error {

	coll := cr.database.Collection(collection)

	cursor, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := coll.FindOne(cursor, bson.M{"email": email}).Decode(structure)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			cr.logger.WithCtx(ctx).Error(err.Error())
			return errors.New(exception.CodeNotFound)
		}
		cr.logger.WithCtx(ctx).Error(err.Error())
		return errors.New(exception.CodeDatabaseFailed)
	}

	return nil
}
