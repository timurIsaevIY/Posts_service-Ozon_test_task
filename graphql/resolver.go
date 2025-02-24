package graph

import (
	"Ozon_Post_comment_system/internal/models"
	commentDelivery "Ozon_Post_comment_system/internal/pkg/comments/delivery"
	postDelivery "Ozon_Post_comment_system/internal/pkg/posts/delivery"
	"context"
)

// Resolver is the root resolver structure
type Resolver struct {
	postResolvers    *postDelivery.PostResolvers
	commentResolvers *commentDelivery.CommentResolvers
}

// NewResolver creates a new resolver with necessary services
func NewResolver(
	postService *postDelivery.PostResolvers,
	commentService *commentDelivery.CommentResolvers,
) *Resolver {
	return &Resolver{
		postResolvers:    postService,
		commentResolvers: commentService,
	}
}

// CommentResolver implementation
type commentResolver struct{ *Resolver }

func (r *commentResolver) Replies(ctx context.Context, obj *models.Comment) ([]*models.Comment, error) {
	return r.commentResolvers.Replies(ctx, obj)
}

// MutationResolver implementation
type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreatePost(ctx context.Context, post models.InputPost) (*models.Post, error) {
	return r.postResolvers.CreatePost(ctx, post)
}

func (r *mutationResolver) CreateComment(ctx context.Context, input models.InputComment) (*models.Comment, error) {
	return r.commentResolvers.CreateComment(ctx, input)
}

// PostResolver implementation
type postResolver struct{ *Resolver }

func (r *postResolver) Comments(ctx context.Context, obj *models.Post, page, pageSize int) ([]*models.Comment, error) {
	return r.commentResolvers.Comments(ctx, obj, page, pageSize)
}

// QueryResolver implementation
type queryResolver struct{ *Resolver }

func (r *queryResolver) GetAllPosts(ctx context.Context, page, pageSize int) ([]*models.Post, error) {
	return r.postResolvers.GetAllPosts(ctx, page, pageSize)
}

func (r *queryResolver) GetPostByID(ctx context.Context, id int) (*models.Post, error) {
	return r.postResolvers.GetPostByID(ctx, id)
}

// SubscriptionResolver implementation
type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) CommentsSubscription(ctx context.Context, postID int) (<-chan *models.Comment, error) {
	return r.postResolvers.CommentsSubscription(ctx, postID)
}

// Helper methods to return typed resolvers
func (r *Resolver) Comment() CommentResolver           { return &commentResolver{r} }
func (r *Resolver) Mutation() MutationResolver         { return &mutationResolver{r} }
func (r *Resolver) Post() PostResolver                 { return &postResolver{r} }
func (r *Resolver) Query() QueryResolver               { return &queryResolver{r} }
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }
