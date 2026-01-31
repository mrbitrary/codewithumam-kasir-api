package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAPIResponse(t *testing.T) {
	data := map[string]string{"key": "value"}
	
	resp := NewAPIResponse(data)

	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.ID)
	assert.NotEmpty(t, resp.Etag)
	assert.Equal(t, data, resp.Data)
	assert.Nil(t, resp.Error)
}

func TestNewAPIResponseWithItems_WithData(t *testing.T) {
	items := []string{"item1", "item2"}
	
	resp := NewAPIResponseWithItems(items)

	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.ID)
	assert.NotEmpty(t, resp.Etag)
	
	dataMap, ok := resp.Data.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, items, dataMap["items"])
}

func TestNewAPIResponseWithItems_EmptySlice(t *testing.T) {
	items := []string{}
	
	resp := NewAPIResponseWithItems(items)

	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.ID)
	// Empty slice should still have data
	assert.NotNil(t, resp.Data)
}

func TestNewAPIResponseWithItems_Nil(t *testing.T) {
	resp := NewAPIResponseWithItems(nil)

	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.ID)
	assert.Nil(t, resp.Data)
	assert.Empty(t, resp.Etag)
}

func TestNewAPIError(t *testing.T) {
	resp := NewAPIError(404, "Not found")

	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.ID)
	require.NotNil(t, resp.Error)
	assert.Equal(t, 404, resp.Error.Code)
	assert.Equal(t, "Not found", resp.Error.Message)
	assert.Nil(t, resp.Data)
}

func TestNewAPIErrorWithErrors(t *testing.T) {
	errors := []ErrorItem{
		{Message: "Field is required", Reason: ReasonRequired},
		{Message: "Invalid value", Reason: ReasonInvalidValue},
	}
	
	resp := NewAPIErrorWithErrors(400, errors)

	require.NotNil(t, resp)
	require.NotNil(t, resp.Error)
	assert.Equal(t, 400, resp.Error.Code)
	assert.Len(t, resp.Error.Errors, 2)
}

func TestErrorDetail_MarshalJSON_AutoPopulateMessage(t *testing.T) {
	errorDetail := ErrorDetail{
		Code: 400,
		Errors: []ErrorItem{
			{Message: "First error", Reason: ReasonRequired},
		},
	}

	data, err := json.Marshal(errorDetail)
	require.NoError(t, err)

	var result map[string]any
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)
	
	assert.Equal(t, "First error", result["message"])
}

func TestAPIResponse_WithContext(t *testing.T) {
	resp := NewAPIResponse(map[string]string{"test": "data"})
	resp = resp.WithContext("test-context")

	assert.Equal(t, "test-context", resp.Context)
}

func TestAPIResponse_WithKind(t *testing.T) {
	resp := NewAPIResponse(map[string]string{"test": "data"})
	resp = resp.WithKind("kasir#product")

	assert.Equal(t, "kasir#product", resp.Kind)
}

func TestNewErrorItem(t *testing.T) {
	item := NewErrorItem("Test error")

	assert.Equal(t, "Test error", item.Message)
	assert.Empty(t, item.Reason)
}

func TestErrorItem_WithReason(t *testing.T) {
	item := NewErrorItem("Test error").WithReason(ReasonNotFound)

	assert.Equal(t, "Test error", item.Message)
	assert.Equal(t, ReasonNotFound, item.Reason)
}

func TestGenerateETag_Consistency(t *testing.T) {
	data := map[string]string{"key": "value"}
	
	resp1 := NewAPIResponse(data)
	resp2 := NewAPIResponse(data)

	assert.Equal(t, resp1.Etag, resp2.Etag, "Same data should generate same ETag")
}

func TestGenerateETag_DifferentData(t *testing.T) {
	resp1 := NewAPIResponse(map[string]string{"key": "value1"})
	resp2 := NewAPIResponse(map[string]string{"key": "value2"})

	assert.NotEqual(t, resp1.Etag, resp2.Etag, "Different data should generate different ETags")
}
