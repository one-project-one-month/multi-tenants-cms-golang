package middleware

import "net/http"

type Role string

const (
	LMSAdmin   Role = "LMS_ADMIN"
	Student    Role = "STUDENT"
	Instructor Role = "INSTRUCTOR"
)

type UserContext struct {
	UserID string
	Roles  []Role
}

type contextKey string

const (
	UserContextKey contextKey = "user_context"
)

func GetUserContext(r *http.Request) *UserContext {
	if ctx := r.Context().Value(UserContextKey); ctx != nil {
		return ctx.(*UserContext)
	}
	return nil
}
