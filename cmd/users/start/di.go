package start

import (
	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/services/users/app"
	"github.com/italoservio/braz_ecommerce/services/users/infra/http"
	"github.com/italoservio/braz_ecommerce/services/users/infra/storage"
)

func InjectionsContainer(database *database.Database) *http.UserControllerImpl {
	userRepositoryImpl := storage.NewUserRepositoryImpl(database)
	getUserByIdImpl := app.NewGetUserByIdImpl(userRepositoryImpl)
	userControllerImpl := http.NewUserControllerImpl(getUserByIdImpl)

	return userControllerImpl
}
