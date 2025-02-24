package inMemory

import (
	"Ozon_Post_comment_system/internal/models"
	"context"
	"sync"
	"time"
)

type InMemoryCommentsRepository struct {
	mu       sync.RWMutex
	comments map[int]*models.Comment
	nextID   int
}

func NewInMemoryCommentsRepository() *InMemoryCommentsRepository {
	return &InMemoryCommentsRepository{
		comments: make(map[int]*models.Comment),
		nextID:   1,
	}
}

func (r *InMemoryCommentsRepository) CreateComment(ctx context.Context, comment models.InputComment) (*models.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID
	r.nextID++

	newComment := &models.Comment{
		ID:        id,
		Content:   comment.Content,
		Author:    comment.Author,
		ReplyTo:   comment.ReplyTo,
		Post:      comment.Post,
		CreatedAt: time.Now(),
	}

	if replyTo, exists := r.comments[comment.Post]; exists {
		if newComment.ReplyTo != nil {
			replyTo.Replies = append(replyTo.Replies, newComment)
		}
	}

	r.comments[id] = newComment
	return newComment, nil
}

func (r *InMemoryCommentsRepository) GetRepliesByCommentID(ctx context.Context, commentId int) ([]*models.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	comment, exists := r.comments[commentId]
	if !exists {
		return []*models.Comment{}, nil
	}

	return comment.Replies, nil
}

func (r *InMemoryCommentsRepository) GetCommentsByPostID(ctx context.Context, postID, limit, offset int) ([]*models.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var comments []*models.Comment
	for _, comment := range r.comments {
		if comment.Post == postID && comment.ReplyTo == nil {
			comments = append(comments, comment)
		}
	}

	start := offset
	end := offset + limit
	if start > len(comments) {
		return []*models.Comment{}, nil
	}
	if end > len(comments) {
		end = len(comments)
	}

	return comments[start:end], nil
}
