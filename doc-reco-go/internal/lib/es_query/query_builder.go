package es_query

import (
	"bytes"
	"context"
)

type QueryBuilder interface {
	CreateQuery(ctx context.Context, searchDict map[string]interface{}, additionalFilter map[string]interface{}) (bytes.Buffer, error)
}
