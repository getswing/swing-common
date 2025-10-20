package utils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func UuidToStringValue(id *uuid.UUID) *wrapperspb.StringValue {
	if id == nil {
		return nil
	}
	return wrapperspb.String(id.String())
}

func StringPointerToStringValue(s *string) *wrapperspb.StringValue {
	if s == nil {
		return nil
	}
	return wrapperspb.String(*s)
}

func BoolPointerToBoolValue(b *bool) *wrapperspb.BoolValue {
	if b == nil {
		return nil
	}
	return wrapperspb.Bool(*b)
}

func IntPointerToInt32Value(i *int) *wrapperspb.Int32Value {
	if i == nil {
		return nil
	}
	return wrapperspb.Int32(int32(*i))
}

func TimePointerToTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil || t.IsZero() {
		return nil
	}
	return timestamppb.New(*t)
}

func PathToS3File(path string) string {
	return fmt.Sprintf("%s%s", "__s3Url__", path)
}
