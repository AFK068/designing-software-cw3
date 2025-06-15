package domain

import (
	"encoding/json"

	"github.com/google/uuid"
)

type OutboxRecord struct {
	ID      uuid.UUID
	Payload json.RawMessage
}
