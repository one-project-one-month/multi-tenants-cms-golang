package types

import (
	"github.com/google/uuid"
	"time"
)

/*
type PageResponse struct {
	ID            uuid.UUID  `json:"id"`
	PageRequestID uuid.UUID  `json:"pageRequestId"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	ImageURL      string     `json:"imageUrl"`
	Status        PageStatus `json:"status"`
	OwnerID       uuid.UUID  `json:"ownerId"`
	PublisherID   *uuid.UUID `json:"publisherId"`
	CreatedAt     time.Time  `json:"createdAt"`
}
*/

type PageResponse struct {
	ID    uuid.UUID  `json:"id"`
	Page  PageInner  `json:"page"`
	Owner OwnerInner `json:"owner"`
}

type PageInner struct {
	PageRequestID uuid.UUID  `json:"pageRequestId"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	ImageURL      *string    `json:"imageUrl,omitempty"`
	Status        PageStatus `json:"status"`
}

type OwnerInner struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	PublisherID *uuid.UUID `json:"publisherId,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
}
