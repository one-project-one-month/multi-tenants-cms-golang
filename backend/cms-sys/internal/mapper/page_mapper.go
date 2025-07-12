package mapper

import (
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/utils"
)

func ToPageResponse(page *types.Page) *types.PageResponse {
	return &types.PageResponse{
		ID:            page.PageID,
		PageRequestID: page.PageRequestID,
		Title:         page.Title,
		Content:       utils.SafeString(&page.Content),
		ImageURL:      *page.ImageURL,
		Status:        page.Status,
		OwnerID:       page.OwnerID,
		PublisherID:   page.PublishedByStaffID,
		CreatedAt:     page.CreatedAt,
	}
}
