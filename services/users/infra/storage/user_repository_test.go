package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/packages/logger"
	"github.com/italoservio/braz_ecommerce/services/users/infra/storage"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

type MockStructure struct {
	Id  string `bson:"_id"`
	Foo string `bson:"foo"`
}

const (
	MOCK_DB_NAME   = "foo"
	MOCK_COLL_NAME = "users"
	MOCK_NS        = "foo.users"
)

func TestUserRepository_NewUserRepository(t *testing.T) {

	logger := logger.NewLogger()
	rootMt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	rootMt.Run("should return a new instance when all right", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockDB := &database.Database{Database: nestedMt.Client.Database(mockDbName)}

		instance := storage.NewUserRepositoryImpl(logger, mockDB)

		assert.Equal(
			t,
			fmt.Sprintf("%T", instance),
			"*storage.UserRepositoryImpl",
			"should be a pointer to UserRepositoryImpl",
		)
	})
}

func TestCrudRepository_GetByEmail(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewLogger()
	rootMt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	rootMt.Run("should return the document when call database with success", func(nestedMt *mtest.T) {
		mockId := primitive.NewObjectID()

		nestedMt.AddMockResponses(mtest.CreateCursorResponse(
			1,
			MOCK_NS,
			mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: mockId},
				{Key: "foo", Value: "bar"},
			},
		))
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(MOCK_DB_NAME)}
		crudRepository := storage.NewUserRepositoryImpl(logger, mockDB)

		var result MockStructure

		err := crudRepository.GetByEmail(ctx, MOCK_COLL_NAME, "", &result)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.Nil(t, err, "should not return error")
		assert.Equal(t, "bar", result.Foo, "should return the expected object by id")
	})

	rootMt.Run("should return error when no document is found", func(nestedMt *mtest.T) {
		nestedMt.AddMockResponses(mtest.CreateCursorResponse(
			0,
			MOCK_NS,
			mtest.FirstBatch,
		))

		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(MOCK_DB_NAME)}
		crudRepository := storage.NewUserRepositoryImpl(logger, mockDB)

		var result MockStructure

		err := crudRepository.GetByEmail(ctx, MOCK_COLL_NAME, "", &result)
		if err == nil {
			t.Fail()
		}

		assert.Equal(t, exception.CodeNotFound, err.Error(), "should return the expected error")
	})

	rootMt.Run("should return error when failed to call database", func(nestedMt *mtest.T) {
		nestedMt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(MOCK_DB_NAME)}
		crudRepository := storage.NewUserRepositoryImpl(logger, mockDB)

		var result MockStructure

		err := crudRepository.GetByEmail(ctx, MOCK_COLL_NAME, "", &result)
		if err == nil {
			t.Fail()
		}

		assert.Equal(t, exception.CodeDatabaseFailed, err.Error(), "should return database call error")
	})
}
