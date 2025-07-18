package mapper

import (
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
)

/*
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
*/
func ToPageResponse(p *types.Page) *types.PageResponse {
	return &types.PageResponse{
		ID: p.PageID,
		Page: types.PageInner{
			PageRequestID: p.PageRequestID,
			Title:         p.Title,
			Content:       p.Content,
			ImageURL:      p.ImageURL,
			Status:        p.Status,
		},
		Owner: types.OwnerInner{
			ID:        p.OwnerUser.CMSUserID,
			Name:      p.OwnerUser.CMSUserName,
			Email:     p.OwnerUser.CMSUserEmail,
			CreatedAt: p.OwnerUser.CreatedAt,
		},
	}
}
