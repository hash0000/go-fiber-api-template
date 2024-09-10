package helpers

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

func Int64Ptr(n sql.NullInt64) *int64 {
	if n.Valid {
		return &n.Int64
	}
	return nil
}

func Int16Ptr(n sql.NullInt16) *int16 {
	if n.Valid {
		return &n.Int16
	}
	return nil
}

func BoolPtr(b sql.NullBool) *bool {
	if b.Valid {
		return &b.Bool
	}
	return nil
}

func StringPtr(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

func UuidPtr(s sql.NullString) *uuid.UUID {
	if s.Valid {
		u, err := uuid.Parse(s.String)
		if err == nil {
			return &u
		}
	}
	return nil
}

func TimePtr(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func JsonRawMessagePtr(ns sql.NullString) *json.RawMessage {
	if ns.Valid {
		raw := json.RawMessage(ns.String)
		return &raw
	}
	return nil
}
