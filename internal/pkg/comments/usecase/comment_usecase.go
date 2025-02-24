package commentsUsecase

import (
	"Ozon_Post_comment_system/internal/models"
	"Ozon_Post_comment_system/internal/notifications"
	"Ozon_Post_comment_system/internal/pkg/comments"
	"Ozon_Post_comment_system/internal/tools/pagination"
	"Ozon_Post_comment_system/internal/tools/validation"
	"context"
)

type CommentsUsecaseImpl struct {
	repo     comments.CommentsRepository
	observer notifications.Observers
}

// NewCommentsUsecaseImpl создаёт новый экземпляр CommentsUsecaseImpl
func NewCommentsUsecaseImpl(repo comments.CommentsRepository, observer notifications.Observers) *CommentsUsecaseImpl {
	return &CommentsUsecaseImpl{repo, observer}
}

// CreateComment создаёт новый комментарий
func (uc *CommentsUsecaseImpl) CreateComment(ctx context.Context, input models.InputComment) (*models.Comment, error) {
	if err := validation.ValidateText(input.Content, 2000); err != nil {
		return nil, err
	}
	comment, err := uc.repo.CreateComment(ctx, input)
	if err != nil {
		return nil, err
	}
	uc.observer.Notify(comment.Post, comment)
	return comment, nil
}

// GetRepliesByCommentID возвращает ответы на комментарий по его ID
func (uc *CommentsUsecaseImpl) GetRepliesByCommentID(ctx context.Context, commentId int) ([]*models.Comment, error) {
	if err := validation.ValidateID(commentId); err != nil {
		return nil, err
	}
	return uc.repo.GetRepliesByCommentID(ctx, commentId)
}

func (u *CommentsUsecaseImpl) GetCommentsByPostID(ctx context.Context, postID, page, pageSize int) ([]*models.Comment, error) {
	if err := validation.ValidateID(postID); err != nil {
		return nil, err
	}
	if err := validation.ValidatePagination(page, pageSize); err != nil {
		return nil, err
	}
	limit, offset := pagination.Pagination(page, pageSize)
	return u.repo.GetCommentsByPostID(ctx, postID, limit, offset)
}
