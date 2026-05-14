package ginhttp

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fahrishih/tasks/internal/task/domain"
)

// statusForError maps domain sentinel errors to HTTP status codes.
// This is THE place HTTP semantics meet domain semantics. Keeping it
// in one function means we can't accidentally return inconsistent
// statuses for the same error across different handlers.
func statusForError(err error) int {
	switch {
	case errors.Is(err, domain.ErrTaskNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrTaskDuplicate):
		return http.StatusConflict
	case errors.Is(err, domain.ErrInvalidTitle),
		errors.Is(err, domain.ErrTitleTooLong):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func writeError(c *gin.Context, err error) {
	c.JSON(statusForError(err), gin.H{"error": err.Error()})
}
