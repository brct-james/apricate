// Package uuid defines helper functions for generating uuids
package uuid

import (
	"github.com/google/uuid"
)

func NewUUID() (string) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return uuid.String()
}
