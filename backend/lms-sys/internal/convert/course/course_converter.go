package course

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/global"
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	cpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/courses"
)

func ConvertCourseProto(courseDBObj db.Course) *cpb.CreateCourseResponse {
	return &cpb.CreateCourseResponse{
		CourseCreated: &cpb.Course{
			CourseId:    courseDBObj.CourseID.String(),
			CourseTitle: courseDBObj.CourseTitle,
			CreatedAt:   global.ConvertTimestamp(pgtype.Timestamptz(courseDBObj.CreatedAt)),
			UpdatedAt:   global.ConvertTimestamp(pgtype.Timestamptz(courseDBObj.UpdatedAt)),
		},
	}
}
