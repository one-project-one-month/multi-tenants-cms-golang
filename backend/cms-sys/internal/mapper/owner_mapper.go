package mapper

import (
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
)

func ToOwnerResponse(user *types.CMSUser) *types.OwnerResponse {

	ownerResp := &types.OwnerResponse{
		ID:                   user.CMSUserID,
		Name:                 user.CMSUserName,
		Email:                user.CMSUserEmail,
		NameSpace:            derefString(user.CMSNameSpace),
		Verified:             user.Verified,
		NumberOfRequestPages: &user.NumberOfRequestPages,
		NumberOfPagesOwned:   &user.NumberOfPagesOwned,
	}

	if user.CMSUserRole != "" {
		roleName := user.Role.RoleName
		ownerResp.Role = &roleName
	}

	return ownerResp
}

func ToOwnerListResponse(users []types.CMSUser) []*types.OwnerResponse {
	var ownerList []*types.OwnerResponse
	for _, user := range users {
		ownerList = append(ownerList, ToOwnerResponse(&user))
	}
	return ownerList
}

func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
