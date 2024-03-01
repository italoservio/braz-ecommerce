package database_test

import (
	"context"
	"testing"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/packages/logger"
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

func TestCrudRepository_GetById(t *testing.T) {
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
		crudRepository := database.NewCrudRepository(logger, mockDB)

		var result MockStructure

		err := crudRepository.GetById(ctx, MOCK_COLL_NAME, "mockId.Hex()", &result)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.Nil(t, err, "should not return error")
		assert.Equal(t, "bar", result.Foo, "should return the expected object by id")
	})

	rootMt.Run("should return error when no document is found", func(nestedMt *mtest.T) {
		mockId := primitive.NewObjectID()

		nestedMt.AddMockResponses(mtest.CreateCursorResponse(
			0,
			MOCK_NS,
			mtest.FirstBatch,
		))

		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(MOCK_DB_NAME)}
		crudRepository := database.NewCrudRepository(logger, mockDB)

		var result MockStructure

		err := crudRepository.GetById(ctx, MOCK_COLL_NAME, mockId.Hex(), &result)
		if err == nil {
			t.Fail()
		}

		assert.Equal(t, exception.CodeNotFound, err.Error(), "should return the expected error")
	})

	rootMt.Run("should return error when failed to call database", func(nestedMt *mtest.T) {

		mockId := primitive.NewObjectID()

		nestedMt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(MOCK_DB_NAME)}
		crudRepository := database.NewCrudRepository(logger, mockDB)

		var result MockStructure

		err := crudRepository.GetById(ctx, MOCK_COLL_NAME, mockId.Hex(), &result)
		if err == nil {
			t.Fail()
		}

		assert.Equal(t, exception.CodeDatabaseFailed, err.Error(), "should return database call error")
	})

	rootMt.Run("should return error when wrong object id is provided", func(nestedMt *mtest.T) {

		mockWrongId := "something_wrong"

		var mockStructure MockStructure
		mockDB := &database.Database{nestedMt.Client.Database(MOCK_DB_NAME)}
		crudRepository := database.NewCrudRepository(logger, mockDB)

		err := crudRepository.GetById(ctx, MOCK_COLL_NAME, mockWrongId, &mockStructure)
		if err == nil {
			t.Fail()
		}

		assert.Equal(t, exception.CodeValidationFailed, err.Error(), "should return object id error")
	})
}

func TestCrudRepository_DeleteById(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewLogger()
	rootMt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	rootMt.Run("should return nil when call database with success", func(nestedMt *mtest.T) {

		mockId := primitive.NewObjectID()

		nestedMt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "nModified", Value: 1},
		})
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(MOCK_DB_NAME)}
		crudRepository := database.NewCrudRepository(logger, mockDB)

		err := crudRepository.DeleteById(ctx, MOCK_COLL_NAME, mockId.Hex())
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.Nil(t, err, "should not return error")
	})

	rootMt.Run("should return error when failed to call database", func(nestedMt *mtest.T) {

		mockId := primitive.NewObjectID()

		nestedMt.AddMockResponses(bson.D{
			{Key: "ok", Value: 0},
			{Key: "nModified", Value: 0},
		})
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(MOCK_DB_NAME)}
		crudRepository := database.NewCrudRepository(logger, mockDB)

		err := crudRepository.DeleteById(ctx, MOCK_COLL_NAME, mockId.Hex())
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return database call error")
	})

	rootMt.Run("should return error when wrong object id is provided", func(nestedMt *mtest.T) {

		mockWrongId := "something_wrong"

		mockDB := &database.Database{nestedMt.Client.Database(MOCK_DB_NAME)}
		crudRepository := database.NewCrudRepository(logger, mockDB)

		err := crudRepository.DeleteById(ctx, MOCK_COLL_NAME, mockWrongId)
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return object id error")
	})
}

func TestCrudRepository_CreateOne(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewLogger()
	rootMt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	type MockStructureB struct {
		Id  *primitive.ObjectID `bson:"_id"`
		Foo string              `bson:"foo"`
	}

	rootMt.Run("should return the inserted id when created with success", func(nestedMt *mtest.T) {

		mockId := primitive.NewObjectID()

		nestedMt.AddMockResponses(mtest.CreateSuccessResponse())
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(MOCK_DB_NAME)}
		crudRepository := database.NewCrudRepository(logger, mockDB)

		id, err := crudRepository.CreateOne(
			ctx,
			MOCK_COLL_NAME,
			MockStructureB{Foo: "bar", Id: &mockId},
		)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.Nil(t, err, "should not return error")
		assert.Equal(t, mockId.Hex(), id, "should return the expected id")
	})

	rootMt.Run("should return error when failed to call database", func(nestedMt *mtest.T) {

		nestedMt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(MOCK_DB_NAME)}
		crudRepository := database.NewCrudRepository(logger, mockDB)

		_, err := crudRepository.CreateOne(
			ctx,
			MOCK_COLL_NAME,
			MockStructure{Foo: "bar"},
		)
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return error")
	})
}

func TestCrudRepository_UpdateById(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewLogger()
	rootMt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	rootMt.Run("should return error when failed to call database", func(nestedMt *mtest.T) {

		mockId := primitive.NewObjectID().Hex()

		nestedMt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(MOCK_DB_NAME)}
		crudRepository := database.NewCrudRepository(logger, mockDB)

		err := crudRepository.UpdateById(
			ctx,
			MOCK_COLL_NAME,
			mockId,
			MockStructure{Foo: "bar", Id: mockId},
			MockStructure{Foo: "bar", Id: mockId},
		)
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return database call error")
	})

	rootMt.Run("should return error when wrong object id is provided", func(nestedMt *mtest.T) {

		mockWrongId := "something_wrong"

		mockDB := &database.Database{nestedMt.Client.Database(MOCK_DB_NAME)}
		crudRepository := database.NewCrudRepository(logger, mockDB)

		err := crudRepository.UpdateById(
			ctx,
			MOCK_COLL_NAME,
			mockWrongId,
			MockStructure{Foo: "bar", Id: mockWrongId},
			MockStructure{Foo: "bar", Id: mockWrongId},
		)
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return object id error")
	})

	rootMt.Run("should return error when failed to parse struct to document", func(nestedMt *mtest.T) {

		mockId := primitive.NewObjectID().Hex()

		mockDB := &database.Database{nestedMt.Client.Database(MOCK_DB_NAME)}
		crudRepository := database.NewCrudRepository(logger, mockDB)

		err := crudRepository.UpdateById(
			ctx,
			MOCK_COLL_NAME,
			mockId,
			"something_wrong",
			"something_wrong",
		)
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return parse error")
	})

	rootMt.Run("testesteste", func(nestedMt *mtest.T) {
		mockId := primitive.NewObjectID()
		nestedMt.AddMockResponses(mtest.CreateCursorResponse(1, MOCK_NS, MOCK_DB_NAME, bson.D{
			{"id", mockId.Hex()},
			{"field-1", MOCK_DB_NAME},
			{"field-2", MOCK_COLL_NAME},
		}))

		mockDB := &database.Database{nestedMt.Client.Database(MOCK_DB_NAME)}
		crudRepository := database.NewCrudRepository(logger, mockDB)

		var input MockStructure
		var output MockStructure
		err := crudRepository.UpdateById(ctx, MOCK_COLL_NAME, mockId.Hex(), input, output)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.Nil(t, err, "should not return error")
	})
}

func TestCrudRepository_GetPaginated(t *testing.T) {
	ctx := context.TODO()
	logger := logger.NewLogger()
	rootMt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	rootMt.Run("should fill the array when call database with success", func(nestedMt *mtest.T) {

		mockId1 := primitive.NewObjectID().Hex()
		mockId2 := primitive.NewObjectID().Hex()
		page := 1
		perPage := 10
		filters := map[string]interface{}{"foo": []string{"bar", "buzz"}}
		projections := map[string]int{"foo": 1}
		sort := map[string]int{"foo": 1}
		structures := []MockStructure{}

		ns := MOCK_NS
		nestedMt.AddMockResponses()
		nestedMt.AddMockResponses(mtest.CreateCursorResponse(
			1,
			ns,
			mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: mockId1},
				{Key: "foo", Value: "bar"},
			},
		), mtest.CreateCursorResponse(
			1,
			ns,
			mtest.NextBatch,
			bson.D{
				{Key: "_id", Value: mockId2},
				{Key: "foo", Value: "bar"},
			},
		), mtest.CreateCursorResponse(
			0,
			ns,
			mtest.NextBatch,
		))
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(MOCK_DB_NAME)}
		crudRepository := database.NewCrudRepository(logger, mockDB)

		err := crudRepository.GetPaginated(
			ctx,
			MOCK_COLL_NAME,
			page,
			perPage,
			filters,
			projections,
			sort,
			&structures,
		)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.Equal(t, 2, len(structures), "should return the expected number of items")
		assert.Nil(t, err, "should not return error")
	})

	rootMt.Run("should fill the array when no filters are provided and call database with success", func(nestedMt *mtest.T) {

		mockId1 := primitive.NewObjectID().Hex()
		mockId2 := primitive.NewObjectID().Hex()
		page := 1
		perPage := 10
		filters := map[string]interface{}{}
		projections := map[string]int{}
		sort := map[string]int{"foo": 1}
		structures := []MockStructure{}

		ns := MOCK_NS
		nestedMt.AddMockResponses(mtest.CreateCursorResponse(
			1,
			ns,
			mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: mockId1},
				{Key: "foo", Value: "bar"},
			},
		), mtest.CreateCursorResponse(
			1,
			ns,
			mtest.NextBatch,
			bson.D{
				{Key: "_id", Value: mockId2},
				{Key: "foo", Value: "bar"},
			},
		), mtest.CreateCursorResponse(
			0,
			ns,
			mtest.NextBatch,
		))
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(MOCK_DB_NAME)}
		crudRepository := database.NewCrudRepository(logger, mockDB)

		err := crudRepository.GetPaginated(
			ctx,
			MOCK_COLL_NAME,
			page,
			perPage,
			filters,
			projections,
			sort,
			&structures,
		)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.Equal(t, 2, len(structures), "should return the expected number of items")
		assert.Nil(t, err, "should not return error")
	})

	rootMt.Run("should return error when failed to call database", func(nestedMt *mtest.T) {

		page := 1
		perPage := 10
		filters := map[string]interface{}{}
		projections := map[string]int{}
		sort := map[string]int{}
		structures := []MockStructure{}

		nestedMt.AddMockResponses(bson.D{{Key: "ok", Value: 0}})
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(MOCK_DB_NAME)}
		crudRepository := database.NewCrudRepository(logger, mockDB)

		err := crudRepository.GetPaginated(
			ctx,
			MOCK_COLL_NAME,
			page,
			perPage,
			filters,
			projections,
			sort,
			&structures,
		)
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return database call error")
	})

	rootMt.Run("should return error when failed to fill the array", func(nestedMt *mtest.T) {

		mockId1 := primitive.NewObjectID().Hex()
		mockId2 := primitive.NewObjectID().Hex()
		page := 1
		perPage := 10
		filters := map[string]interface{}{}
		projections := map[string]int{}
		sort := map[string]int{}
		structures := "something_wrong"

		ns := MOCK_NS
		nestedMt.AddMockResponses(mtest.CreateCursorResponse(
			1,
			ns,
			mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: mockId1},
				{Key: "foo", Value: "bar"},
			},
		), mtest.CreateCursorResponse(
			1,
			ns,
			mtest.NextBatch,
			bson.D{
				{Key: "_id", Value: mockId2},
				{Key: "foo", Value: "bar"},
			},
		), mtest.CreateCursorResponse(
			0,
			ns,
			mtest.NextBatch,
		))
		defer nestedMt.ClearMockResponses()

		mockDB := &database.Database{nestedMt.Client.Database(MOCK_DB_NAME)}
		crudRepository := database.NewCrudRepository(logger, mockDB)

		err := crudRepository.GetPaginated(
			ctx,
			MOCK_COLL_NAME,
			page,
			perPage,
			filters,
			projections,
			sort,
			structures,
		)
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return fill error")
	})
}
