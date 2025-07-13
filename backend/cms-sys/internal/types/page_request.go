package types

import "github.com/google/uuid"

type PageCreateRequest struct {
	PageRequestID uuid.UUID `json:"pageRequestId" validate:"required,uuid"`
	Title         string    `json:"title" validate:"required,min=2"`
	Content       string    `json:"content" validate:"required,min=3"`
	ImageUrl      string    `json:"imageUrl" validate:"required,url"`
	Status        *string   `json:"status" validate:"omitempty,oneof=DRAFT PUBLISHED ARCHIVED"`
	OwnerId       uuid.UUID `json:"ownerId" validate:"required,uuid"`
	PublisherId   uuid.UUID `json:"publisherId" validate:"required,uuid"`
}

type PageUpdateRequest struct {
	PageRequestID uuid.UUID `json:"pageRequestId" validate:"required,uuid"`
	Title         string    `json:"title" validate:"required,min=2"`
	Content       string    `json:"content" validate:"required,min=3"`
	ImageUrl      string    `json:"imageUrl" validate:"required,url"`
	Status        string    `json:"status" validate:"required,oneof=DRAFT PUBLISHED ARCHIVED"`
	OwnerId       uuid.UUID `json:"ownerId" validate:"required,uuid"`
	PublisherId   uuid.UUID `json:"publisherId" validate:"required,uuid"`
}
