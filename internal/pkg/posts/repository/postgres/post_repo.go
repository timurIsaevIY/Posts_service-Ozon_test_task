package postgres

import (
	"Ozon_Post_comment_system/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type PostsRepositoryImpl struct {
	db *sql.DB
}

func NewPostsRepositoryImpl(db *sql.DB) *PostsRepositoryImpl {
	return &PostsRepositoryImpl{db: db}
}

func (r *PostsRepositoryImpl) CreatePost(ctx context.Context, post models.InputPost) (*models.Post, error) {
	query := `
		INSERT INTO posts (created_at, name, author, content, comments_allowed)
		VALUES (NOW(), $1, $2, $3, $4)
		RETURNING id, created_at
	`
	var id int
	var createdAt time.Time
	err := r.db.QueryRowContext(ctx, query, post.Name, post.Author, post.Content, post.CommentsAllowed).Scan(&id, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("error creating post: %w", err)
	}

	return &models.Post{
		ID:              id,
		CreatedAt:       createdAt,
		Name:            post.Name,
		Author:          post.Author,
		Content:         post.Content,
		CommentsAllowed: post.CommentsAllowed,
	}, nil
}

func (r *PostsRepositoryImpl) GetPostByID(ctx context.Context, id int) (*models.Post, error) {
	query := `
		SELECT id, created_at, name, author, content, comments_allowed
		FROM posts
		WHERE id = $1
	`
	var post models.Post
	err := r.db.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.CreatedAt, &post.Name, &post.Author, &post.Content, &post.CommentsAllowed)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("post not found :%w", models.ErrNotFound)
		}
		return nil, fmt.Errorf("error getting post: %w", err)
	}
	return &post, nil
}

func (r *PostsRepositoryImpl) GetAllPosts(ctx context.Context, limit, offset int) ([]*models.Post, error) {
	query := `
		SELECT id, created_at, name, author, content, comments_allowed
		FROM posts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting posts: %w", err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.CreatedAt, &post.Name, &post.Author, &post.Content, &post.CommentsAllowed)
		if err != nil {
			return nil, fmt.Errorf("error getting posts: %w", err)
		}
		posts = append(posts, &post)
	}
	return posts, nil
}
