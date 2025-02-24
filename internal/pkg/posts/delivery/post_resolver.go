package posts

import (
	"Ozon_Post_comment_system/internal/models"
	"Ozon_Post_comment_system/internal/pkg/posts"
	"Ozon_Post_comment_system/internal/tools/errorChecker"

	"context"
)

type PostResolvers struct {
	uc posts.PostsUsecase
}

func NewPostResolvers(uc posts.PostsUsecase) *PostResolvers {
	return &PostResolvers{uc: uc}
}

// Query resolvers
func (r *PostResolvers) GetAllPosts(ctx context.Context, page, pageSize int) ([]*models.Post, error) {
	post, err := r.uc.GetAllPosts(ctx, page, pageSize)
	if err != nil {
		return nil, errorChecker.ErrorResponse(err)
	}
	return post, nil
}

func (r *PostResolvers) GetPostByID(ctx context.Context, id int) (*models.Post, error) {
	post, err := r.uc.GetPostById(ctx, id)
	if err != nil {
		return nil, errorChecker.ErrorResponse(err)
	}
	return post, nil
}

// Mutation resolvers
func (r *PostResolvers) CreatePost(ctx context.Context, post models.InputPost) (*models.Post, error) {
	postAns, err := r.uc.CreatePost(ctx, post)
	if err != nil {
		return nil, errorChecker.ErrorResponse(err)
	}
	return postAns, nil
}

// Subscription resolvers
func (r *PostResolvers) CommentsSubscription(ctx context.Context, postID int) (<-chan *models.Comment, error) {
	res, err := r.uc.SubscribeToComments(ctx, postID)
	if err != nil {
		return nil, errorChecker.ErrorResponse(err)
	}
	return res, nil
}
