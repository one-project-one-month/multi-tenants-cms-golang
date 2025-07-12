package mapper

import (
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/utils"
)

func ToPageRequestResponse(model *types.PageRequest) *types.PageRequestResponse {
	return &types.PageRequestResponse{
		ID:          model.RequestID,
		OwnerID:     model.OwnerID,
		RequestType: model.RequestType,
		Title:       model.Title,
		Description: utils.SafeString(model.Description),
		Status:      model.Status,
		PageUrl:     model.PageUrl,
		LogoUrl:     model.LogoUrl,
		CreatedAt:   model.CreatedAt,
	}
}
