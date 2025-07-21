package lesson

import (
	"context"
	"github.com/google/uuid"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/global"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/lesson"
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	lspb "github.com/multi-tenants-cms-golang/lms-sys/protogen/lesson"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (ls *LessonService) CreateLesson(ctx context.Context, req *lspb.CreateLessonRequest) (*lspb.CreateLessonResponse, error) {
	ls.logger.WithFields(logrus.Fields{
		"method": "CreateLesson",
	}).Info("Creating Lesson")

	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}
	if req.Content == "" {
		return nil, status.Error(codes.InvalidArgument, "content is required")
	}
	pgModuleId, err := global.ConvertStringToUUID(req.ModuleId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	pgModuleUUID := uuid.MustParse(pgModuleId.String())

	module, err := ls.store.GetModuleByID(dbCtx, pgModuleUUID)

	if err != nil {
		ls.logger.WithError(err).Error("Failed to fetch module")
		return nil, status.Error(codes.Internal, "Failed to fetch module info")
	}

	params := db.CreateLessonParams{
		Title:        req.Title,
		Content:      *global.ConvertStringToPgText(req.Content),
		MaterialType: db.NullMaterialType{MaterialType: db.MaterialType(req.MaterialType), Valid: true},
		ModuleID:     module.ModuleID,
	}
	l, err := ls.store.CreateLesson(dbCtx, params)

	if err != nil {
		ls.logger.WithError(err).Error("Failed to create Lesson")
		return nil, status.Error(codes.Internal, "Failed to create Lesson")
	}

	return &lspb.CreateLessonResponse{
		Lesson: lesson.ConvertLessonToProtoFromLesson(l, module.ModuleName),
	}, nil
}
