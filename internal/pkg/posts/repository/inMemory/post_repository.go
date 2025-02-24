package inMemory

import (
	"Ozon_Post_comment_system/internal/models"
	"context"
	"errors"
	"sync"
	"time"
)

type InMemoryPostsRepository struct {
	mu     sync.RWMutex
	posts  map[int]*models.Post
	nextID int
}

func NewInMemoryPostsRepository() *InMemoryPostsRepository {
	return &InMemoryPostsRepository{
		posts:  make(map[int]*models.Post),
		nextID: 1,
	}
}

func (r *InMemoryPostsRepository) CreatePost(ctx context.Context, post models.InputPost) (*models.Post, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID
	r.nextID++

	newPost := &models.Post{
		ID:              id,
		CreatedAt:       time.Now(),
		Name:            post.Name,
		Author:          post.Author,
		Content:         post.Content,
		CommentsAllowed: post.CommentsAllowed,
	}

	r.posts[id] = newPost
	return newPost, nil
}

func (r *InMemoryPostsRepository) GetPostByID(ctx context.Context, id int) (*models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, exists := r.posts[id]
	if !exists {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (r *InMemoryPostsRepository) GetAllPosts(ctx context.Context, limit, offset int) ([]*models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var posts []*models.Post
	count := 0

	for i := r.nextID - 1 - offset; i > 0 && count <= limit; i-- {
		post, exist := r.posts[i]
		if exist {
			posts = append(posts, post)
			count++
		}
	}
	return posts, nil
}
