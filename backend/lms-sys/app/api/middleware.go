package api

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func (s *Server) unaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		s.logger.WithFields(logrus.Fields{
			"method": info.FullMethod,
		}).Info("incoming gRPC request")
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"method": info.FullMethod,
				"error":  err,
			}).Error("gRPC request failed")
		}
		return resp, err
	}
}
