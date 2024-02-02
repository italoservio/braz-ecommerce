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
	getUserByIdImpl := app.NewGetUserByIdImpl(userRepositoryImpl, crudRepositoryImpl)
	userControllerImpl := http.NewUserControllerImpl(getUserByIdImpl)

	return userControllerImpl
}
