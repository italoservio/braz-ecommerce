package app_test

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/services/users/app"
	"github.com/italoservio/braz_ecommerce/services/users/domain"
	"github.com/italoservio/braz_ecommerce/services/users/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

type TestingDependencies_TestGetUserPaginated struct {
	ctx                  context.Context
	ctrl                 *gomock.Controller
	mockCrudRepository   *mocks.MockCrudRepositoryInterface
	getUserPaginatedImpl *app.GetUserPaginatedImpl
}

func BeforeEach_TestGetUserPaginated(t *testing.T) *TestingDependencies_TestGetUserPaginated {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	mockCrudRepository := mocks.NewMockCrudRepositoryInterface(ctrl)

	getUserPaginatedImpl := app.NewGetUserPaginatedImpl(mockCrudRepository)

	return &TestingDependencies_TestGetUserPaginated{
		ctx:                  ctx,
		ctrl:                 ctrl,
		mockCrudRepository:   mockCrudRepository,
		getUserPaginatedImpl: getUserPaginatedImpl,
	}
}

func TestGetUserPaginated_Do(t *testing.T) {
	t.Run("should return error when failed to call database", func(t *testing.T) {
		deps := BeforeEach_TestGetUserPaginated(t)

		mockExpectedError := errors.New("something goes wrong")
		input := &app.GetUserPaginatedInput{}

		filters := map[string]any{}
		projection := map[string]int{
			"password":   0,
			"cipher_key": 0,
		}
		sorting := map[string]int{
			"created_at": -1,
		}

		deps.mockCrudRepository.
			EXPECT().
			GetPaginated(gomock.Any(), database.UsersCollection, input.Page, input.PerPage, filters, projection, sorting, gomock.Any()).
			Times(1).
			Return(mockExpectedError)

		_, err := deps.getUserPaginatedImpl.Do(deps.ctx, true, input)
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return error")
	})

	t.Run("should return empty error when executed successfully", func(t *testing.T) {
		deps := BeforeEach_TestGetUserPaginated(t)

		objId1 := primitive.NewObjectID()
		objId2 := primitive.NewObjectID()

		input := &app.GetUserPaginatedInput{
			Page:    1,
			PerPage: 10,
			Ids:     []string{objId1.Hex(), objId2.Hex()},
			Emails:  []string{"foo@bar.net"},
		}

		filters := map[string]any{
			"_id":   []primitive.ObjectID{objId1, objId2},
			"email": []string{"foo@bar.net"},
		}
		projection := map[string]int{
			"password":   0,
			"cipher_key": 0,
		}
		sorting := map[string]int{
			"created_at": -1,
		}

		deps.mockCrudRepository.
			EXPECT().
			GetPaginated(gomock.Any(), database.UsersCollection, input.Page, input.PerPage, filters, projection, sorting, gomock.Any()).
			Times(1).
			DoAndReturn(func(
				ctx context.Context,
				collection string,
				page int,
				perPage int,
				filters map[string]any,
				projection map[string]int,
				sorting map[string]int,
				structures any,
			) error {
				*structures.(*[]app.GetUserPaginatedOutput) = []app.GetUserPaginatedOutput{
					{UserDatabaseNoPassword: &domain.UserDatabaseNoPassword{
						DatabaseIdentifier: &database.DatabaseIdentifier{
							Id: "123",
						},
					}},
				}

				return nil
			})

		structures, err := deps.getUserPaginatedImpl.Do(deps.ctx, true, input)
		if err != nil {
			log.Fatal(err)
		}

		assert.Nil(t, err, "should not return an error")
		assert.Equal(t, 1, len(*structures.Items), "should return one structure")
	})

	t.Run("should return empty error when executed successfully", func(t *testing.T) {
		deps := BeforeEach_TestGetUserPaginated(t)

		objId1 := primitive.NewObjectID()
		objId2 := primitive.NewObjectID()

		input := &app.GetUserPaginatedInput{
			Page:    1,
			PerPage: 10,
			Ids:     []string{objId1.Hex(), objId2.Hex()},
			Emails:  []string{"foo@bar.net"},
		}

		filters := map[string]any{
			"_id":   []primitive.ObjectID{objId1, objId2},
			"email": []string{"foo@bar.net"},
		}
		projection := map[string]int{
			"password":   0,
			"cipher_key": 0,
		}
		sorting := map[string]int{
			"created_at": -1,
		}

		deps.mockCrudRepository.
			EXPECT().
			GetPaginated(gomock.Any(), database.UsersCollection, input.Page, input.PerPage, filters, projection, sorting, gomock.Any()).
			Times(1).
			DoAndReturn(func(
				ctx context.Context,
				collection string,
				page int,
				perPage int,
				filters map[string]any,
				projection map[string]int,
				sorting map[string]int,
				structures any,
			) error {
				*structures.(*[]app.GetUserPaginatedOutput) = []app.GetUserPaginatedOutput{
					{UserDatabaseNoPassword: &domain.UserDatabaseNoPassword{
						DatabaseIdentifier: &database.DatabaseIdentifier{
							Id: "123",
						},
					}},
				}

				return nil
			})

		structures, err := deps.getUserPaginatedImpl.Do(deps.ctx, false, input)
		if err != nil {
			log.Fatal(err)
		}

		assert.Nil(t, err, "should not return an error")
		assert.Equal(t, 1, len(*structures.Items), "should return one structure")
	})

	t.Run("should return error when the filter id isn't a valid database id", func(t *testing.T) {
		deps := BeforeEach_TestGetUserPaginated(t)

		input := &app.GetUserPaginatedInput{
			Page:    1,
			PerPage: 10,
			Ids:     []string{"foo"},
		}

		_, err := deps.getUserPaginatedImpl.Do(deps.ctx, true, input)
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return error")
	})
}
