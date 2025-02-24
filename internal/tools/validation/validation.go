package validation

import (
	"Ozon_Post_comment_system/internal/models"
	"fmt"
)

func ValidateID(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid ID: must be greater than 0: %w", models.ErrValidation)
	}
	return nil
}

func ValidatePagination(page, pageSize int) error {
	if pageSize <= 0 || pageSize > 100 {
		return fmt.Errorf("pageSize must be between 1 and 100: %w", models.ErrValidation)
	}
	if page <= 0 {
		return fmt.Errorf("page number must be 1 or greater: %w", models.ErrValidation)
	}
	return nil
}
func ValidateText(text string, length int) error {
	if len(text) > length {
		return fmt.Errorf("text must be no more than %d characters: %w", length, models.ErrValidation)
	}
	return nil
}
