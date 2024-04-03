package utils

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		resp, err = handler(ctx, req)
		duration := time.Since(start)

		p, _ := peer.FromContext(ctx)

		// Log the request
		log.Info().
			Str("method", info.FullMethod).
			Str("peer", p.Addr.String()).
			Str("duration", duration.String()).
			Interface("request", req).
			Msg("gRPC request")

		// Log the response
		if err != nil {
			log.Error().
				Str("method", info.FullMethod).
				Str("duration", duration.String()).
				Err(err).
				Msg("gRPC response error")
		} else {
			log.Info().
				Str("method", info.FullMethod).
				Str("duration", duration.String()).
				Interface("response", resp).
				Msg("gRPC response")
		}

		return resp, err
	}
}
