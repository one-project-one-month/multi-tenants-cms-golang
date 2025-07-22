package lesson

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/global"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils"
	lspb "github.com/multi-tenants-cms-golang/lms-sys/protogen/lesson"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (ls *LessonService) UpdateLesson(ctx context.Context, req *lspb.UpdateLessonRequest) (*lspb.UpdateLessonResponse, error) {
	ls.logger.WithFields(logrus.Fields{
		"method": "UpdateLessonByID",
		"params": req,
	}).Info("Updating Lesson by ID")

		md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, utils.ErrMissingOrganization().ToGRPCStatus()
	}

	// INFO: To be deleted
	fmt.Printf("rpc_get_module: org FromIncomingContext: %+v\n", md)

	orgValues := md.Get("x-organisation")
	if len(orgValues) <= 0 {
		return nil, utils.ErrMissingOrganization().ToGRPCStatus()
	}

	orgName := orgValues[0]

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id required")
	}

	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	lessonID := global.ConvertStringToGoogleUUID(req.Id)

	_, err := ls.store.GetLessonById(dbCtx, lessonID)
	if err != nil {
		ls.logger.WithError(err).Errorf("Failed to fetch lesson with ID %s", req.Id)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "lesson not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch lesson with id %s", req.Id)
	}

	args := repo.GetModuleByIDWithTenantParams{
		ModuleID: global.ConvertStringToGoogleUUID(req.ModuleId),
		Namespace: orgName,
	}

	m, err := ls.store.GetModuleByIDWithTenant(dbCtx, args)
	if err != nil {
		ls.logger.WithError(err).Errorf("Failed to fetch module with ID %s", req.ModuleId)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "module not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch module with id %s", req.ModuleId)
	}

	updateArgs := db.UpdateLessonParams{
		LessonID: lessonID,
		Title:    req.Title,
		Content:  *global.ConvertStringToPgText(req.Content),
		MaterialType: db.NullMaterialType{
			MaterialType: db.MaterialType(req.MaterialType),
			Valid:        req.MaterialType != "",
		},
		ModuleID: global.ConvertStringToGoogleUUID(req.ModuleId),
	}

	l, err := ls.store.UpdateLesson(dbCtx, updateArgs)
	if err != nil {
		ls.logger.WithError(err).Errorf("Failed to update lesson with ID %s", req.Id)
		return nil, status.Errorf(codes.Internal, "failed to update lesson with id %s", req.Id)
	}

	return &lspb.UpdateLessonResponse{
		Lesson: &lspb.Lesson{
			Id:           l.LessonID.String(),
			Title:        l.Title,
			Content:      l.Content.String,
			MaterialType: string(l.MaterialType.MaterialType),
			Module: &lspb.Module{
				Id:   l.ModuleID.String(),
				Name: m.ModuleName,
			},
			CreatedAt: timestamppb.New(l.CreatedAt.Time),
			UpdatedAt: timestamppb.New(l.UpdatedAt.Time),
		},
	}, nil
}
