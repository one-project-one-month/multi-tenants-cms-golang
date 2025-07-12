package utils

import "github.com/multi-tenants-cms-golang/cms-sys/internal/types"

func IsValidRequestStatus(status types.RequestStatus) bool {
	switch status {
	case types.RequestStatusPending,
		types.RequestStatusApproved,
		types.RequestStatusRejected:
		return true
	default:
		return false
	}
}
