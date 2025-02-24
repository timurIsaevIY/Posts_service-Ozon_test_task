package comments

import (
	"Ozon_Post_comment_system/internal/models"
	"Ozon_Post_comment_system/internal/pkg/comments"
	"Ozon_Post_comment_system/internal/tools/errorChecker"
	"context"
)

type CommentResolvers struct {
	uc comments.CommentsUsecase
}

func NewCommentResolvers(uc comments.CommentsUsecase) *CommentResolvers {
	return &CommentResolvers{uc: uc}
}

// Comment resolvers
func (r *CommentResolvers) Replies(ctx context.Context, obj *models.Comment) ([]*models.Comment, error) {
	comments, err := r.uc.GetRepliesByCommentID(ctx, obj.ID)
	if err != nil {
		return nil, errorChecker.ErrorResponse(err)
	}
	return comments, nil
}

// Mutation resolvers
func (r *CommentResolvers) CreateComment(ctx context.Context, input models.InputComment) (*models.Comment, error) {
	res, err := r.uc.CreateComment(ctx, input)
	if err != nil {
		return nil, errorChecker.ErrorResponse(err)
	}
	return res, nil
}

// Comment resolvers
func (r *CommentResolvers) Comments(ctx context.Context, obj *models.Post, page, pageSize int) ([]*models.Comment, error) {
	res, err := r.uc.GetCommentsByPostID(ctx, obj.ID, page, pageSize)
	if err != nil {
		return nil, errorChecker.ErrorResponse(err)
	}
	return res, nil
}
