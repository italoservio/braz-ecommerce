package start

import (
	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/services/users/app"
	"github.com/italoservio/braz_ecommerce/services/users/infra/http"
	"github.com/italoservio/braz_ecommerce/services/users/infra/storage"
)

func InjectionsContainer(db *database.Database) *http.UserControllerImpl {
	userRepositoryImpl := storage.NewUserRepositoryImpl(db)
	crudRepositoryImpl := database.NewCrudRepository(db)
	getUserByIdImpl := app.NewGetUserByIdImpl(crudRepositoryImpl, userRepositoryImpl)
	deleteUserByIdImpl := app.NewDeleteUserByIdImpl(crudRepositoryImpl, userRepositoryImpl)
	createUserImpl := app.NewCreateUserImpl(crudRepositoryImpl, userRepositoryImpl)
	updateUserImpl := app.NewUpdateUserImpl(crudRepositoryImpl, userRepositoryImpl)
	userControllerImpl := http.NewUserControllerImpl(getUserByIdImpl, deleteUserByIdImpl, createUserImpl, updateUserImpl)

	return userControllerImpl
}
