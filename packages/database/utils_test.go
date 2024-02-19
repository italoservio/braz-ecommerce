package database_test

import (
	"testing"

	"github.com/italoservio/braz_ecommerce/packages/database"
	"github.com/stretchr/testify/assert"
)

func TestUtils_ParseDocument(t *testing.T) {
	type SampleStruct struct {
		Name  string `bson:"name"`
		Value int    `bson:"value"`
	}

	t.Run("should parse structure when document is valid", func(t *testing.T) {
		structure := SampleStruct{
			Name:  "TestName",
			Value: 42,
		}

		_, err := database.ParseToDocument(structure)
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.Nil(t, err, "should not return error")
	})

	t.Run("should fail marshaling when inputted structure is invalid", func(t *testing.T) {
		wrongStructure := ""

		_, err := database.ParseToDocument(wrongStructure)
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return marshal error")
	})
}

func TestUtils_ParseToDatabaseId(t *testing.T) {
	t.Run("should parse string to object id when input is valid", func(t *testing.T) {
		objectIds, err := database.ParseToDatabaseId("5f3b3e3a7f2e1c1e9b3c9c8d", "5f3b3e3a7f2e1c1e9b3c9c8d")
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}

		assert.Nil(t, err, "should not return error")
		assert.Equal(t, 2, len(objectIds), "should return 2 object ids")
	})

	t.Run("should fail when input is invalid", func(t *testing.T) {
		objectIds, err := database.ParseToDatabaseId("invalid_id")
		if err == nil {
			t.Fail()
		}

		assert.NotNil(t, err, "should return error")
		assert.Nil(t, objectIds, "should not return object ids")
	})
}
