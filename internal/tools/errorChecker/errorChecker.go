package errorChecker

import (
	"Ozon_Post_comment_system/internal/models"
	"errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"net/http"
	"strconv"
)

func ErrorResponse(err error) *gqlerror.Error {
	if errors.Is(err, models.ErrNotFound) {
		return graphQLError(models.ErrNotFound, strconv.Itoa(http.StatusNotFound))
	} else if errors.Is(err, models.ErrValidation) {
		return graphQLError(err, strconv.Itoa(http.StatusBadRequest))
	}
	return graphQLError(models.ErrInternal, strconv.Itoa(http.StatusInternalServerError))
}

func graphQLError(err error, code string) *gqlerror.Error {
	if err == nil {
		return nil
	}
	return &gqlerror.Error{
		Message: err.Error(),
		Extensions: map[string]interface{}{
			"code": code,
		},
	}
}
