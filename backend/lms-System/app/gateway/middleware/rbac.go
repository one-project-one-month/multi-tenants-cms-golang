package middleware

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type RBACConfig struct {
	Logger     *logrus.Logger
	RouteRoles map[string][]Role
}

func (c *RBACConfig) RBACMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//if shouldSkipRBAC(r.URL.Path) {
		//	next.ServeHTTP(w, r)
		//	return
		//}

		userCtx := GetUserContext(r)
		if userCtx == nil {
			c.Logger.Warn("Missing user context in RBAC middleware")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		requiredRoles := c.getRequiredRoles(r)
		if len(requiredRoles) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		if !hasAnyRole(userCtx.Roles, requiredRoles) {
			c.Logger.Warnf("User %s lacks required roles for %s", userCtx.UserID, r.URL.Path)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (c *RBACConfig) getRequiredRoles(r *http.Request) []Role {
	if roles, ok := c.RouteRoles[r.URL.Path]; ok {
		return roles
	}

	for route, roles := range c.RouteRoles {
		if strings.HasPrefix(r.URL.Path, route) {
			return roles
		}
	}

	return nil
}

func hasAnyRole(userRoles, requiredRoles []Role) bool {
	for _, userRole := range userRoles {
		for _, requiredRole := range requiredRoles {
			if userRole == requiredRole {
				return true
			}
		}
	}
	return false
}
