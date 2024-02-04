package app_test

import (
	"errors"
	"testing"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/services/users/app"
	"github.com/italoservio/braz_ecommerce/services/users/mocks"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestDeleteUserById(t *testing.T) {
	t.Run("should return nil when deleted with success", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		cr := mocks.NewMockCrudRepositoryInterface(ctrl)
		ur := mocks.NewMockUserRepositoryInterface(ctrl)

		deleteUserByIdImpl := app.NewDeleteUserByIdImpl(cr, ur)

		id := primitive.NewObjectID().Hex()

		cr.EXPECT().
			DeleteById(database.UsersCollection, id).
			Times(1).
			Return(nil)

		err := deleteUserByIdImpl.Do(id)

		assert.Nil(t, err, "should return nil")
	})

	t.Run("should return the error when failed to delete", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		cr := mocks.NewMockCrudRepositoryInterface(ctrl)
		ur := mocks.NewMockUserRepositoryInterface(ctrl)

		deleteUserByIdImpl := app.NewDeleteUserByIdImpl(cr, ur)

		id := primitive.NewObjectID().Hex()

		cr.EXPECT().
			DeleteById(database.UsersCollection, id).
			Times(1).
			Return(errors.New(exception.CodeDatabaseFailed))

		err := deleteUserByIdImpl.Do(id)

		assert.NotNil(t, err, "should return the error")
	})
}
