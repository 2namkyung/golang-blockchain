package uuid

import (
	"github.com/google/uuid"
)

func MakingUUID() string {
	id := uuid.New()

	return id.String()
}
