package uuid

import "github.com/google/uuid"

func GenerateUID() string {
	uid := uuid.New()

	return uid.String()
}
