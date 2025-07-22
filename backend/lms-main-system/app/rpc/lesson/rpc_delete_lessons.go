package lesson

import (
	"context"
	"github.com/google/uuid"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/global"
	lspb "github.com/multi-tenants-cms-golang/lms-sys/protogen/lesson"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (ls *LessonService) DeleteLessons(ctx context.Context, req *lspb.DeleteLessonsRequest) (*lspb.DeleteLessonsResponse, error) {
	ls.logger.WithFields(logrus.Fields{
		"method": "DeleteLessons",
		"params": req,
	}).Info("Deleting lessons")

	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if len(req.Ids) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no lesson IDs provided for deletion : %v", req.Ids)
	}

	var lessonUUIDs []uuid.UUID
	for _, id := range req.Ids {
		uid := global.ConvertStringToGoogleUUID(id)
		lessonUUIDs = append(lessonUUIDs, uid)
	}

	_, err := ls.store.DeleteLessons(dbCtx, lessonUUIDs)

	if err != nil {
		ls.logger.WithError(err).Errorf("Failed to delete lessons: %v", req.Ids)
		return nil, status.Errorf(codes.Internal, "failed to delete lessons")
	}

	return &lspb.DeleteLessonsResponse{
		Success: true,
	}, nil
}
