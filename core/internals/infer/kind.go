package infer

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DetectKind(v any) Kind {
	if v == nil {
		return KindNull
	}

	switch v.(type) {
	case string:
		return KindString
	case bool:
		return KindBoolean
	case int, int32, int64, float32, float64:
		return KindNumber

	case primitive.DateTime:
		return KindDate
	case primitive.ObjectID:
		return KindObjectID

	case []any, primitive.A:
		return KindArray

	case map[string]any, *bson.M, primitive.M, *bson.D, primitive.D:
		return KindObject
	}

	rk := reflect.TypeOf(v).Kind()
	if rk == reflect.Map {
		return KindObject
	}
	if rk == reflect.Slice || rk == reflect.Array {
		return KindArray
	}

	return KindUnknown
}
