package database_test

import (
	"testing"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/stretchr/testify/assert"
)

func TestPaginatedSlice_NewPaginatedSlice(t *testing.T) {
	type MockStructure struct {
		Foo  string
		Fizz string
	}

	t.Run("should encapsulate in page when a slice is provided", func(t *testing.T) {
		structures := []MockStructure{
			{Foo: "bar1", Fizz: "buzz2"},
			{Foo: "bar2", Fizz: "buzz2"},
		}

		paginated := database.NewPaginatedSlice[MockStructure](1, 10, &structures)

		assert.Equal(t, 1, paginated.Page, "should return the expected page")
		assert.Equal(t, 10, paginated.PerPage, "should return the expected per page")
		assert.Equal(t, 2, len(*paginated.Items), "should contain 2 items")
	})

	t.Run("should encapsulate in page when an empty slice is provided", func(t *testing.T) {
		structures := []MockStructure{}

		paginated := database.NewPaginatedSlice[MockStructure](1, 10, &structures)

		assert.Equal(t, 1, paginated.Page, "should return the expected page")
		assert.Equal(t, 10, paginated.PerPage, "should return the expected per page")
		assert.Equal(t, 0, len(*paginated.Items), "should contain 0 items")
	})
}
