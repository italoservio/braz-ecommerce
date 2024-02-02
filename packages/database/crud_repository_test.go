package database

import (
	"fmt"
	"testing"

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

		mockDB := &Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := NewCrudRepository(mockDB)

		var result MockStructure

		decoder, err := crudRepository.GetById(mockCollName, mockId.Hex())
		if err != nil {
			t.Fatal(err)
		}

		err = decoder.Decode(&result)

		assert.Nil(t, err, "should not return error")
		assert.Equal(t, result.Foo, "bar", "should return the expected object by id")
	})

	rootMt.Run("should return error when failed to call database", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockId := primitive.NewObjectID()

		nestedMt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})
		defer nestedMt.ClearMockResponses()

		mockDB := &Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := NewCrudRepository(mockDB)

		var result MockStructure

		decoder, err := crudRepository.GetById(mockCollName, mockId.Hex())
		if err != nil {
			t.Fatal(err)
		}

		err = decoder.Decode(&result)

		assert.NotNil(t, err, "should return database call error")
	})

	rootMt.Run("should return error when wrong object id is provided", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockWrongId := "something_wrong"

		mockDB := &Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := NewCrudRepository(mockDB)

		_, err := crudRepository.GetById(mockCollName, mockWrongId)
		if err == nil {
			t.Fatal(err)
		}

		assert.NotNil(t, err, "should return object id error")
	})
}

func TestCrudRepository_DeleteById(t *testing.T) {
	rootMt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	rootMt.Run("should return nil when call database with success", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockId := primitive.NewObjectID()

		nestedMt.AddMockResponses(bson.D{{Key: "ok", Value: 1}})
		defer nestedMt.ClearMockResponses()

		mockDB := &Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := NewCrudRepository(mockDB)

		err := crudRepository.DeleteById(mockCollName, mockId.Hex())
		if err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, err, "should not return error")
	})

	rootMt.Run("should return error when failed to call database", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockId := primitive.NewObjectID()

		nestedMt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})
		defer nestedMt.ClearMockResponses()

		mockDB := &Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := NewCrudRepository(mockDB)

		err := crudRepository.DeleteById(mockCollName, mockId.Hex())
		if err == nil {
			t.Fatal(err)
		}

		assert.NotNil(t, err, "should return database call error")
	})

	rootMt.Run("should return error when wrong object id is provided", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockWrongId := "something_wrong"

		mockDB := &Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := NewCrudRepository(mockDB)

		err := crudRepository.DeleteById(mockCollName, mockWrongId)
		if err == nil {
			t.Fatal(err)
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

		mockDB := &Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := NewCrudRepository(mockDB)

		id, err := crudRepository.CreateOne(
			mockCollName,
			MockStructure{Foo: "bar", Id: mockId},
		)
		if err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, err, "should not return error")
		assert.Equal(t, id, mockId, "should return the expected id")
	})

	rootMt.Run("should return error when failed to parse object id", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"

		nestedMt.AddMockResponses(mtest.CreateSuccessResponse())
		defer nestedMt.ClearMockResponses()

		mockDB := &Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := NewCrudRepository(mockDB)

		_, err := crudRepository.CreateOne(
			mockCollName,
			MockStructure{Foo: "bar"},
		)
		if err == nil {
			t.Fatal(err)
		}

		assert.NotNil(t, err, "should return error")
	})

	rootMt.Run("should return error when failed to call database", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"

		nestedMt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})
		defer nestedMt.ClearMockResponses()

		mockDB := &Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := NewCrudRepository(mockDB)

		_, err := crudRepository.CreateOne(
			mockCollName,
			MockStructure{Foo: "bar"},
		)
		if err == nil {
			t.Fatal(err)
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

		mockDB := &Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := NewCrudRepository(mockDB)

		err := crudRepository.UpdateById(
			mockCollName,
			mockId,
			MockStructure{Foo: "bar", Id: mockId},
		)
		if err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, err, "should not return error")
	})

	rootMt.Run("should return error when failed to call database", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockId := primitive.NewObjectID().Hex()

		nestedMt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})
		defer nestedMt.ClearMockResponses()

		mockDB := &Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := NewCrudRepository(mockDB)

		err := crudRepository.UpdateById(
			mockCollName,
			mockId,
			MockStructure{Foo: "bar", Id: mockId},
		)
		if err == nil {
			t.Fatal(err)
		}

		assert.NotNil(t, err, "should return database call error")
	})

	rootMt.Run("should return error when wrong object id is provided", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockWrongId := "something_wrong"

		mockDB := &Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := NewCrudRepository(mockDB)

		err := crudRepository.UpdateById(
			mockCollName,
			mockWrongId,
			MockStructure{Foo: "bar", Id: mockWrongId},
		)
		if err == nil {
			t.Fatal(err)
		}

		assert.NotNil(t, err, "should return object id error")
	})

	rootMt.Run("should return error when failed to parse struct to document", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockCollName := "users"
		mockId := primitive.NewObjectID().Hex()

		mockDB := &Database{nestedMt.Client.Database(mockDbName)}
		crudRepository := NewCrudRepository(mockDB)

		err := crudRepository.UpdateById(
			mockCollName,
			mockId,
			"something_wrong",
		)
		if err == nil {
			t.Fatal(err)
		}

		assert.NotNil(t, err, "should return parse error")
	})
}
