package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/italoservio/braz_ecommerce/packages/exception"
	"github.com/italoservio/braz_ecommerce/services/users/domain"
)

type GetUserPaginatedInterface interface {
	Do(ctx context.Context, input *GetUserPaginatedInput) (*database.PaginatedSlice[GetUserPaginatedOutput], error)
}

type GetUserPaginatedImpl struct {
	crudRepository database.CrudRepositoryInterface
}

func NewGetUserPaginatedImpl(
	cr database.CrudRepositoryInterface,
) *GetUserPaginatedImpl {
	return &GetUserPaginatedImpl{
		crudRepository: cr,
	}
}

type GetUserPaginatedInput struct {
	Page    int
	PerPage int
	Emails  []string
	Ids     []string
	Deleted bool
}

type GetUserPaginatedOutput struct {
	*domain.UserDatabaseNoPassword `bson:",inline"`
}

func (gup *GetUserPaginatedImpl) Do(
	ctx context.Context,
	input *GetUserPaginatedInput,
) (*database.PaginatedSlice[GetUserPaginatedOutput], error) {
	sorting := mountSorting()
	projection := mountProjection()
	filters, err := mountFilters(input)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%v", filters)
	users := []GetUserPaginatedOutput{}

	err = gup.crudRepository.GetPaginated(
		ctx,
		database.UsersCollection,
		input.Page,
		input.PerPage,
		filters,
		projection,
		sorting,
		&users,
	)
	if err != nil {
		return nil, err
	}

	output := database.NewPaginatedSlice[GetUserPaginatedOutput](
		input.Page,
		input.PerPage,
		&users,
	)

	return output, nil
}

func mountFilters(input *GetUserPaginatedInput) (map[string]any, error) {
	filters := make(map[string]any)
	if len(input.Emails) > 0 {
		filters["email"] = input.Emails
	}

	if len(input.Ids) > 0 {
		ids, err := database.ParseToDatabaseId(input.Ids...)
		if err != nil {
			return nil, errors.New(exception.CodeValidationFailed)
		}

		filters["_id"] = ids
	}

	if !input.Deleted {
		filters["deleted_at"] = nil
	}

	return filters, nil
}

func mountProjection() map[string]int {
	projection := make(map[string]int)
	projection["password"] = 0
	projection["cipher_key"] = 0

	return projection
}

func mountSorting() map[string]int {
	sorting := make(map[string]int)
	sorting["created_at"] = -1

	return sorting
}
