package lesson

import (
	"context"
	"database/sql"
	"errors"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/global"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/lesson"
	lpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/lesson"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (ls *LessonService) GetAllLessonsByModuleId(ctx context.Context, req *lpb.GetLessonsByModuleIdRequest) (*lpb.GetLessonsByModuleIdResponse, error) {
	ls.logger.WithFields(logrus.Fields{
		"method": "GetLessonsByModuleId",
		"params": req,
	}).Info("Get lessons by module id")

	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	moduleId := global.ConvertStringToGoogleUUID(req.ModuleId)

	lessons, err := ls.store.GetLessonsByModuleId(dbCtx, moduleId)

	if err != nil {
		ls.logger.WithError(err).Errorf("Failed to get lessons by module id")
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "Lessons not found")
		}
		return nil, status.Error(codes.Internal, "Failed to get lessons by module id")
	}

	return &lpb.GetLessonsByModuleIdResponse{
		Lessons: lesson.ConvertLessonsToProto(lessons),
	}, nil
}
