package posts

import (
	"Ozon_Post_comment_system/internal/models"
	"context"
)

type PostsUsecase interface {
	CreatePost(ctx context.Context, post models.InputPost) (*models.Post, error)
	GetPostById(ctx context.Context, id int) (*models.Post, error)
	GetAllPosts(ctx context.Context, page, pageSize int) ([]*models.Post, error)
	SubscribeToComments(ctx context.Context, postId int) (<-chan *models.Comment, error)
}

type PostsRepository interface {
	CreatePost(ctx context.Context, post models.InputPost) (*models.Post, error)
	GetPostByID(ctx context.Context, id int) (*models.Post, error)
	GetAllPosts(ctx context.Context, limit, offset int) ([]*models.Post, error)
}
