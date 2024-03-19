package storage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/packages/logger"
	"github.com/italoservio/braz_ecommerce/services/users/domain"
	"github.com/italoservio/braz_ecommerce/services/users/infra/storage"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

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

type TestingDependencies_TestGetByEmail struct {
	ctx            context.Context
	mockDB         *database.Database
	userRepository *storage.UserRepositoryImpl
}

func BeforeEach_TestGetByEmail(mt *mtest.T) *TestingDependencies_TestGetByEmail {
	ctx := context.TODO()
	mockDB := &database.Database{Database: mt.Client.Database(MOCK_DB_NAME)}
	userRepository := storage.NewUserRepositoryImpl(logger.NewLogger(), mockDB)

	return &TestingDependencies_TestGetByEmail{
		ctx:            ctx,
		mockDB:         mockDB,
		userRepository: userRepository,
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	rootMt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	rootMt.Run("should return the document when call database with success", func(nestedMt *mtest.T) {
		deps := BeforeEach_TestGetByEmail(nestedMt)
		mockId := primitive.NewObjectID()

		nestedMt.AddMockResponses(mtest.CreateCursorResponse(
			1,
			MOCK_NS,
			mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: mockId},
				{Key: "first_name", Value: "bar"},
			},
		))
		defer nestedMt.ClearMockResponses()

		var result domain.UserDatabaseNoPassword

		err := deps.userRepository.GetByEmail(
			deps.ctx,
			MOCK_COLL_NAME,
			"",
			&result,
		)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.Nil(t, err, "should not return error")
		assert.Equal(t, "bar", result.FirstName, "should return the expected first name")
	})

	rootMt.Run("should return empty when no document is found", func(nestedMt *mtest.T) {
		deps := BeforeEach_TestGetByEmail(nestedMt)
		nestedMt.AddMockResponses(mtest.CreateCursorResponse(
			0,
			MOCK_NS,
			mtest.FirstBatch,
		))

		defer nestedMt.ClearMockResponses()

		var result domain.UserDatabaseNoPassword

		err := deps.userRepository.GetByEmail(
			deps.ctx,
			MOCK_COLL_NAME,
			"",
			&result,
		)
		if err != nil {
			t.Fail()
		}

		assert.Equal(t, (domain.UserDatabaseNoPassword{}), result, "should return empty result")
	})

	rootMt.Run("should return error when failed to call database", func(nestedMt *mtest.T) {
		deps := BeforeEach_TestGetByEmail(nestedMt)

		nestedMt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})
		defer nestedMt.ClearMockResponses()

		var result domain.UserDatabaseNoPassword

		err := deps.userRepository.GetByEmail(
			deps.ctx,
			MOCK_COLL_NAME,
			"",
			&result,
		)
		if err == nil {
			t.Fail()
		}

		assert.Equal(t, exception.CodeDatabaseFailed, err.Error(), "should return database call error")
	})
}
