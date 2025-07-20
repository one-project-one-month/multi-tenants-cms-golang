package global

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
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

// ConvertUUIDToString converts pgtype.UUID to string
func ConvertUUIDToString(pgUUID pgtype.UUID) string {
	if !pgUUID.Valid {
		return ""
	}

	u := uuid.UUID(pgUUID.Bytes)
	return u.String()
}

func ConvertStringToPgText(uuidStr string) *pgtype.Text {
	if uuidStr == "" {
		return nil
	}
	return &pgtype.Text{
		String: uuidStr,
		Valid:  true,
	}
}

func ConvertStringToGoogleUUID(uuidStr string) uuid.UUID {
	if uuidStr == "" {
		return uuid.Nil
	}
	return uuid.MustParse(uuidStr)
}
