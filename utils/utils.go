package utils

import (
	"github.com/gofrs/uuid"
	jsoniter "github.com/json-iterator/go"
)

func RandomUUID() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}

func ToJSON(a any) string {
	r, _ := jsoniter.MarshalToString(a)
	return r
}
