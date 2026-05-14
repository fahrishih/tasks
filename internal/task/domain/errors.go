package domain

import "errors"

var (
	ErrTaskNotFound  = errors.New("task not found")
	ErrTaskDuplicate = errors.New("task already exists")
	ErrInvalidTitle  = errors.New("title is required")
	ErrTitleTooLong  = errors.New("title must be 200 characters or fewer")
)
