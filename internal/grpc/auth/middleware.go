package auth

import (
	"context"
	"strings"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func EnsureValidToken(secret string, requireAuth map[string]bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if !requireAuth[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.Log().Debug(ctx, "metadata is not provided")
			return nil, status.Error(codes.Unauthenticated, core.ErrUnauthorized.Error())
		}

		authorization := md.Get("authorization")
		if len(authorization) == 0 {
			logger.Log().Debug(ctx, "token is not provided")
			return nil, status.Error(codes.Unauthenticated, core.ErrUnauthorized.Error())
		}

		tokenString := strings.TrimPrefix(authorization[0], "Bearer")
		tokenString = strings.TrimSpace(tokenString)

		id, err := validToken(ctx, tokenString, secret)
		if err != nil {
			logger.Log().Debug(ctx, err.Error())
			return nil, status.Error(codes.Unauthenticated, core.ErrUnauthorized.Error())
		}

		ctx = context.WithValue(ctx, userIDContextKey, *id)

		return handler(ctx, req)
	}
}
