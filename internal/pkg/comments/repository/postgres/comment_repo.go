package postgres

import (
	"Ozon_Post_comment_system/internal/models"
	"context"
	"database/sql"
	"fmt"
)

type СommentsRepositoryImpl struct {
	db *sql.DB
}

// NewCommentsRepository создаёт новый экземпляр CommentsRepository
func NewCommentsRepository(db *sql.DB) *СommentsRepositoryImpl {
	return &СommentsRepositoryImpl{db}
}

// CreateComment сохраняет комментарий в базе данных
func (r *СommentsRepositoryImpl) CreateComment(ctx context.Context, comment models.InputComment) (*models.Comment, error) {
	query := `
		INSERT INTO comments (content, author, reply_to, post_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, content, author, reply_to, post_id, created_at
	`
	row := r.db.QueryRowContext(ctx, query, comment.Content, comment.Author, comment.ReplyTo, comment.Post)

	var createdComment models.Comment
	err := row.Scan(
		&createdComment.ID,
		&createdComment.Content,
		&createdComment.Author,
		&createdComment.ReplyTo,
		&createdComment.Post,
		&createdComment.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании комментария: %w", err)
	}
	return &createdComment, nil
}

// GetRepliesByCommentID возвращает ответы на комментарий по его ID
func (r *СommentsRepositoryImpl) GetRepliesByCommentID(ctx context.Context, commentId int) ([]*models.Comment, error) {
	query := `
		SELECT id, content, author, reply_to, post_id, created_at
		FROM comments
		WHERE reply_to = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, commentId)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении ответов: %w", err)
	}
	defer rows.Close()

	var replies []*models.Comment
	for rows.Next() {
		var reply models.Comment
		err := rows.Scan(
			&reply.ID,
			&reply.Content,
			&reply.Author,
			&reply.ReplyTo,
			&reply.Post,
			&reply.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании ответа: %w", err)
		}
		replies = append(replies, &reply)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по ответам: %w", err)
	}

	return replies, nil
}

func (r *СommentsRepositoryImpl) GetCommentsByPostID(ctx context.Context, postID, limit, offset int) ([]*models.Comment, error) {
	query := `
		SELECT id, created_at, author, content, post_id, reply_to
		FROM comments
		WHERE post_id = $1 AND reply_to IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting comments: %w", err)
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.CreatedAt, &comment.Author, &comment.Content, &comment.Post, &comment.ReplyTo)
		if err != nil {
			return nil, fmt.Errorf("error getting comments: %w", err)
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}
