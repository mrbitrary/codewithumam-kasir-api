package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIResponseWithID(t *testing.T) {
	resp := NewAPIResponse(map[string]string{"test": "data"})
	resp = resp.WithID("custom-id")

	assert.Equal(t, "custom-id", resp.ID)
}

func TestAPIResponseWithEtag(t *testing.T) {
	resp := NewAPIResponse(map[string]string{"test": "data"})
	resp = resp.WithEtag("custom-etag")

	assert.Equal(t, "custom-etag", resp.Etag)
}

func TestAPIResponseWithUpdated(t *testing.T) {
	resp := NewAPIResponse(map[string]string{"test": "data"})
	resp = resp.WithUpdated("2024-01-01T00:00:00Z")

	assert.Equal(t, "2024-01-01T00:00:00Z", resp.Updated)
}

func TestGenerateETagWithNilData(t *testing.T) {
	resp := NewAPIResponseWithItems(nil)

	assert.Empty(t, resp.Etag)
	assert.Nil(t, resp.Data)
}

func TestErrorDetailMarshalJSONWithMultipleErrors(t *testing.T) {
	errorDetail := ErrorDetail{
		Code: 400,
		Errors: []ErrorItem{
			{Message: "First error", Reason: ReasonRequired},
			{Message: "Second error", Reason: ReasonInvalidValue},
		},
	}

	data, err := json.Marshal(errorDetail)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// Should use first error message
	assert.Equal(t, "First error", result["message"])
}

func TestErrorDetailMarshalJSONWithNoErrors(t *testing.T) {
	errorDetail := ErrorDetail{
		Code:   500,
		Errors: []ErrorItem{},
	}

	data, err := json.Marshal(errorDetail)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// Should have empty message
	assert.Equal(t, "", result["message"])
}

func TestCategoryEntityToModelWithNilDeletedAt(t *testing.T) {
	entity := &CategoryEntity{
		Name:        "Test",
		Description: "Test Desc",
		DeletedAt:   nil,
	}

	model := entity.ToModel()

	require.NotNil(t, model)
	assert.Nil(t, model.DeletedAt)
}

func TestProductEntityToModelWithNilDeletedAt(t *testing.T) {
	entity := &ProductEntity{
		Name:      "Test Product",
		Price:     1000,
		Stocks:    10,
		DeletedAt: nil,
	}

	model := entity.ToModel()

	require.NotNil(t, model)
	assert.Nil(t, model.DeletedAt)
}

func TestCreateCategoryRequestToEntityGeneratesUUID(t *testing.T) {
	req := &CreateCategoryRequest{
		Name:        "Test",
		Description: "Test Desc",
	}

	entity1 := req.ToEntity()
	entity2 := req.ToEntity()

	// Should generate different UUIDs each time
	assert.NotEqual(t, entity1.ID, entity2.ID)
}

func TestCreateProductRequestToEntityGeneratesUUID(t *testing.T) {
	req := &CreateProductRequest{
		Name:     "Test",
		Price:    1000,
		Stocks:   10,
		Category: "Electronics",
	}

	entity1 := req.ToEntity()
	entity2 := req.ToEntity()

	// Should generate different UUIDs each time
	assert.NotEqual(t, entity1.ID, entity2.ID)
}

func TestJSONMarshalingCategory(t *testing.T) {
	category := Category{
		ID:          "test-id",
		Name:        "Electronics",
		Description: "Devices",
		Version:     1,
	}

	data, err := json.Marshal(category)
	require.NoError(t, err)

	var unmarshaled Category
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, category.ID, unmarshaled.ID)
	assert.Equal(t, category.Name, unmarshaled.Name)
	assert.Equal(t, category.Description, unmarshaled.Description)
}

func TestJSONMarshalingProduct(t *testing.T) {
	product := Product{
		ID:       "test-id",
		Name:     "Laptop",
		Price:    1000,
		Stocks:   5,
		Category: "Electronics",
		Version:  1,
	}

	data, err := json.Marshal(product)
	require.NoError(t, err)

	var unmarshaled Product
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, product.ID, unmarshaled.ID)
	assert.Equal(t, product.Name, unmarshaled.Name)
	assert.Equal(t, product.Price, unmarshaled.Price)
	assert.Equal(t, product.Stocks, unmarshaled.Stocks)
}
