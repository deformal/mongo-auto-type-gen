package infer

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Flatten(
	doc map[string]any,
	root map[string]*FieldStats,
	totalDocs *int,
) {
	*totalDocs++
	flattenObject("", doc, root)
}

func flattenObject(prefix string, obj map[string]any, stats map[string]*FieldStats) {
	for key, value := range obj {
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}

		fs, ok := stats[path]
		if !ok {
			fs = NewFieldStats(path)
			stats[path] = fs
		}

		fs.Count++

		kind := DetectKind(value)
		fs.Types[kind]++

		switch kind {
		case KindObject:
			if m, ok := asObject(value); ok {
				flattenObject(path, m, stats)
			}

		case KindArray:
			arr, ok := asArray(value)
			if !ok {
				continue
			}
			for _, elem := range arr {
				elemKind := DetectKind(elem)
				fs.ArrayTypes[elemKind]++

				if elemKind == KindObject {
					if m, ok := asObject(elem); ok {
						flattenObject(path+"[]", m, stats)
					}
				}
			}
		}
	}
}

func asObject(v any) (map[string]any, bool) {
	switch x := v.(type) {
	case map[string]any:
		return x, true
	case *bson.M:
		return map[string]any(*x), true
	case primitive.M:
		return map[string]any(x), true
	case *bson.D:
		return dToMap(*x), true
	case primitive.D:
		return dToMap(x), true
	default:
		return nil, false
	}
}

func dToMap(d primitive.D) map[string]any {
	m := make(map[string]any, len(d))
	for _, e := range d {
		m[e.Key] = e.Value
	}
	return m
}

func asArray(v any) ([]any, bool) {
	switch x := v.(type) {
	case []any:
		return x, true
	case primitive.A:
		return []any(x), true
	default:
		return nil, false
	}
}
