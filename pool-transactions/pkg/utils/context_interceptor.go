package utils

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey int

const UserIdKey contextKey = 0

func UserIdInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "No metadata provided")
	}

	userIDs := md.Get("user_id")
	if len(userIDs) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "User ID not provided")
	}

	newCtx := context.WithValue(ctx, UserIdKey, userIDs[0])
	return handler(newCtx, req)
}
