package infer

type Kind string

const (
	KindString   Kind = "string"
	KindNumber   Kind = "number"
	KindBoolean  Kind = "boolean"
	KindNull     Kind = "null"
	KindObject   Kind = "object"
	KindArray    Kind = "array"
	KindDate     Kind = "date"
	KindObjectID Kind = "objectId"
	KindUnknown  Kind = "unknown"
)
