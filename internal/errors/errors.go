package errors

import (
	"errors"
	"fmt"
)

// Common error types
var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrAlreadyExists      = errors.New("already exists")
	ErrOperationFailed    = errors.New("operation failed")
	ErrMissingRequirement = errors.New("missing requirement")
)

// UserError represents an error that should be displayed to the user
type UserError struct {
	Err     error
	Message string
	Action  string
}

func (e *UserError) Error() string {
	if e.Action != "" {
		return fmt.Sprintf("%s\n提案: %s", e.Message, e.Action)
	}
	return e.Message
}

func (e *UserError) Unwrap() error {
	return e.Err
}

// Helper functions to create user-friendly errors

func NewNotFoundError(resource string, suggestion string) error {
	return &UserError{
		Err:     ErrNotFound,
		Message: fmt.Sprintf("%sが見つかりませんでした", resource),
		Action:  suggestion,
	}
}

func NewInvalidArgumentError(arg string, reason string, suggestion string) error {
	return &UserError{
		Err:     ErrInvalidArgument,
		Message: fmt.Sprintf("引数%sが不正です: %s", arg, reason),
		Action:  suggestion,
	}
}

func NewPermissionDeniedError(resource string, suggestion string) error {
	return &UserError{
		Err:     ErrPermissionDenied,
		Message: fmt.Sprintf("%sへのアクセス権限がありません", resource),
		Action:  suggestion,
	}
}

func NewAlreadyExistsError(resource string, suggestion string) error {
	return &UserError{
		Err:     ErrAlreadyExists,
		Message: fmt.Sprintf("%sは既に存在します", resource),
		Action:  suggestion,
	}
}

func NewOperationFailedError(operation string, reason error, suggestion string) error {
	return &UserError{
		Err:     ErrOperationFailed,
		Message: fmt.Sprintf("%sに失敗しました: %v", operation, reason),
		Action:  suggestion,
	}
}

func NewMissingRequirementError(requirement string, suggestion string) error {
	return &UserError{
		Err:     ErrMissingRequirement,
		Message: fmt.Sprintf("%sが必要です", requirement),
		Action:  suggestion,
	}
}

// WrapError wraps an error with a user-friendly message and suggested action
func WrapError(err error, message string, suggestion string) error {
	return &UserError{
		Err:     err,
		Message: message,
		Action:  suggestion,
	}
}
