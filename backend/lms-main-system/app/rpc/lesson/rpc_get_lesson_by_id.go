package lesson

import (
	"context"
	"database/sql"
	"errors"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/global"
	lspb "github.com/multi-tenants-cms-golang/lms-sys/protogen/lesson"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (ls *LessonService) GetLessonById(ctx context.Context, req *lspb.GetLessonByIdRequest) (*lspb.GetLessonByIdResponse, error) {
	ls.logger.WithFields(logrus.Fields{
		"method": "GetLessonById",
		"params": req,
	}).Info("GetLessonById")

	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	pgUUID := global.ConvertStringToGoogleUUID(req.Id)

	l, err := ls.store.GetLessonById(dbCtx, pgUUID)
	if err != nil {
		ls.logger.WithError(err).Errorf("Failed to get lesson by id %s", req.Id)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "lesson not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get lesson by id %s", req.Id)
	}

	return &lspb.GetLessonByIdResponse{
		Lesson: &lspb.Lesson{
			Id:           l.LessonID.String(),
			Title:        l.Title,
			Content:      l.Content.String,
			MaterialType: string(l.MaterialType.MaterialType),
			Module: &lspb.Module{
				Id:   l.ModuleID.String(),
				Name: l.ModuleName,
			},
			CreatedAt: timestamppb.New(l.CreatedAt.Time),
			UpdatedAt: timestamppb.New(l.UpdatedAt.Time),
		},
	}, nil
}
