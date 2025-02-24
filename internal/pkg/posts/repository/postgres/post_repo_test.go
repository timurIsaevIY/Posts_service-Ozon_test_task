package postgres_test

import (
	"Ozon_Post_comment_system/internal/models"
	postsRepo "Ozon_Post_comment_system/internal/pkg/posts/repository/postgres"
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := postsRepo.NewPostsRepositoryImpl(db)
	ctx := context.Background()
	timeNow := time.Now()

	tests := []struct {
		desc        string
		mock        func()
		expectError bool
	}{
		{
			desc: "Successful CreatePost",
			mock: func() {
				mock.ExpectQuery("INSERT INTO posts").WithArgs("Test Post", "Test Author", "Test Content", true).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, timeNow))
			},
			expectError: false,
		},
		{
			desc: "CreatePost with DB Error",
			mock: func() {
				mock.ExpectQuery("INSERT INTO posts").WillReturnError(errors.New("db error"))
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			test.mock()
			_, err := repo.CreatePost(ctx, models.InputPost{Name: "Test Post", Author: "Test Author", Content: "Test Content", CommentsAllowed: true})
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetPostByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := postsRepo.NewPostsRepositoryImpl(db)
	ctx := context.Background()
	timeNow := time.Now()

	tests := []struct {
		desc        string
		mock        func()
		expectError bool
	}{
		{
			desc: "Successful GetPostByID",
			mock: func() {
				mock.ExpectQuery(`SELECT id, created_at, name, author, content, comments_allowed FROM posts WHERE id = \$1`).WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "name", "author", "content", "comments_allowed"}).
						AddRow(1, timeNow, "Test Post", "Test Author", "Test Content", true))
			},
			expectError: false,
		},
		{
			desc: "GetPostByID Not Found",
			mock: func() {
				mock.ExpectQuery(`SELECT id, created_at, name, author, content, comments_allowed FROM posts WHERE id = \$1`).WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			test.mock()
			_, err := repo.GetPostByID(ctx, 1)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetAllPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := postsRepo.NewPostsRepositoryImpl(db)
	ctx := context.Background()

	tests := []struct {
		desc        string
		mock        func()
		expectError bool
	}{
		{
			desc: "Successful GetAllPosts",
			mock: func() {
				mock.ExpectQuery(`SELECT id, created_at, name, author, content, comments_allowed FROM posts ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).
					WithArgs(10, 0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "name", "author", "content", "comments_allowed"}).
						AddRow(1, time.Now(), "Post1", "Author1", "Content1", true).
						AddRow(2, time.Now(), "Post2", "Author2", "Content2", false))
			},
			expectError: false,
		},
		{
			desc: "GetAllPosts DB Error",
			mock: func() {
				mock.ExpectQuery(`SELECT id, created_at, name, author, content, comments_allowed FROM posts ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).WillReturnError(errors.New("db error"))
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			test.mock()
			_, err := repo.GetAllPosts(ctx, 10, 0)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
