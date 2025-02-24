package comments

import (
	"Ozon_Post_comment_system/internal/models"
	"context"
)

type CommentsUsecase interface {
	CreateComment(ctx context.Context, comment models.InputComment) (*models.Comment, error)
	GetRepliesByCommentID(ctx context.Context, commentId int) ([]*models.Comment, error)
	GetCommentsByPostID(ctx context.Context, postId, page, pageSize int) ([]*models.Comment, error)
}

type CommentsRepository interface {
	CreateComment(ctx context.Context, comment models.InputComment) (*models.Comment, error)
	GetRepliesByCommentID(ctx context.Context, commentId int) ([]*models.Comment, error)
	GetCommentsByPostID(ctx context.Context, postID, limit, offset int) ([]*models.Comment, error)
}
