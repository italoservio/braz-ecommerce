package database_test

import (
	"fmt"
	"testing"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

type MockStructure struct {
	Id  string `bson:"_id"`
	Foo string `bson:"foo"`
}

func TestCrudRepository_GetById(t *testing.T) {
	rootMt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	rootMt.Run("should return the document when call database with success", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockId := primitive.NewObjectID()

		nestedMt.AddMockResponses(mtest.CreateCursorResponse(
			1,
			fmt.Sprintf("%s.%s", mockDbName, mockCollName),
			mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: mockId},
				{Key: "foo", Value: "bar"},
			},
		))
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := database.NewCrudRepository(mockDB)

		var result MockStructure

		err := crudRepository.GetById(mockCollName, mockId.Hex(), &result)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.Nil(t, err, "should not return error")
		assert.Equal(t, "bar", result.Foo, "should return the expected object by id")
	})

	rootMt.Run("should return error when no document is found", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockId := primitive.NewObjectID()

		nestedMt.AddMockResponses(mtest.CreateCursorResponse(
			0,
			fmt.Sprintf("%s.%s", mockDbName, mockCollName),
			mtest.FirstBatch,
		))

		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := database.NewCrudRepository(mockDB)

		var result MockStructure

		err := crudRepository.GetById(mockCollName, mockId.Hex(), &result)
		if err == nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.Equal(t, exception.CodeNotFound, err.Error(), "should return the expected error")
	})

	rootMt.Run("should return error when failed to call database", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockId := primitive.NewObjectID()

		nestedMt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := database.NewCrudRepository(mockDB)

		var result MockStructure

		err := crudRepository.GetById(mockCollName, mockId.Hex(), &result)
		if err == nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.Equal(t, exception.CodeDatabaseFailed, err.Error(), "should return database call error")
	})

	rootMt.Run("should return error when wrong object id is provided", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockWrongId := "something_wrong"

		var mockStructure MockStructure
		mockDB := &database.Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := database.NewCrudRepository(mockDB)

		err := crudRepository.GetById(mockCollName, mockWrongId, &mockStructure)
		if err == nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.Equal(t, exception.CodeValidationFailed, err.Error(), "should return object id error")
	})
}

func TestCrudRepository_DeleteById(t *testing.T) {
	rootMt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	rootMt.Run("should return nil when call database with success", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockId := primitive.NewObjectID()

		nestedMt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "nModified", Value: 1},
		})
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := database.NewCrudRepository(mockDB)

		err := crudRepository.DeleteById(mockCollName, mockId.Hex())
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.Nil(t, err, "should not return error")
	})

	rootMt.Run("should return error when failed to call database", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockId := primitive.NewObjectID()

		nestedMt.AddMockResponses(bson.D{
			{Key: "ok", Value: 0},
			{Key: "nModified", Value: 0},
		})
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := database.NewCrudRepository(mockDB)

		err := crudRepository.DeleteById(mockCollName, mockId.Hex())
		if err == nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.NotNil(t, err, "should return database call error")
	})

	rootMt.Run("should return error when wrong object id is provided", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockWrongId := "something_wrong"

		mockDB := &database.Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := database.NewCrudRepository(mockDB)

		err := crudRepository.DeleteById(mockCollName, mockWrongId)
		if err == nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.NotNil(t, err, "should return object id error")
	})
}

func TestCrudRepository_CreateOne(t *testing.T) {
	rootMt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	rootMt.Run("should return the inserted id when created with success", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockId := primitive.NewObjectID().Hex()

		nestedMt.AddMockResponses(mtest.CreateSuccessResponse())
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := database.NewCrudRepository(mockDB)

		id, err := crudRepository.CreateOne(
			mockCollName,
			MockStructure{Foo: "bar", Id: mockId},
		)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.Nil(t, err, "should not return error")
		assert.Equal(t, mockId, id, "should return the expected id")
	})

	rootMt.Run("should return error when failed to parse object id", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"

		nestedMt.AddMockResponses(mtest.CreateSuccessResponse())
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := database.NewCrudRepository(mockDB)

		_, err := crudRepository.CreateOne(
			mockCollName,
			MockStructure{Foo: "bar"},
		)
		if err == nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.NotNil(t, err, "should return error")
	})

	rootMt.Run("should return error when failed to call database", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"

		nestedMt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := database.NewCrudRepository(mockDB)

		_, err := crudRepository.CreateOne(
			mockCollName,
			MockStructure{Foo: "bar"},
		)
		if err == nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.NotNil(t, err, "should return error")
	})
}

func TestCrudRepository_UpdateById(t *testing.T) {
	rootMt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	rootMt.Run("should return nil when call database with success", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockId := primitive.NewObjectID().Hex()

		nestedMt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "nModified", Value: 1},
		})
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := database.NewCrudRepository(mockDB)

		err := crudRepository.UpdateById(
			mockCollName,
			mockId,
			MockStructure{Foo: "bar", Id: mockId},
		)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.Nil(t, err, "should not return error")
	})

	rootMt.Run("should return error when failed to call database", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockId := primitive.NewObjectID().Hex()

		nestedMt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := database.NewCrudRepository(mockDB)

		err := crudRepository.UpdateById(
			mockCollName,
			mockId,
			MockStructure{Foo: "bar", Id: mockId},
		)
		if err == nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.NotNil(t, err, "should return database call error")
	})

	rootMt.Run("should return error when wrong object id is provided", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockWrongId := "something_wrong"

		mockDB := &database.Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := database.NewCrudRepository(mockDB)

		err := crudRepository.UpdateById(
			mockCollName,
			mockWrongId,
			MockStructure{Foo: "bar", Id: mockWrongId},
		)
		if err == nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.NotNil(t, err, "should return object id error")
	})

	rootMt.Run("should return error when failed to parse struct to document", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockId := primitive.NewObjectID().Hex()

		mockDB := &database.Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := database.NewCrudRepository(mockDB)

		err := crudRepository.UpdateById(
			mockCollName,
			mockId,
			"something_wrong",
		)
		if err == nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.NotNil(t, err, "should return parse error")
	})
}
