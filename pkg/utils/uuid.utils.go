package utils

import (
	fuuid "github.com/gofrs/uuid/v5"
	guuid "github.com/google/uuid"
)

func GenerateUUIDv7() guuid.UUID {
	id, err := fuuid.NewV7()
	if err != nil {
		panic(err)
	}
	return guuid.UUID(id)
}
