package module

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	"google.golang.org/protobuf/types/known/timestamppb"

	mpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/modules"
)

func ConvertModuleToProto(module db.Module) *mpb.Module {
	return &mpb.Module{
		ModuleId:    module.ModuleID.String(),
		ModuleName:  module.ModuleName,
		CourseId:    module.CourseID.String(),
		Description: module.Description.String,
		CreatedAt:   timestamppb.New(module.CreatedAt.Time),
		UpdatedAt:   timestamppb.New(module.UpdatedAt.Time),
	}
}

func ConvertModulesToProto(modules []db.Module) []*mpb.Module {
	protoModules := make([]*mpb.Module, len(modules))

	for _, module := range modules {
		protoModules = append(protoModules, ConvertModuleToProto(module))
	}

	return protoModules
}

func ConvertStringToGoogleUUID(uuidStr string) uuid.UUID {
	if uuidStr == "" {
		return uuid.Nil
	}
	return uuid.MustParse(uuidStr)
}

func ConvertStringToUUID(uuidStr string) (pgtype.UUID, error) {
	if uuidStr == "" {
		return pgtype.UUID{Valid: false}, nil
	}

	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return pgtype.UUID{Valid: false}, err
	}

	pgUUID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}

	return pgUUID, nil
}

func ConvertStringToText(str string) pgtype.Text {
	pgText := pgtype.Text{
		String: str,
		Valid:  false,
	}

	if str != "" {
		pgText = pgtype.Text{
			Valid: true,
		}
	}

	return pgText
}

func ConvertPgBoolToBool(pgbool pgtype.Bool) bool {
	if pgbool.Valid {
		return pgbool.Bool
	}

	return false
}
