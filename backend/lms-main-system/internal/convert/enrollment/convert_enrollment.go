package enrollment

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	epb "github.com/multi-tenants-cms-golang/lms-sys/protogen/enrollments"
)

func ConvertStringToGoogleUUID(uuidStr string) uuid.UUID {
	if uuidStr == "" {
		return uuid.Nil
	}
	return uuid.MustParse(uuidStr)
}

func ConvertGoogleUUIDToString(u uuid.UUID) string {
	return u.String()
}

func ConvertTimeToPgTime(t time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{
		Time:  t,
		Valid: t.Compare(time.Now()) == 1,
	}
}

func ConvertEnrollmentDetailsToProto(ed db.EnrollmentDetails) *epb.EnrollmentResponse {
	return &epb.EnrollmentResponse{
		Id: ConvertGoogleUUIDToString(ed.EnrollmentID),
		Student: &epb.Student{
			Id:    ConvertGoogleUUIDToString(ed.StudentID),
			Name:  ed.StudentName,
			Email: ed.StudentEmail,
		},
		Course: &epb.Course{
			Id:    ConvertGoogleUUIDToString(ed.CourseID),
			Title: ed.CourseTitle,
		},
		Category: &epb.Category{
			Id:   ConvertGoogleUUIDToString(ed.CategoryID),
			Name: ed.CategoryName,
		},
	}
}
