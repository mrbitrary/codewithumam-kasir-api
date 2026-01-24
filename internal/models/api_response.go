package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"reflect"

	"github.com/google/uuid"
	"kasir-api/internal/utils"
)

// Reference: https://google.github.io/styleguide/jsoncstyleguide.xml
// Google JSON Style Guide defines standard structures for API responses

// APIResponse represents an API response following Google JSON Style
type APIResponse struct {
	// Top-level reserved properties
	Context string `json:"context,omitempty"` // Context for the request (e.g., session info)
	ID      string `json:"id,omitempty"`      // Unique identifier for this particular response (auto-generated UUIDv7)

	// Reserved data properties (common for successful responses)
	Kind    string `json:"kind,omitempty"`    // Resource type identifier (e.g., "kasir#product")
	Etag    string `json:"etag,omitempty"`    // Entity tag for caching and versioning (auto-generated from data)
	Updated string `json:"updated,omitempty"` // RFC3339 timestamp of last update

	// Data payload
	// For collections, Data should be an object with an "items" field, e.g.:
	// Data: map[string]any{"items": []Product{...}}
	Data any `json:"data,omitempty"` // Data payload (single item or object containing items array)

	// Error information (populated only if there's an error)
	Error *ErrorDetail `json:"error,omitempty"` // Detailed error information
}

// ErrorDetail contains detailed information about an error
type ErrorDetail struct {
	Code    int         `json:"code"`             // HTTP status code or application-specific error code
	Message string      `json:"message"`          // Human-readable error message (auto-populated from Errors[0].Message if empty)
	Errors  []ErrorItem `json:"errors,omitempty"` // Array of detailed error information
}

// MarshalJSON implements custom JSON marshaling for ErrorDetail
// It automatically populates Message from Errors[0].Message if Message is empty
func (e ErrorDetail) MarshalJSON() ([]byte, error) {
	type Alias ErrorDetail
	aux := struct {
		Alias
	}{
		Alias: (Alias)(e),
	}

	// Auto-populate Message from first error if Message is empty
	if aux.Message == "" && len(aux.Errors) > 0 {
		aux.Message = aux.Errors[0].Message
	}

	return json.Marshal(aux)
}

// ErrorItem represents a single error item with detailed information
type ErrorItem struct {
	Reason  string `json:"reason,omitempty"` // Short error reason (e.g., "invalidParameter", "required")
	Message string `json:"message"`          // Human-readable error message
}

// generateETag generates an ETag hash from the data
func generateETag(data any) string {
	if data == nil {
		return ""
	}

	// Marshal data to JSON for consistent hashing
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	// Generate SHA-256 hash
	hash := sha256.Sum256(jsonData)
	return hex.EncodeToString(hash[:])
}

// NewAPIResponse creates a new successful API response with a single data item
// Automatically generates ID (UUIDv7) and ETag (SHA-256 hash of data)
func NewAPIResponse(data any) *APIResponse {
	id, err := uuid.NewV7()
	if err != nil {
		return nil
	}

	return &APIResponse{
		ID:   utils.EncodeBase62(id.String()),
		Data: data,
		Etag: generateETag(data),
	}
}

// NewAPIResponseWithItems creates a new successful API response with a collection of items
// This is a convenience function that wraps items in a map with "items" key
// Automatically generates ID (UUIDv7) and ETag (SHA-256 hash of data)
// If items is nil or an empty slice, both "items" and "data" fields will be omitted from the response
func NewAPIResponseWithItems(items any) *APIResponse {
	id, err := uuid.NewV7()
	if err != nil {
		return nil
	}

	response := &APIResponse{
		ID: utils.EncodeBase62(id.String()),
	}

	// Only set Data and Etag if items is not nil and not empty
	if items != nil {
		v := reflect.ValueOf(items)
		if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
			// Only include if slice/array has elements
			if v.Len() >= 0 {
				dataMap := map[string]any{"items": items}
				response.Data = dataMap
				response.Etag = generateETag(dataMap)
			}
		} else {
			// Non-slice/array, include as-is
			dataMap := map[string]any{"items": items}
			response.Data = dataMap
			response.Etag = generateETag(dataMap)
		}
	}

	return response
}

// WithContext sets the context
func (r *APIResponse) WithContext(context string) *APIResponse {
	r.Context = context
	return r
}

// WithID sets the response ID (overrides auto-generated ID)
func (r *APIResponse) WithID(id string) *APIResponse {
	r.ID = id
	return r
}

// WithKind sets the resource kind
func (r *APIResponse) WithKind(kind string) *APIResponse {
	r.Kind = kind
	return r
}

// WithEtag sets the entity tag (overrides auto-generated ETag)
func (r *APIResponse) WithEtag(etag string) *APIResponse {
	r.Etag = etag
	return r
}

// WithUpdated sets the updated timestamp
func (r *APIResponse) WithUpdated(updated string) *APIResponse {
	r.Updated = updated
	return r
}

// NewAPIError creates a new API error response
// Automatically generates ID (UUIDv7)
// The message will be auto-populated from the first error item if not provided
func NewAPIError(code int, message string) *APIResponse {
	id, err := uuid.NewV7()
	if err != nil {
		return nil
	}

	return &APIResponse{
		ID: utils.EncodeBase62(id.String()),
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}

// NewAPIErrorWithErrors creates a new API error response with error items
// Automatically generates ID (UUIDv7)
// The top-level message will be auto-populated from the first error item
func NewAPIErrorWithErrors(code int, errors []ErrorItem) *APIResponse {
	id, err := uuid.NewV7()
	if err != nil {
		return nil
	}

	return &APIResponse{
		ID: utils.EncodeBase62(id.String()),
		Error: &ErrorDetail{
			Code:   code,
			Errors: errors,
			// Message will be auto-populated from Errors[0].Message during marshaling
		},
	}
}

// AddError adds a detailed error item to the error response
func (r *APIResponse) AddError(err ErrorItem) *APIResponse {
	if r.Error != nil {
		r.Error.Errors = append(r.Error.Errors, err)
	}
	return r
}

// AddErrors adds multiple detailed error items to the error response
func (r *APIResponse) AddErrors(errors []ErrorItem) *APIResponse {
	if r.Error != nil {
		r.Error.Errors = append(r.Error.Errors, errors...)
	}
	return r
}

// NewErrorItem creates a new error item
func NewErrorItem(message string) ErrorItem {
	return ErrorItem{
		Message: message,
	}
}

// WithReason sets the reason for the error item
func (ei ErrorItem) WithReason(reason string) ErrorItem {
	ei.Reason = reason
	return ei
}

// Common error reasons
const (
	ReasonInvalidParameter   = "invalidParameter"
	ReasonInvalidValue       = "invalidValue"
	ReasonRequired           = "required"
	ReasonNotFound           = "notFound"
	ReasonAlreadyExists      = "alreadyExists"
	ReasonUnauthorized       = "unauthorized"
	ReasonForbidden          = "forbidden"
	ReasonRateLimitExceeded  = "rateLimitExceeded"
	ReasonBackendError       = "backendError"
	ReasonServiceUnavailable = "serviceUnavailable"
)
