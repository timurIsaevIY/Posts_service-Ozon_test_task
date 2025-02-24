package middleware

import (
	"Ozon_Post_comment_system/internal/logger"
	"context"
	"github.com/99designs/gqlgen/graphql"
)

// Middleware для логирования входящих GraphQL-запросов
func LoggingMiddleware(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	operationCtx := graphql.GetOperationContext(ctx)
	ctx = logger.LogGraphQLStart(ctx, operationCtx.OperationName, operationCtx.RawQuery, operationCtx.Variables)
	return next(ctx)
}
