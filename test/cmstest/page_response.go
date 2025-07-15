package cmstest

import (
	"github.com/google/uuid"
	"time"
)

type (
	RequestStatus string
	PageStatus    string
)

const (
	RequestStatusPending  RequestStatus = "PENDING"
	RequestStatusApproved RequestStatus = "APPROVED"
	RequestStatusRejected RequestStatus = "REJECTED"
)

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
