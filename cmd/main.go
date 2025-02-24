package main

import (
	"Ozon_Post_comment_system/graphql"
	"Ozon_Post_comment_system/internal/config"
	"Ozon_Post_comment_system/internal/logger"
	"Ozon_Post_comment_system/internal/middleware"
	"Ozon_Post_comment_system/internal/notifications"
	comments2 "Ozon_Post_comment_system/internal/pkg/comments"
	"Ozon_Post_comment_system/internal/pkg/comments/delivery"
	inMemoryComment "Ozon_Post_comment_system/internal/pkg/comments/repository/inMemory"
	commentsRepo "Ozon_Post_comment_system/internal/pkg/comments/repository/postgres"
	commentsUsecase "Ozon_Post_comment_system/internal/pkg/comments/usecase"
	posts2 "Ozon_Post_comment_system/internal/pkg/posts"
	posts "Ozon_Post_comment_system/internal/pkg/posts/delivery"
	inMemoryPost "Ozon_Post_comment_system/internal/pkg/posts/repository/inMemory"
	postsRepo "Ozon_Post_comment_system/internal/pkg/posts/repository/postgres"
	postsUsecase "Ozon_Post_comment_system/internal/pkg/posts/usecase"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
	"github.com/vektah/gqlparser/v2/ast"
	"log"
	"net/http"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger := logger.NewPrettyHandler(os.Stdout, logger.PrettyHandlerOptions{})
	slog.SetDefault(slog.New(logger))

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Database.DbHost, cfg.Database.DbPort, cfg.Database.DbUser, cfg.Database.DbPass, cfg.Database.DbName))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(10 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)
	log.Printf("DB_HOST: %s, DB_PORT: %d, DB_USER: %s, DB_PASS: %s, DB_NAME: %s", cfg.Database.DbHost, cfg.Database.DbPort, cfg.Database.DbUser, cfg.Database.DbPass, cfg.Database.DbName)
	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	observer := notifications.NewObserver()
	var postRepo posts2.PostsRepository
	var commentRepo comments2.CommentsRepository

	if cfg.StorageType == "postgres" {
		postRepo = postsRepo.NewPostsRepositoryImpl(db)
		commentRepo = commentsRepo.NewCommentsRepository(db)
	} else if cfg.StorageType == "in-memory" {
		postRepo = inMemoryPost.NewInMemoryPostsRepository()
		commentRepo = inMemoryComment.NewInMemoryCommentsRepository()
	} else {
		log.Fatalf("storage_type must be postgres/in-memory")
	}

	postUsecase := postsUsecase.NewPostsUsecaseImpl(postRepo, observer)
	postResolvers := posts.NewPostResolvers(postUsecase)

	commentUsecase := commentsUsecase.NewCommentsUsecaseImpl(commentRepo, observer)
	commentResolvers := comments.NewCommentResolvers(commentUsecase)

	resolver := graph.NewResolver(postResolvers, commentResolvers)

	// Настройка GraphQL сервера
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	srv.AroundResponses(middleware.LoggingMiddleware)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	graphqlServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HttpServer.Address),
		Handler: nil,
	}

	go func() {
		slog.Info(fmt.Sprintf("GraphQL server listening on :%s", cfg.HttpServer.Address))
		if err := graphqlServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to serve GraphQL", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	slog.Info("Shutting down HTTP server...")
	if err := graphqlServer.Shutdown(context.Background()); err != nil {
		slog.Error("HTTP server shutdown failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("HTTP server gracefully stopped")
}
