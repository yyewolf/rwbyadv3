package api

import (
	"net/http"

	"github.com/astaxie/beego/utils/pagination"
	"github.com/labstack/echo/v4"
)

// ErrorCode represents standard API error codes
type ErrorCode string

const (
	ErrorBadRequest          ErrorCode = "BAD_REQUEST"
	ErrorUnauthorized        ErrorCode = "UNAUTHORIZED"
	ErrorForbidden           ErrorCode = "FORBIDDEN"
	ErrorNotFound            ErrorCode = "NOT_FOUND"
	ErrorConflict            ErrorCode = "CONFLICT"
	ErrorInternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrorValidation          ErrorCode = "VALIDATION_ERROR"
	ErrorRedirected          ErrorCode = "REDIRECTED"
)

// Error represents a standardized API error
type Error struct {
	Code     ErrorCode `json:"code"`
	Message  string    `json:"message"`
	Details  any       `json:"details,omitempty"`
	Redirect string    `json:"redirect,omitempty"` // Optional redirect URL for errors that require user action
}

// ResponseMeta contains metadata about the response
type ResponseMeta struct {
	Success bool   `json:"success"`
	TraceID string `json:"trace_id,omitempty"`
}

// PaginationMeta contains pagination information
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
}

// Response is a generic API response structure
type Response[T any] struct {
	Meta       ResponseMeta    `json:"meta"`
	Data       T               `json:"data"`
	Error      *Error          `json:"error,omitempty"`
	Pagination *PaginationMeta `json:"pagination,omitempty"`
}

// NewSuccessResponse creates a new successful response with data
func NewSuccessResponse[T any](data T) Response[T] {
	return Response[T]{
		Meta: ResponseMeta{
			Success: true,
		},
		Data: data,
	}
}

// NewPaginatedResponse creates a new successful paginated response
func NewPaginatedResponse[T any](data T, paginator *pagination.Paginator) Response[T] {
	resp := NewSuccessResponse(data)

	if paginator != nil {
		resp.Pagination = &PaginationMeta{
			CurrentPage: paginator.Page(),
			TotalPages:  paginator.PageNums(),
			PerPage:     paginator.PerPageNums,
			TotalItems:  paginator.Nums(),
			HasNext:     paginator.HasNext(),
			HasPrev:     paginator.HasPrev(),
		}
	}

	return resp
}

// NewErrorResponse creates a new error response
func NewErrorResponse[T any](code ErrorCode, message string) Response[T] {
	return Response[T]{
		Meta: ResponseMeta{
			Success: false,
		},
		Error: &Error{
			Code:    code,
			Message: message,
		},
	}
}

// NewRedirectErrorResponse creates a new error response with a redirect URL
func NewRedirectErrorResponse[T any](code ErrorCode, message, redirect string) Response[T] {
	return Response[T]{
		Meta: ResponseMeta{
			Success: false,
		},
		Error: &Error{
			Code:     code,
			Message:  message,
			Redirect: redirect,
		},
	}
}

// WithTraceID adds a trace ID to the response
func (r Response[T]) WithTraceID(traceID string) Response[T] {
	r.Meta.TraceID = traceID
	return r
}

// WithErrorDetails adds details to an error response
func (r Response[T]) WithErrorDetails(details any) Response[T] {
	if r.Error != nil {
		r.Error.Details = details
	}
	return r
}

// JSON sends the response as JSON using Echo's context
func (r Response[T]) JSON(c echo.Context, statusCode int) error {
	return c.JSON(statusCode, r)
}

// Common HTTP status code helpers

// SendOK sends a 200 OK response
func SendOK[T any](c echo.Context, data T) error {
	return NewSuccessResponse(data).JSON(c, http.StatusOK)
}

// SendCreated sends a 201 Created response
func SendCreated[T any](c echo.Context, data T) error {
	return NewSuccessResponse(data).JSON(c, http.StatusCreated)
}

// SendPaginated sends a paginated response
func SendPaginated[T any](c echo.Context, data T, paginator *pagination.Paginator) error {
	return NewPaginatedResponse(data, paginator).JSON(c, http.StatusOK)
}

// SendBadRequest sends a 400 Bad Request error response
func SendBadRequest(c echo.Context, message string) error {
	return NewErrorResponse[any](ErrorBadRequest, message).JSON(c, http.StatusBadRequest)
}

// SendUnauthorized sends a 401 Unauthorized error response
func SendUnauthorized(c echo.Context, message string) error {
	if message == "" {
		message = "Authentication required"
	}
	return NewErrorResponse[any](ErrorUnauthorized, message).JSON(c, http.StatusUnauthorized)
}

// SendForbidden sends a 403 Forbidden error response
func SendForbidden(c echo.Context, message string) error {
	if message == "" {
		message = "You don't have permission to access this resource"
	}
	return NewErrorResponse[any](ErrorForbidden, message).JSON(c, http.StatusForbidden)
}

// SendNotFound sends a 404 Not Found error response
func SendNotFound(c echo.Context, message string) error {
	if message == "" {
		message = "Resource not found"
	}
	return NewErrorResponse[any](ErrorNotFound, message).JSON(c, http.StatusNotFound)
}

// SendConflict sends a 409 Conflict error response
func SendConflict(c echo.Context, message string) error {
	return NewErrorResponse[any](ErrorConflict, message).JSON(c, http.StatusConflict)
}

// SendInternalError sends a 500 Internal Server Error response
func SendInternalError(c echo.Context, message string) error {
	if message == "" {
		message = "Internal server error"
	}
	return NewErrorResponse[any](ErrorInternalServerError, message).JSON(c, http.StatusInternalServerError)
}

// OkRedirect sends a redirect response with a 200 OK status
func OkRedirect(c echo.Context, redirectURL string) error {
	return NewRedirectErrorResponse[any](ErrorRedirected, "Redirecting to login", redirectURL).
		JSON(c, http.StatusOK)
}

// SendRedirect sends a redirect response with a 401 Unauthorized status
func SendRedirect(c echo.Context, redirectURL string) error {
	return NewRedirectErrorResponse[any](ErrorUnauthorized, "Redirecting to login", redirectURL).
		JSON(c, http.StatusUnauthorized)
}

// SendValidationError sends a 422 Unprocessable Entity error response with validation details
func SendValidationError(c echo.Context, message string, details any) error {
	return NewErrorResponse[any](ErrorValidation, message).
		WithErrorDetails(details).
		JSON(c, http.StatusUnprocessableEntity)
}
