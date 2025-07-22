package types

import "time"

type CourseStatus string

const (
	StatusPending     CourseStatus = "Pending"
	StatusPublished   CourseStatus = "Published"
	StatusUnpublished CourseStatus = "Unpublished"
	StatusArchived    CourseStatus = "Archived"
)

type Course struct {
	ID               int64        `json:"id" gorm:"primaryKey;autoIncrement"`
	Title            string       `json:"title"`
	Description      string       `json:"description"`
	CategoryID       int64        `json:"category_id"`
	InstructorID     int64        `json:"instructor_id"`
	OverallRating    int          `json:"overall_rating"`
	Status           CourseStatus `json:"status" gorm:"type:enum('Pending','Published','Unpublished','Archived')"`
	DurationDayCount int          `json:"duration_day_count"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
	DeletedAt        *time.Time   `json:"deleted_at" gorm:"index"`
}
