package utils

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func ConvertStringToUUID(uuidStr string) (pgtype.UUID, error) {
	if uuidStr == "" {
		return pgtype.UUID{Valid: false}, nil
	}

	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return pgtype.UUID{Valid: false}, err
	}

	var pgUUID pgtype.UUID
	pgUUID.Bytes = parsedUUID
	pgUUID.Valid = true

	return pgUUID, nil
}

// ConvertTimestamp converts pgtype.Timestamptz to *timestamppb.Timestamp
func ConvertTimestamp(ts pgtype.Timestamptz) *timestamppb.Timestamp {
	if !ts.Valid {
		return nil
	}
	return timestamppb.New(ts.Time)
}
func NullableStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// NullableTimeToProto safely converts sql.NullTime to *timestamppb.Timestamp
func NullableTimeToProto(nt sql.NullTime) *timestamppb.Timestamp {
	if nt.Valid {
		return timestamppb.New(nt.Time)
	}
	return nil
}

// ConvertUUIDToString converts pgtype.UUID to string
func ConvertUUIDToString(pgUUID pgtype.UUID) string {
	if !pgUUID.Valid {
		return ""
	}

	u := uuid.UUID(pgUUID.Bytes)
	return u.String()
}

// NullTimeToTimestamptz converts sql.NullTime to pgtype.Timestamptz
func NullTimeToTimestamptz(nt sql.NullTime) pgtype.Timestamptz {
	if !nt.Valid {
		return pgtype.Timestamptz{}
	}

	return pgtype.Timestamptz{
		Time:  nt.Time,
		Valid: true,
	}
}

// TimestamptzToNullTime converts pgtype.Timestamptz to sql.NullTime
func TimestamptzToNullTime(ts pgtype.Timestamptz) sql.NullTime {
	if !ts.Valid {
		return sql.NullTime{}
	}

	return sql.NullTime{
		Time:  ts.Time,
		Valid: true,
	}
}

// TimeToTimestamptz converts time.Time to pgtype.Timestamptz
func TimeToTimestamptz(t time.Time) pgtype.Timestamptz {
	if t.IsZero() {
		return pgtype.Timestamptz{}
	}

	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}

// TimestamptzToTime converts pgtype.Timestamptz to time.Time
func TimestamptzToTime(ts pgtype.Timestamptz) time.Time {
	if !ts.Valid {
		return time.Time{}
	}
	return ts.Time
}

// NullTimeToProtoTimestamp converts sql.NullTime to *timestamppb.Timestamp
func NullTimeToProtoTimestamp(nt sql.NullTime) *timestamppb.Timestamp {
	if !nt.Valid {
		return nil
	}
	return timestamppb.New(nt.Time)
}

// ProtoTimestampToNullTime converts *timestamppb.Timestamp to sql.NullTime
func ProtoTimestampToNullTime(ts *timestamppb.Timestamp) sql.NullTime {
	if ts == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{
		Time:  ts.AsTime(),
		Valid: true,
	}
}

// TimePointerToProtoTimestamp converts *time.Time to *timestamppb.Timestamp
func TimePointerToProtoTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

// ProtoTimestampToTimePointer converts *timestamppb.Timestamp to *time.Time
func ProtoTimestampToTimePointer(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

// TimeToProtoTimestamp converts time.Time to *timestamppb.Timestamp
func TimeToProtoTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

// ProtoTimestampToTime converts *timestamppb.Timestamp to time.Time
func ProtoTimestampToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

// Additional utility functions for common operations

// NullTimeToTimePointer converts sql.NullTime to *time.Time
func NullTimeToTimePointer(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

// TimePointerToNullTime converts *time.Time to sql.NullTime
func TimePointerToNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{
		Time:  *t,
		Valid: true,
	}
}

// SafeTimeToProto safely converts time.Time to *timestamppb.Timestamp with zero-value check
func SafeTimeToProto(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

// Batch conversion functions for slices

// NullTimesToProtoTimestamps converts []sql.NullTime to []*timestamppb.Timestamp
func NullTimesToProtoTimestamps(nts []sql.NullTime) []*timestamppb.Timestamp {
	if nts == nil {
		return nil
	}

	result := make([]*timestamppb.Timestamp, len(nts))
	for i, nt := range nts {
		result[i] = NullTimeToProtoTimestamp(nt)
	}
	return result
}

// ProtoTimestampsToNullTimes converts []*timestamppb.Timestamp to []sql.NullTime
func ProtoTimestampsToNullTimes(tss []*timestamppb.Timestamp) []sql.NullTime {
	if tss == nil {
		return nil
	}

	result := make([]sql.NullTime, len(tss))
	for i, ts := range tss {
		result[i] = ProtoTimestampToNullTime(ts)
	}
	return result
}
