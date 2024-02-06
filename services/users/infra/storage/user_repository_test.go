package storage_test

import (
	"fmt"
	"testing"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/services/users/infra/storage"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestUserRepository_NewUserRepository(t *testing.T) {
	rootMt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	rootMt.Run("should return a new instance when all right", func(nestedMt *mtest.T) {
		mockDbName := "foo"
		mockDB := &database.Database{Database: nestedMt.Client.Database(mockDbName)}

		instance := storage.NewUserRepositoryImpl(mockDB)

		assert.Equal(
			t,
			fmt.Sprintf("%T", instance),
			"*storage.UserRepositoryImpl",
			"should be a pointer to UserRepositoryImpl",
		)
	})
}
