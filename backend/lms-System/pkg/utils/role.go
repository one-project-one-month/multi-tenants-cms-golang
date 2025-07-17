package utils

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func ParseTimestamp(value interface{}) *timestamppb.Timestamp {
	if value == nil {
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return nil
	}
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return nil
	}
	return timestamppb.New(t)
}
