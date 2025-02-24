package postsUsecase

import (
	"Ozon_Post_comment_system/internal/models"
	"Ozon_Post_comment_system/internal/notifications"
	"Ozon_Post_comment_system/internal/pkg/posts"
	"Ozon_Post_comment_system/internal/tools/pagination"
	"Ozon_Post_comment_system/internal/tools/validation"
	"context"
)

type PostsUsecaseImpl struct {
	repo     posts.PostsRepository
	observer notifications.Observers
}

func NewPostsUsecaseImpl(repo posts.PostsRepository, observer notifications.Observers) *PostsUsecaseImpl {
	return &PostsUsecaseImpl{repo, observer}
}

func (u *PostsUsecaseImpl) CreatePost(ctx context.Context, post models.InputPost) (*models.Post, error) {
	return u.repo.CreatePost(ctx, post)
}

func (u *PostsUsecaseImpl) GetPostById(ctx context.Context, id int) (*models.Post, error) {
	if err := validation.ValidateID(id); err != nil {
		return nil, err
	}
	return u.repo.GetPostByID(ctx, id)
}

func (u *PostsUsecaseImpl) GetAllPosts(ctx context.Context, page, pageSize int) ([]*models.Post, error) {
	if err := validation.ValidatePagination(page, pageSize); err != nil {
		return nil, err
	}
	limit, offset := pagination.Pagination(page, pageSize)
	return u.repo.GetAllPosts(ctx, limit, offset)
}

func (u *PostsUsecaseImpl) SubscribeToComments(ctx context.Context, postID int) (<-chan *models.Comment, error) {
	if err := validation.ValidateID(postID); err != nil {
		return nil, err
	}
	return u.observer.Subscribe(ctx, postID)
}
