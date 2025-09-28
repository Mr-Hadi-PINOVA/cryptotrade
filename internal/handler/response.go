package handler

import (
    "errors"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"

    "cryptotrade/internal/repository"
    "cryptotrade/internal/service"
)

func respondError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, service.ErrValidation):
        c.JSON(http.StatusBadRequest, gin.H{"error": sanitizeValidationMessage(err)})
    case errors.Is(err, repository.ErrNotFound):
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
    case errors.Is(err, repository.ErrConflict):
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
    default:
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
        _ = c.Error(err)
    }
}

func sanitizeValidationMessage(err error) string {
    msg := err.Error()
    prefix := service.ErrValidation.Error() + ": "
    if strings.HasPrefix(msg, prefix) {
        return strings.TrimPrefix(msg, prefix)
    }
    return msg
}
