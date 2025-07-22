package lesson

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/global"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/lesson"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils"
	lspb "github.com/multi-tenants-cms-golang/lms-sys/protogen/lesson"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (ls *LessonService) CreateLesson(ctx context.Context, req *lspb.CreateLessonRequest) (*lspb.CreateLessonResponse, error) {
	ls.logger.WithFields(logrus.Fields{
		"method": "CreateLesson",
	}).Info("Creating Lesson") 

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

	args := repo.GetModuleByIDWithTenantParams{
		ModuleID: pgModuleUUID,
		Namespace: orgName,
	}

	module, err := ls.store.GetModuleByIDWithTenant(dbCtx, args)

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
