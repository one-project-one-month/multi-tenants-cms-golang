package lesson

import (
	lpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/lesson"
	"google.golang.org/protobuf/types/known/timestamppb"

	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
)

func ConvertLessonToProto(lesson *db.GetLessonsByModuleIdRow) *lpb.Lesson {
	materialType := ""

	if lesson.MaterialType.Valid {
		materialType = string(lesson.MaterialType.MaterialType)
	}

	return &lpb.Lesson{
		Id:           lesson.LessonID.String(),
		Title:        lesson.Title,
		Content:      lesson.Content.String,
		MaterialType: materialType,
		Module: &lpb.Module{
			Id:   lesson.ModuleID.String(),
			Name: lesson.ModuleName,
		},
		CreatedAt: timestamppb.New(lesson.CreatedAt.Time),
		UpdatedAt: timestamppb.New(lesson.UpdatedAt.Time),
	}
}

func ConvertLessonToProtoFromLesson(l db.Lesson, moduleName string) *lpb.Lesson {
	materialType := ""
	if l.MaterialType.Valid {
		materialType = string(l.MaterialType.MaterialType)
	}

	return &lpb.Lesson{
		Id:           l.LessonID.String(),
		Title:        l.Title,
		Content:      l.Content.String,
		MaterialType: materialType,
		Module: &lpb.Module{
			Id:   l.ModuleID.String(),
			Name: moduleName,
		},
		CreatedAt: timestamppb.New(l.CreatedAt.Time),
		UpdatedAt: timestamppb.New(l.UpdatedAt.Time),
	}
}

func ConvertLessonsToProto(lessons []db.GetLessonsByModuleIdRow) []*lpb.Lesson {
	protoLessons := make([]*lpb.Lesson, len(lessons))
	for _, lesson := range lessons {
		protoLessons = append(protoLessons, ConvertLessonToProto(&lesson))
	}
	return protoLessons
}
