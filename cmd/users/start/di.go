package start

import (
	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/packages/encryption"
	"github.com/italoservio/braz_ecommerce/packages/logger"
	"github.com/italoservio/braz_ecommerce/services/users/app"
	"github.com/italoservio/braz_ecommerce/services/users/infra/http"
	"github.com/italoservio/braz_ecommerce/services/users/infra/storage"
)

func InjectionsContainer(db *database.Database) *http.UserControllerImpl {
	loggerImpl := logger.NewLogger()
	encryptionImpl := encryption.NewEncryptionImpl(loggerImpl)

	userRepositoryImpl := storage.NewUserRepositoryImpl(db)
	crudRepositoryImpl := database.NewCrudRepository(loggerImpl, db)
	getUserByIdImpl := app.NewGetUserByIdImpl(crudRepositoryImpl, userRepositoryImpl)
	deleteUserByIdImpl := app.NewDeleteUserByIdImpl(crudRepositoryImpl, userRepositoryImpl)
	createUserImpl := app.NewCreateUserImpl(encryptionImpl, crudRepositoryImpl, userRepositoryImpl)
	getUserPaginatedImpl := app.NewGetUserPaginatedImpl(crudRepositoryImpl)

	userControllerImpl := http.NewUserControllerImpl(
		loggerImpl,
		getUserByIdImpl,
		deleteUserByIdImpl,
		createUserImpl,
		getUserPaginatedImpl,
	)

	return userControllerImpl
}
