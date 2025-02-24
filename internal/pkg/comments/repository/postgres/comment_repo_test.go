package postgres

import (
	"Ozon_Post_comment_system/internal/models"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewCommentsRepository(db)
	ctx := context.Background()
	timeNow := time.Now()

	tests := []struct {
		desc        string
		mock        func()
		expectError bool
	}{
		{
			desc: "Successful CreateComment",
			mock: func() {
				mock.ExpectQuery("INSERT INTO comments").WithArgs("Test Content", "Test Author", nil, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "content", "author", "reply_to", "post_id", "created_at"}).
						AddRow(1, "Test Content", "Test Author", nil, 1, timeNow))
			},
			expectError: false,
		},
		{
			desc: "CreateComment with DB Error",
			mock: func() {
				mock.ExpectQuery("INSERT INTO comments").WillReturnError(errors.New("db error"))
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			test.mock()
			_, err := repo.CreateComment(ctx, models.InputComment{Author: "Test Author", Content: "Test Content", Post: 1, ReplyTo: nil})
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetCommentsByPostID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewCommentsRepository(db)
	ctx := context.Background()
	timeNow := time.Now()

	tests := []struct {
		desc        string
		mock        func()
		expectError bool
	}{
		{
			desc: "GetCommentsByPostID Success",
			mock: func() {
				mock.ExpectQuery(`SELECT id, created_at, author, content, post_id, reply_to FROM comments WHERE post_id = \$1 AND reply_to IS NULL ORDER BY created_at DESC LIMIT \$2 OFFSET \$3`).
					WithArgs(1, 10, 0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "author", "content", "post_id", "reply_to"}).
						AddRow(1, timeNow, "Commenter1", "Comment1", 1, nil))
			},
			expectError: false,
		},
		{
			desc: "GetCommentsByPostID Not Found",
			mock: func() {
				mock.ExpectQuery(`SELECT id, created_at, author, content, post_id, reply_to FROM comments WHERE post_id = \$1 AND reply_to IS NULL ORDER BY created_at DESC LIMIT \$2 OFFSET \$3`).
					WillReturnError(errors.New("db error"))
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			test.mock()
			_, err := repo.GetCommentsByPostID(ctx, 1, 10, 0)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
